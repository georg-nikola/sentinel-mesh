# Sentinel Mesh
> A personal hobby project exploring distributed system monitoring with predictive security intelligence

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![Vue.js](https://img.shields.io/badge/vue.js-3.x-green.svg)](https://vuejs.org)
[![Python](https://img.shields.io/badge/python-3.9+-yellow.svg)](https://python.org)
[![Kubernetes](https://img.shields.io/badge/kubernetes-1.25+-blue.svg)](https://kubernetes.io)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Documentation](https://github.com/georg-nikola/sentinel-mesh/workflows/Documentation%20Check/badge.svg)](https://github.com/georg-nikola/sentinel-mesh/actions/workflows/docs-drift.yml)

## Overview

Sentinel Mesh is a cloud-native distributed system monitoring platform that combines real-time observability with machine learning-powered security intelligence. Built for Kubernetes environments, it provides comprehensive monitoring, anomaly detection, and automated incident response capabilities.

> **Personal Project**: This is a hobby project developed for learning and experimentation with Kubernetes, microservices, Vue.js, and ML-powered observability. While functional, it's not intended for production use.

## 🚀 Key Features

### Real-Time Monitoring
- **Service Mesh Observability**: Deep integration with Istio/Linkerd for traffic analysis
- **Multi-Cluster Support**: Cross-cluster visibility and management
- **High-Performance Data Collection**: Go-based collectors with minimal overhead
- **Real-Time Streaming**: Kafka-based data pipeline for sub-second latency

### ML-Powered Intelligence
- **Anomaly Detection**: Behavioral analysis using TensorFlow models
- **Predictive Scaling**: ML-based auto-scaling recommendations
- **Security Correlation**: Event correlation for threat hunting
- **Pattern Recognition**: Automated baseline learning and drift detection

### Operational Excellence
- **SLO/SLI Tracking**: Custom service level objectives monitoring
- **Incident Response**: Automated escalation and context-aware alerts
- **Performance Optimization**: Resource usage optimization recommendations
- **Audit & Compliance**: Comprehensive security audit logging

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Data Sources  │    │  Processing     │    │   Analytics     │
│                 │    │                 │    │                 │
│  • Kubernetes   │    │  • Collector    │    │  • ML Engine    │
│  • Service Mesh │────│  • Processor    │────│  • Analyzer     │
│  • Applications │    │  • Kafka        │    │  • Predictor    │
│  • Infrastructure │   │  • InfluxDB     │    │  • Correlator   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   Presentation  │
                       │                 │
                       │  • Vue.js UI    │
                       │  • REST API     │
                       │  • WebSocket    │
                       │  • Alerting     │
                       └─────────────────┘
```

### Core Components

#### Data Collection Layer
- **Collector Service**: High-performance Go service for metrics and logs collection
- **Kubernetes Integration**: Native K8s API integration with service discovery
- **Custom Metrics**: Extensible metric collection framework

#### Processing Layer
- **Stream Processor**: Real-time data processing with Apache Kafka
- **Data Storage**: InfluxDB for time-series, Elasticsearch for logs
- **Event Bus**: Pub/sub architecture for loose coupling

#### Analytics Engine
- **ML Pipeline**: Python-based machine learning for anomaly detection
- **Security Intelligence**: Threat detection and behavioral analysis
- **Predictive Analytics**: Resource forecasting and capacity planning

#### Presentation Layer
- **Vue.js Dashboard**: Modern, responsive monitoring interface
- **REST API**: Comprehensive API for external integrations
- **Real-time Updates**: WebSocket connections for live data

## 🛠️ Technology Stack

### Backend Services
- **Go 1.21+**: High-performance microservices
- **Apache Kafka**: Real-time streaming platform
- **InfluxDB**: Time-series database
- **Elasticsearch**: Log storage and search
- **Redis**: Caching and session storage

### Machine Learning
- **TensorFlow**: Deep learning framework
- **Python 3.9+**: ML pipeline and analytics
- **NumPy/Pandas**: Data processing libraries
- **Scikit-learn**: Classical ML algorithms

### Frontend
- **Vue.js 3**: Progressive web framework
- **TypeScript**: Type-safe JavaScript
- **Chart.js**: Data visualization
- **WebSocket**: Real-time communications

### Infrastructure
- **Kubernetes 1.25+**: Container orchestration
- **Helm 3**: Package management
- **Docker**: Containerization
- **Prometheus**: Metrics collection
- **Jaeger**: Distributed tracing

## 🚦 Quick Start

### Prerequisites
- Kubernetes cluster (1.25+)
- Helm 3.0+
- kubectl configured
- Docker (for local development)

### Installation

#### 1. Deploy with Helm
```bash docs-drift:skip
# Install from local chart
helm install sentinel-mesh ./deployments/helm/sentinel-mesh \
  --namespace sentinel-mesh \
  --create-namespace

# Install with custom values
helm install sentinel-mesh ./deployments/helm/sentinel-mesh \
  --namespace sentinel-mesh \
  --create-namespace \
  --values values-production.yaml
```

#### 2. Port Forward for Local Access
```bash docs-drift:skip
kubectl port-forward -n sentinel-mesh svc/sentinel-mesh-api 8080:8080
kubectl port-forward -n sentinel-mesh svc/sentinel-mesh-ui 3000:80
```

#### 3. Access the Dashboard
Open your browser to `http://localhost:3000`

### Configuration

#### Basic Configuration
```yaml docs-drift:skip
# values.yaml
global:
  clusterName: "production"
  environment: "prod"

collector:
  replicas: 3
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 500m
      memory: 512Mi

ml:
  enabled: true
  models:
    anomaly_detection: true
    predictive_scaling: true
```

## 📊 Monitoring Capabilities

### Infrastructure Monitoring
- **Node Metrics**: CPU, memory, disk, network usage
- **Pod Monitoring**: Resource utilization and health status
- **Cluster Overview**: Multi-cluster resource visualization
- **Network Traffic**: Service-to-service communication analysis

### Application Performance
- **Request Latency**: P50, P90, P95, P99 percentiles
- **Error Rates**: 4xx, 5xx error tracking
- **Throughput**: Requests per second across services
- **Dependency Mapping**: Service dependency visualization

### Security Intelligence
- **Anomaly Detection**: Behavioral baseline and deviation alerts
- **Threat Correlation**: Security event correlation across services
- **Access Patterns**: Unusual access pattern detection
- **Compliance Monitoring**: Policy violation tracking

### Business Metrics
- **SLO Tracking**: Custom service level objective monitoring
- **Cost Analysis**: Resource cost tracking and optimization
- **Capacity Planning**: Growth trend analysis and forecasting
- **Performance Insights**: Automated performance recommendations

## 🔧 Development

### Local Development Setup

#### 1. Clone Repository
```bash docs-drift:skip
git clone https://github.com/georg-nikola/sentinel-mesh.git
cd sentinel-mesh
```

#### 2. Start Development Environment
```bash docs-drift:skip
# Start local dependencies
make dev-up

# Run backend services
make run-collector
make run-processor
make run-api

# Start frontend development server
cd web
npm install
npm run dev
```

#### 3. Run Tests
```bash docs-drift:skip
# Run all tests
make test

# Run specific test suites
make test-unit
make test-integration
make test-e2e
```

### Project Structure
```
sentinel-mesh/
├── cmd/                    # Application entry points
│   ├── collector/         # Data collection service
│   ├── processor/         # Stream processing service
│   ├── analyzer/          # ML analysis service
│   ├── alerting/          # Alert management service
│   └── api/               # REST API service
├── pkg/                   # Shared Go packages
│   ├── config/           # Configuration management
│   ├── database/         # Database clients
│   ├── metrics/          # Metrics collection
│   ├── security/         # Security utilities
│   └── utils/            # Common utilities
├── internal/             # Internal Go packages
│   ├── handlers/         # HTTP handlers
│   ├── services/         # Business logic
│   └── models/           # Data models
├── deployments/          # Deployment configurations
│   ├── helm/             # Helm charts
│   ├── kubernetes/       # K8s manifests
│   └── docker/           # Docker configurations
├── configs/              # Configuration files
├── web/                  # Vue.js frontend
├── ml/                   # Python ML components
├── docs/                 # Documentation
└── scripts/              # Build and utility scripts
```

## 🔐 Security

### Authentication & Authorization
- **RBAC Integration**: Kubernetes Role-Based Access Control
- **JWT Tokens**: Secure API authentication
- **Multi-tenancy**: Namespace-based isolation
- **Audit Logging**: Comprehensive access audit trails

### Data Security
- **Encryption at Rest**: Database encryption
- **TLS Everywhere**: End-to-end encryption
- **Secret Management**: Kubernetes secrets integration
- **Data Retention**: Configurable data lifecycle policies

### Network Security
- **Network Policies**: Kubernetes network segmentation
- **Service Mesh**: mTLS communication
- **Egress Control**: Outbound traffic filtering
- **Ingress Protection**: WAF integration support

## 📈 Performance

### Scalability
- **Horizontal Scaling**: Auto-scaling based on metrics
- **Resource Optimization**: Dynamic resource allocation
- **Load Balancing**: Intelligent traffic distribution
- **Cache Strategy**: Multi-level caching architecture

### High Availability
- **Multi-AZ Deployment**: Cross-availability zone setup
- **Health Checks**: Comprehensive health monitoring
- **Circuit Breakers**: Fault tolerance patterns
- **Backup & Recovery**: Automated backup strategies

## 🤝 Contributing

Contributions are welcome! This is a learning project, so feel free to experiment and propose improvements.

### How to Contribute
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Ensure tests pass (`./tests/e2e/run-all.sh`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to your fork (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Style
- **Go**: Use `gofmt` and follow standard Go conventions
- **Python**: Follow PEP 8, format with Black
- **TypeScript/Vue**: Follow ESLint rules (see `.eslintrc`)
- **Commits**: Use conventional commit format when possible (e.g., `feat:`, `fix:`, `docs:`)

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📚 Documentation

- [E2E Test Plan](docs/E2E_TEST_PLAN.md) - Comprehensive testing strategy
- [E2E Test Results](docs/E2E_TEST_RESULTS.md) - Latest test execution results
- [CLAUDE.md](CLAUDE.md) - Development guide for Claude Code
- [CHANGELOG.md](CHANGELOG.md) - Release history and changes

## 🐛 Issues & Contributions

This is a personal hobby project, but contributions are welcome!

- [Report Issues](https://github.com/georg-nikola/sentinel-mesh/issues)
- [View Source](https://github.com/georg-nikola/sentinel-mesh)
- [Latest Release](https://github.com/georg-nikola/sentinel-mesh/releases)

Feel free to fork, experiment, and submit pull requests!

---

**Sentinel Mesh** - *A personal project exploring Kubernetes monitoring and ML-powered observability*

> **Note**: This is a hobby project under active development. Not recommended for production use.