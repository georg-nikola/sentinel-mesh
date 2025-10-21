package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sentinel-mesh/sentinel-mesh/internal/services"
	"github.com/sentinel-mesh/sentinel-mesh/pkg/config"
	"github.com/sentinel-mesh/sentinel-mesh/pkg/metrics"
)

var (
	configPath string
	version    = "dev"
	buildTime  = "unknown"
	gitCommit  = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "collector",
		Short: "Sentinel Mesh Data Collector Service",
		Long:  `A high-performance data collection service for Kubernetes metrics and logs`,
		Run:   run,
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to configuration file")
	rootCmd.AddCommand(versionCmd())

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Failed to execute command")
	}
}

func run(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration")
	}

	// Setup logging
	setupLogging(cfg.Logging)

	logrus.WithFields(logrus.Fields{
		"version":    version,
		"build_time": buildTime,
		"git_commit": gitCommit,
	}).Info("Starting Sentinel Mesh Collector")

	// Initialize metrics
	m := metrics.New("sentinel_mesh", "collector")

	// Create collector service
	collector, err := services.NewCollector(cfg, m)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create collector service")
	}

	// Start metrics server
	metricsServer := startMetricsServer(cfg.Monitoring.Prometheus.Port)

	// Start collector
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	collectorErr := make(chan error, 1)
	go func() {
		if err := collector.Start(ctx); err != nil {
			collectorErr <- fmt.Errorf("collector service failed: %w", err)
		}
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigChan:
		logrus.WithField("signal", sig).Info("Received shutdown signal")
	case err := <-collectorErr:
		logrus.WithError(err).Error("Collector service error")
	}

	// Graceful shutdown
	logrus.Info("Shutting down collector service...")
	cancel()

	// Stop metrics server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		logrus.WithError(err).Error("Failed to shutdown metrics server")
	}

	// Wait for collector to stop
	if err := collector.Stop(shutdownCtx); err != nil {
		logrus.WithError(err).Error("Failed to stop collector gracefully")
	}

	logrus.Info("Collector service stopped")
}

func setupLogging(cfg config.LoggingConfig) {
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		logrus.WithError(err).Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	if cfg.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	if cfg.Output != "stdout" && cfg.Output != "" {
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logrus.WithError(err).Fatal("Failed to open log file")
		}
		logrus.SetOutput(file)
	}
}

func startMetricsServer(port int) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		logrus.WithField("port", port).Info("Starting metrics server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Metrics server failed")
		}
	}()

	return server
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Sentinel Mesh Collector\n")
			fmt.Printf("Version: %s\n", version)
			fmt.Printf("Build Time: %s\n", buildTime)
			fmt.Printf("Git Commit: %s\n", gitCommit)
		},
	}
}