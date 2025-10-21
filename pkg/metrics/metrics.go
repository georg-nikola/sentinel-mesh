package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics holds all Prometheus metrics for the application
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPResponseSize      *prometheus.HistogramVec

	// Data collection metrics
	MetricsCollected      *prometheus.CounterVec
	CollectionDuration    *prometheus.HistogramVec
	CollectionErrors      *prometheus.CounterVec

	// Processing metrics
	MessagesProcessed     *prometheus.CounterVec
	ProcessingDuration    *prometheus.HistogramVec
	ProcessingErrors      *prometheus.CounterVec

	// Storage metrics
	DatabaseWrites        *prometheus.CounterVec
	DatabaseReads         *prometheus.CounterVec
	DatabaseErrors        *prometheus.CounterVec
	DatabaseConnections   *prometheus.GaugeVec

	// ML metrics
	MLPredictions         *prometheus.CounterVec
	MLInferenceDuration   *prometheus.HistogramVec
	AnomaliesDetected     *prometheus.CounterVec

	// System metrics
	GoRoutines            prometheus.Gauge
	MemoryUsage           prometheus.Gauge
	CPUUsage              prometheus.Gauge
}

// New creates a new Metrics instance with all metrics registered
func New(namespace, subsystem string) *Metrics {
	m := &Metrics{
		// HTTP metrics
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status_code"},
		),

		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),

		HTTPResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(1024, 2, 10),
			},
			[]string{"method", "endpoint"},
		),

		// Data collection metrics
		MetricsCollected: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "metrics_collected_total",
				Help:      "Total number of metrics collected",
			},
			[]string{"source", "type"},
		),

		CollectionDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "collection_duration_seconds",
				Help:      "Time spent collecting metrics",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"source", "type"},
		),

		CollectionErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "collection_errors_total",
				Help:      "Total number of collection errors",
			},
			[]string{"source", "type", "error"},
		),

		// Processing metrics
		MessagesProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "messages_processed_total",
				Help:      "Total number of messages processed",
			},
			[]string{"topic", "status"},
		),

		ProcessingDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "processing_duration_seconds",
				Help:      "Time spent processing messages",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"topic", "processor"},
		),

		ProcessingErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "processing_errors_total",
				Help:      "Total number of processing errors",
			},
			[]string{"topic", "processor", "error"},
		),

		// Storage metrics
		DatabaseWrites: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "database_writes_total",
				Help:      "Total number of database writes",
			},
			[]string{"database", "table", "status"},
		),

		DatabaseReads: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "database_reads_total",
				Help:      "Total number of database reads",
			},
			[]string{"database", "table", "status"},
		),

		DatabaseErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "database_errors_total",
				Help:      "Total number of database errors",
			},
			[]string{"database", "operation", "error"},
		),

		DatabaseConnections: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "database_connections_active",
				Help:      "Number of active database connections",
			},
			[]string{"database"},
		),

		// ML metrics
		MLPredictions: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "ml_predictions_total",
				Help:      "Total number of ML predictions made",
			},
			[]string{"model", "result"},
		),

		MLInferenceDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "ml_inference_duration_seconds",
				Help:      "Time spent on ML inference",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"model"},
		),

		AnomaliesDetected: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "anomalies_detected_total",
				Help:      "Total number of anomalies detected",
			},
			[]string{"service", "type", "severity"},
		),

		// System metrics
		GoRoutines: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "goroutines_active",
				Help:      "Number of active goroutines",
			},
		),

		MemoryUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "memory_usage_bytes",
				Help:      "Memory usage in bytes",
			},
		),

		CPUUsage: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "cpu_usage_percent",
				Help:      "CPU usage percentage",
			},
		),
	}

	return m
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(method, endpoint, statusCode string, duration time.Duration, responseSize int) {
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	m.HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	m.HTTPResponseSize.WithLabelValues(method, endpoint).Observe(float64(responseSize))
}

// RecordMetricCollection records metric collection metrics
func (m *Metrics) RecordMetricCollection(source, metricType string, duration time.Duration, count int, err error) {
	m.MetricsCollected.WithLabelValues(source, metricType).Add(float64(count))
	m.CollectionDuration.WithLabelValues(source, metricType).Observe(duration.Seconds())
	
	if err != nil {
		m.CollectionErrors.WithLabelValues(source, metricType, err.Error()).Inc()
	}
}

// RecordMessageProcessing records message processing metrics
func (m *Metrics) RecordMessageProcessing(topic, processor, status string, duration time.Duration, err error) {
	m.MessagesProcessed.WithLabelValues(topic, status).Inc()
	m.ProcessingDuration.WithLabelValues(topic, processor).Observe(duration.Seconds())
	
	if err != nil {
		m.ProcessingErrors.WithLabelValues(topic, processor, err.Error()).Inc()
	}
}

// RecordDatabaseOperation records database operation metrics
func (m *Metrics) RecordDatabaseOperation(database, operation, table, status string, err error) {
	switch operation {
	case "write":
		m.DatabaseWrites.WithLabelValues(database, table, status).Inc()
	case "read":
		m.DatabaseReads.WithLabelValues(database, table, status).Inc()
	}
	
	if err != nil {
		m.DatabaseErrors.WithLabelValues(database, operation, err.Error()).Inc()
	}
}

// RecordMLPrediction records ML prediction metrics
func (m *Metrics) RecordMLPrediction(model, result string, duration time.Duration) {
	m.MLPredictions.WithLabelValues(model, result).Inc()
	m.MLInferenceDuration.WithLabelValues(model).Observe(duration.Seconds())
}

// RecordAnomaly records anomaly detection metrics
func (m *Metrics) RecordAnomaly(service, anomalyType, severity string) {
	m.AnomaliesDetected.WithLabelValues(service, anomalyType, severity).Inc()
}

// UpdateSystemMetrics updates system-level metrics
func (m *Metrics) UpdateSystemMetrics(goroutines int, memoryBytes, cpuPercent float64) {
	m.GoRoutines.Set(float64(goroutines))
	m.MemoryUsage.Set(memoryBytes)
	m.CPUUsage.Set(cpuPercent)
}

// UpdateDatabaseConnections updates database connection metrics
func (m *Metrics) UpdateDatabaseConnections(database string, connections int) {
	m.DatabaseConnections.WithLabelValues(database).Set(float64(connections))
}