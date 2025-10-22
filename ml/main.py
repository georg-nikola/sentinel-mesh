#!/usr/bin/env python3
"""
Sentinel Mesh ML Service
Anomaly Detection and Predictive Analytics for Kubernetes Monitoring
"""

import asyncio
import logging
import os
import signal
import sys
from typing import Dict, Any

from flask import Flask
from flask_cors import CORS

from inference.anomaly_detector import AnomalyDetector
from inference.predictor import ResourcePredictor
from api.routes import create_api_blueprint
from core.config import load_config
from core.kafka_consumer import KafkaConsumerManager
from core.metrics import setup_prometheus_metrics
from core.logging import setup_logging


class MLService:
    """Main ML Service class that orchestrates anomaly detection and prediction."""

    def __init__(self, config: Dict[str, Any]):
        self.config = config
        self.app = Flask(__name__)
        self.logger = logging.getLogger(__name__)

        # Initialize components
        self.anomaly_detector = AnomalyDetector(config)
        self.resource_predictor = ResourcePredictor(config)
        self.kafka_consumer = KafkaConsumerManager(config, self._process_metrics)

        # Setup Flask app
        self._setup_flask_app()

        # Track running state
        self.running = False

    def _setup_flask_app(self):
        """Setup Flask application with routes and middleware."""
        # Enable CORS
        CORS(self.app, origins=self.config.get("cors", {}).get("origins", ["*"]))

        # Setup Prometheus metrics
        setup_prometheus_metrics(self.app)

        # Register API blueprint
        api_blueprint = create_api_blueprint(
            self.anomaly_detector, self.resource_predictor, self.config
        )
        self.app.register_blueprint(api_blueprint, url_prefix="/api/v1")

        # Health check endpoint
        @self.app.route("/health")
        def health_check():
            return {
                "status": "healthy",
                "service": "ml-service",
                "version": self.config.get("version", "1.0.0"),
            }

        # Readiness check endpoint
        @self.app.route("/ready")
        def readiness_check():
            is_ready = (
                self.anomaly_detector.is_ready() and self.resource_predictor.is_ready()
            )

            if is_ready:
                return {"status": "ready"}, 200
            else:
                return {"status": "not ready"}, 503

    async def _process_metrics(self, metrics_batch: list):
        """Process a batch of metrics for anomaly detection and prediction."""
        try:
            self.logger.debug(f"Processing batch of {len(metrics_batch)} metrics")

            # Process metrics for anomaly detection
            anomalies = await self.anomaly_detector.detect_anomalies(metrics_batch)

            if anomalies:
                self.logger.info(f"Detected {len(anomalies)} anomalies")
                # Send anomalies to alerting system via Kafka
                await self._send_anomalies_to_kafka(anomalies)

            # Update prediction models with new data
            await self.resource_predictor.update_models(metrics_batch)

        except Exception as e:
            self.logger.error(f"Error processing metrics batch: {e}", exc_info=True)

    async def _send_anomalies_to_kafka(self, anomalies: list):
        """Send detected anomalies to Kafka for alerting."""
        try:
            # This would be implemented to send to the alerts topic
            self.logger.debug(f"Sending {len(anomalies)} anomalies to Kafka")
            # Implementation would go here
        except Exception as e:
            self.logger.error(f"Error sending anomalies to Kafka: {e}")

    async def start_async_components(self):
        """Start async components like Kafka consumers."""
        self.logger.info("Starting async components...")

        # Start Kafka consumer
        await self.kafka_consumer.start()

        # Start background tasks
        asyncio.create_task(self._periodic_model_training())
        asyncio.create_task(self._periodic_health_check())

        self.logger.info("Async components started")

    async def _periodic_model_training(self):
        """Periodically retrain models with new data."""
        while self.running:
            try:
                await asyncio.sleep(3600)  # Every hour

                self.logger.info("Starting periodic model training")
                await self.anomaly_detector.retrain_models()
                await self.resource_predictor.retrain_models()
                self.logger.info("Periodic model training completed")

            except Exception as e:
                self.logger.error(f"Error in periodic model training: {e}")

    async def _periodic_health_check(self):
        """Periodically check the health of ML components."""
        while self.running:
            try:
                await asyncio.sleep(300)  # Every 5 minutes

                # Check component health
                detector_health = self.anomaly_detector.get_health_status()
                predictor_health = self.resource_predictor.get_health_status()

                if not detector_health["healthy"] or not predictor_health["healthy"]:
                    self.logger.warning("ML components health check failed")

            except Exception as e:
                self.logger.error(f"Error in health check: {e}")

    def start(self):
        """Start the ML service."""
        self.logger.info("Starting Sentinel Mesh ML Service")
        self.running = True

        # Setup signal handlers
        signal.signal(signal.SIGTERM, self._signal_handler)
        signal.signal(signal.SIGINT, self._signal_handler)

        # Start async components in background
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        loop.create_task(self.start_async_components())

        # Start Flask app
        try:
            host = self.config.get("server", {}).get("host", "0.0.0.0")
            port = self.config.get("server", {}).get("port", 5000)
            debug = self.config.get("server", {}).get("debug", False)

            self.logger.info(f"Starting ML service on {host}:{port}")

            # Run Flask with asyncio support
            self.app.run(host=host, port=port, debug=debug, threaded=True)

        except Exception as e:
            self.logger.error(f"Error starting ML service: {e}")
            self.stop()

    def stop(self):
        """Stop the ML service."""
        self.logger.info("Stopping Sentinel Mesh ML Service")
        self.running = False

        # Stop async components
        if hasattr(self, "kafka_consumer"):
            asyncio.create_task(self.kafka_consumer.stop())

    def _signal_handler(self, signum, frame):
        """Handle shutdown signals."""
        self.logger.info(f"Received signal {signum}, shutting down...")
        self.stop()
        sys.exit(0)


def main():
    """Main entry point for the ML service."""
    try:
        # Load configuration
        config_path = os.getenv("CONFIG_PATH", "config/ml_config.yaml")
        config = load_config(config_path)

        # Setup logging
        setup_logging(config.get("logging", {}))

        logger = logging.getLogger(__name__)
        logger.info("Initializing Sentinel Mesh ML Service")

        # Create and start service
        service = MLService(config)
        service.start()

    except KeyboardInterrupt:
        logger.info("Service interrupted by user")
    except Exception as e:
        logger.error(f"Fatal error in ML service: {e}", exc_info=True)
        sys.exit(1)


if __name__ == "__main__":
    main()
