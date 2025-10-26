package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"

	"github.com/sentinel-mesh/sentinel-mesh/internal/models"
	"github.com/sentinel-mesh/sentinel-mesh/pkg/config"
	"github.com/sentinel-mesh/sentinel-mesh/pkg/metrics"
)

// Collector handles data collection from various sources
type Collector struct {
	config        *config.Config
	metrics       *metrics.Metrics
	kubeClient    kubernetes.Interface
	metricsClient versioned.Interface
	kafkaWriter   *kafka.Writer
	
	// Collection intervals
	nodeMetricsInterval time.Duration
	podMetricsInterval  time.Duration
	eventInterval       time.Duration
	
	// Collectors
	nodeCollector  *NodeCollector
	podCollector   *PodCollector
	eventCollector *EventCollector
	
	// State
	running bool
	wg      sync.WaitGroup
	logger  *logrus.Entry
}

// NewCollector creates a new collector service
func NewCollector(cfg *config.Config, m *metrics.Metrics) (*Collector, error) {
	logger := logrus.WithField("service", "collector")
	
	// Initialize Kubernetes clients
	kubeClient, metricsClient, err := initKubernetesClients(cfg.Kubernetes)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Kubernetes clients: %w", err)
	}
	
	// Initialize Kafka writer
	kafkaWriter := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        cfg.Kafka.Topics["metrics"],
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Compression:  kafka.Snappy,
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
	}
	
	collector := &Collector{
		config:              cfg,
		metrics:             m,
		kubeClient:          kubeClient,
		metricsClient:       metricsClient,
		kafkaWriter:         kafkaWriter,
		nodeMetricsInterval: 30 * time.Second,
		podMetricsInterval:  15 * time.Second,
		eventInterval:       5 * time.Second,
		logger:              logger,
	}
	
	// Initialize sub-collectors
	collector.nodeCollector = NewNodeCollector(kubeClient, metricsClient, m, logger)
	collector.podCollector = NewPodCollector(kubeClient, metricsClient, m, logger)
	collector.eventCollector = NewEventCollector(kubeClient, m, logger)
	
	return collector, nil
}

// Start begins the collection process
func (c *Collector) Start(ctx context.Context) error {
	c.logger.Info("Starting collector service")
	c.running = true
	
	// Start system metrics collection
	c.wg.Add(1)
	go c.collectSystemMetrics(ctx)
	
	// Start node metrics collection
	c.wg.Add(1)
	go c.collectNodeMetrics(ctx)
	
	// Start pod metrics collection
	c.wg.Add(1)
	go c.collectPodMetrics(ctx)
	
	// Start event collection
	c.wg.Add(1)
	go c.collectEvents(ctx)
	
	c.logger.Info("Collector service started")
	
	// Wait for context cancellation
	<-ctx.Done()
	c.logger.Info("Collector service stopping...")
	
	return nil
}

// Stop gracefully stops the collector
func (c *Collector) Stop(ctx context.Context) error {
	c.logger.Info("Stopping collector service")
	c.running = false
	
	// Wait for all goroutines to finish
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		c.logger.Info("All collection routines stopped")
	case <-ctx.Done():
		c.logger.Warn("Timeout waiting for collection routines to stop")
	}
	
	// Close Kafka writer
	if err := c.kafkaWriter.Close(); err != nil {
		c.logger.WithError(err).Error("Failed to close Kafka writer")
	}
	
	c.logger.Info("Collector service stopped")
	return nil
}

// collectSystemMetrics collects system-level metrics
func (c *Collector) collectSystemMetrics(ctx context.Context) {
	defer c.wg.Done()
	
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			start := time.Now()
			
			// Collect Go runtime metrics
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			
			goroutines := runtime.NumGoroutine()
			memoryBytes := float64(memStats.Alloc)
			
			// Update metrics
			c.metrics.UpdateSystemMetrics(goroutines, memoryBytes, 0) // CPU collection would need more implementation
			
			duration := time.Since(start)
			c.metrics.RecordMetricCollection("system", "runtime", duration, 3, nil)
			
			c.logger.WithFields(logrus.Fields{
				"goroutines":   goroutines,
				"memory_mb":    memoryBytes / 1024 / 1024,
				"duration_ms":  duration.Milliseconds(),
			}).Debug("Collected system metrics")
		}
	}
}

// collectNodeMetrics collects Kubernetes node metrics
func (c *Collector) collectNodeMetrics(ctx context.Context) {
	defer c.wg.Done()
	
	ticker := time.NewTicker(c.nodeMetricsInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			start := time.Now()
			
			nodeMetrics, err := c.nodeCollector.Collect(ctx)
			if err != nil {
				c.logger.WithError(err).Error("Failed to collect node metrics")
				c.metrics.RecordMetricCollection("kubernetes", "nodes", time.Since(start), 0, err)
				continue
			}
			
			// Send metrics to Kafka
			for _, metric := range nodeMetrics {
				if err := c.sendMetricToKafka(ctx, metric); err != nil {
					c.logger.WithError(err).Error("Failed to send node metric to Kafka")
				}
			}
			
			duration := time.Since(start)
			c.metrics.RecordMetricCollection("kubernetes", "nodes", duration, len(nodeMetrics), nil)
			
			c.logger.WithFields(logrus.Fields{
				"node_count":   len(nodeMetrics),
				"duration_ms":  duration.Milliseconds(),
			}).Debug("Collected node metrics")
		}
	}
}

// collectPodMetrics collects Kubernetes pod metrics
func (c *Collector) collectPodMetrics(ctx context.Context) {
	defer c.wg.Done()
	
	ticker := time.NewTicker(c.podMetricsInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			start := time.Now()
			
			podMetrics, err := c.podCollector.Collect(ctx)
			if err != nil {
				c.logger.WithError(err).Error("Failed to collect pod metrics")
				c.metrics.RecordMetricCollection("kubernetes", "pods", time.Since(start), 0, err)
				continue
			}
			
			// Send metrics to Kafka
			for _, metric := range podMetrics {
				if err := c.sendMetricToKafka(ctx, metric); err != nil {
					c.logger.WithError(err).Error("Failed to send pod metric to Kafka")
				}
			}
			
			duration := time.Since(start)
			c.metrics.RecordMetricCollection("kubernetes", "pods", duration, len(podMetrics), nil)
			
			c.logger.WithFields(logrus.Fields{
				"pod_count":    len(podMetrics),
				"duration_ms":  duration.Milliseconds(),
			}).Debug("Collected pod metrics")
		}
	}
}

// collectEvents collects Kubernetes events
func (c *Collector) collectEvents(ctx context.Context) {
	defer c.wg.Done()
	
	ticker := time.NewTicker(c.eventInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			start := time.Now()
			
			events, err := c.eventCollector.Collect(ctx)
			if err != nil {
				c.logger.WithError(err).Error("Failed to collect events")
				c.metrics.RecordMetricCollection("kubernetes", "events", time.Since(start), 0, err)
				continue
			}
			
			// Send events to Kafka logs topic
			for _, event := range events {
				if err := c.sendEventToKafka(ctx, event); err != nil {
					c.logger.WithError(err).Error("Failed to send event to Kafka")
				}
			}
			
			duration := time.Since(start)
			c.metrics.RecordMetricCollection("kubernetes", "events", duration, len(events), nil)
			
			if len(events) > 0 {
				c.logger.WithFields(logrus.Fields{
					"event_count":  len(events),
					"duration_ms":  duration.Milliseconds(),
				}).Debug("Collected events")
			}
		}
	}
}

// sendMetricToKafka sends a metric to Kafka
func (c *Collector) sendMetricToKafka(ctx context.Context, metric *models.Metric) error {
	data, err := json.Marshal(metric)
	if err != nil {
		return fmt.Errorf("failed to marshal metric: %w", err)
	}
	
	message := kafka.Message{
		Key:   []byte(metric.Name),
		Value: data,
		Time:  metric.Timestamp,
	}
	
	return c.kafkaWriter.WriteMessages(ctx, message)
}

// sendEventToKafka sends an event to Kafka logs topic
func (c *Collector) sendEventToKafka(ctx context.Context, event *models.Event) error {
	// Create a writer for the logs topic
	writer := &kafka.Writer{
		Addr:         kafka.TCP(c.config.Kafka.Brokers...),
		Topic:        c.config.Kafka.Topics["logs"],
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
	}
	defer writer.Close()
	
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	
	message := kafka.Message{
		Key:   []byte(event.Type),
		Value: data,
		Time:  event.Timestamp,
	}
	
	return writer.WriteMessages(ctx, message)
}

// initKubernetesClients initializes Kubernetes clients
func initKubernetesClients(cfg config.KubernetesConfig) (kubernetes.Interface, versioned.Interface, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first if requested
	if cfg.InCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			// Log the in-cluster config error but don't fail yet
			logrus.WithError(err).Warn("Failed to load in-cluster config, will try file-based config")
		}
	}

	// If in-cluster config failed or wasn't requested, try file-based config
	if config == nil {
		configPath := cfg.ConfigPath
		if configPath == "" {
			// Expand HOME environment variable properly
			homeDir := os.Getenv("HOME")
			if homeDir == "" {
				homeDir = "/root"
			}
			configPath = homeDir + "/.kube/config"
		}

		// Check if the config file exists before trying to use it
		if _, statErr := os.Stat(configPath); statErr == nil {
			config, err = clientcmd.BuildConfigFromFlags("", configPath)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to build config from file %s: %w", configPath, err)
			}
		} else if cfg.InCluster {
			// If we're supposed to be in-cluster and file doesn't exist, return the original error
			return nil, nil, fmt.Errorf("in-cluster config failed and no kubeconfig file found: %w", err)
		} else {
			return nil, nil, fmt.Errorf("kubeconfig file not found at %s: %w", configPath, statErr)
		}
	}

	if config == nil {
		return nil, nil, fmt.Errorf("failed to create Kubernetes config: no valid configuration found")
	}

	// Create standard Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// Create metrics client
	metricsClient, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	return kubeClient, metricsClient, nil
}