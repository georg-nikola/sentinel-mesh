package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sentinel-mesh/sentinel-mesh/internal/handlers"
	"github.com/sentinel-mesh/sentinel-mesh/pkg/config"
)

var (
	configPath string
	version    = "dev"
	buildTime  = "unknown"
	gitCommit  = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "api",
		Short: "Sentinel Mesh API Service",
		Long:  `REST API service for querying metrics and analytics`,
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
	}).Info("Starting Sentinel Mesh API")

	// Create router
	router := mux.NewRouter()

	// Register handlers
	healthHandler := handlers.NewHealthHandler()
	metricsHandler := handlers.NewMetricsHandler()
	queryHandler := handlers.NewQueryHandler()

	// Infrastructure handler with fallback
	var infrastructureHandler *handlers.InfrastructureHandler
	if infra, err := handlers.NewInfrastructureHandler(); err == nil {
		infrastructureHandler = infra
	} else {
		logrus.WithError(err).Warn("Failed to initialize infrastructure handler, using fallback")
	}

	// Health endpoints
	router.HandleFunc("/health", healthHandler.Health).Methods("GET")
	router.HandleFunc("/readiness", healthHandler.Readiness).Methods("GET")
	router.HandleFunc("/ready", healthHandler.Readiness).Methods("GET") // Kubernetes-style alias

	// Metrics endpoints
	router.HandleFunc("/v1/metrics", metricsHandler.ListMetrics).Methods("GET")
	router.HandleFunc("/v1/metrics/query", metricsHandler.QueryMetrics).Methods("POST")

	// Query endpoints
	router.HandleFunc("/v1/query/range", queryHandler.QueryRange).Methods("POST")
	router.HandleFunc("/v1/query/instant", queryHandler.QueryInstant).Methods("POST")

	// Analytics endpoints
	router.HandleFunc("/v1/analytics/slo", metricsHandler.GetSLOStatus).Methods("GET")
	router.HandleFunc("/v1/analytics/anomalies", metricsHandler.GetAnomalies).Methods("GET")

	// Infrastructure endpoints
	if infrastructureHandler != nil {
		router.HandleFunc("/v1/infrastructure", infrastructureHandler.GetInfrastructure).Methods("GET")
	}

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler())

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.API.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logrus.WithField("port", cfg.API.Port).Info("Starting API server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("API server failed")
		}
	}()

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logrus.WithField("signal", sig).Info("Received shutdown signal")

	// Graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.WithError(err).Error("Failed to shutdown API server gracefully")
	}

	logrus.Info("API server stopped")
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
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Sentinel Mesh API\n")
			fmt.Printf("Version: %s\n", version)
			fmt.Printf("Build Time: %s\n", buildTime)
			fmt.Printf("Git Commit: %s\n", gitCommit)
		},
	}
}
