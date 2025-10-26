#!/usr/bin/env python3
"""
Simplified Sentinel Mesh ML Service
Simple health endpoint and metrics for testing
"""

import logging
import os
from flask import Flask, jsonify
from flask_cors import CORS

# Setup logging
logging.basicConfig(
    level=logging.INFO, format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)

app = Flask(__name__)
CORS(app)


@app.route("/health")
def health():
    """Health check endpoint"""
    return jsonify({"status": "healthy", "service": "ml-service", "version": "1.0.0"})


@app.route("/ready")
def ready():
    """Readiness check endpoint"""
    return jsonify({"status": "ready"}), 200


@app.route("/metrics")
def metrics():
    """Prometheus metrics endpoint"""
    return """# HELP ml_service_requests_total Total number of requests
# TYPE ml_service_requests_total counter
ml_service_requests_total 42
# HELP ml_service_up Service up status
# TYPE ml_service_up gauge
ml_service_up 1
"""


@app.route("/api/v1/anomalies")
def anomalies():
    """Get detected anomalies"""
    return jsonify(
        {
            "anomalies": [
                {
                    "type": "cpu_spike",
                    "severity": "warning",
                    "resource": "pod/api-7744d86b68-rkgrf",
                    "timestamp": "2025-10-23T10:00:00Z",
                    "value": 95.2,
                },
                {
                    "type": "memory_leak",
                    "severity": "critical",
                    "resource": "pod/collector-76b4b54f77-8zvv4",
                    "timestamp": "2025-10-23T10:15:00Z",
                    "value": 87.5,
                },
            ]
        }
    )


@app.route("/api/v1/predictions")
def predictions():
    """Get resource predictions"""
    return jsonify(
        {
            "predictions": {
                "next_hour": {"cpu_usage": 72.3, "memory_usage": 65.8, "pod_count": 12},
                "next_day": {"cpu_usage": 68.1, "memory_usage": 62.4, "pod_count": 11},
            }
        }
    )


if __name__ == "__main__":
    host = "0.0.0.0"
    port = 8000

    logger.info(f"Starting ML Service on {host}:{port}")
    app.run(host=host, port=port, debug=False)
