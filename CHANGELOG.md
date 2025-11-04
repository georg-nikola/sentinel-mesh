# Changelog

All notable changes to Sentinel Mesh will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2025-10-26

### Added

#### Core Infrastructure
- Initial project setup with distributed system monitoring architecture
- Go-based microservices foundation with support for Go 1.21+
- Comprehensive CI/CD pipeline with GitHub Actions
- Docker multi-stage builds for all services
- Kubernetes deployment manifests and Helm charts
- Security scanning with Trivy and dependency checks

#### Backend Services
- **API Service**: REST API server with health and metrics endpoints
- **Collector Service**: Kubernetes metrics collection service
- **Processor Service**: Data processing service with health checks
- **Analyzer Service**: ML analysis service with monitoring endpoints
- **Alerting Service**: Alert management and notification service
- **ML Service**: Python-based machine learning service with Flask
  - Anomaly detection API (`/api/v1/anomalies`)
  - Resource prediction API (`/api/v1/predictions`)
  - Prometheus metrics endpoint
  - Health and readiness probes

#### Frontend Dashboard
- Modern Vue.js 3 + TypeScript web application
- Real-time monitoring dashboard with live data updates
- Dark mode support with system preference detection
- Responsive design with Tailwind CSS
- Multiple monitoring views:
  - **Dashboard**: Overview with CPU and memory usage charts
  - **Metrics**: Detailed metrics visualization (Request Rate, CPU, Memory, Response Time)
  - **Security**: Security alerts and anomaly detection integration
  - **Nodes**: Kubernetes node monitoring
  - **Services**: Service health and status tracking
- Real-time notifications system with severity levels
- Chart.js integration for data visualization
- WebSocket-based live updates
- ML anomaly integration with auto-refresh

#### Data Visualization
- CPU usage trending with Chart.js
- Memory usage monitoring with real-time updates
- Request rate tracking
- Response time analysis
- Auto-scaling charts (keeps last 20 data points)
- Color-coded severity indicators

#### Monitoring & Observability
- Health check endpoints for all services
- Readiness probes for Kubernetes deployments
- Prometheus metrics exposition
- Service version reporting
- Graceful shutdown handling for all Go services

#### Security Features
- Security alert dashboard with severity levels (info, warning, critical)
- ML-powered anomaly detection integration
- Dark mode optimized notification visibility
- Auto-removal of stale alerts (older than 1 hour)
- Manual alert dismissal functionality
- Real-time security event synchronization

#### Development Experience
- Comprehensive test suite for all components:
  - Go service unit tests
  - Python ML service tests
  - Frontend TypeScript/Vue tests
  - Helm chart validation tests
- Automated code formatting (Black for Python, gofmt for Go)
- Docker Compose setup for local development
- Makefile for common development tasks
- Environment-based configuration

### Fixed
- Go module dependencies resolution and compilation issues
- Python dependency compatibility with Flask and CORS
- CI pipeline reliability and test isolation
- Docker build configurations and multi-stage optimization
- TypeScript type safety improvements
- Chart reactivity issues with Vue 3 shallow refs
- Dark mode visibility for notifications and alerts
- Python code formatting compliance with Black/PEP 8

### Infrastructure
- GitHub Actions CI/CD pipeline with parallel jobs:
  - Go service testing
  - ML service testing (Python)
  - Frontend testing (TypeScript/Vue)
  - Helm chart linting and validation
  - Security scanning with Trivy
  - Docker image builds (7 services)
- Automated dependency caching for faster builds
- SARIF support for security scanning (optional)
- Non-blocking tests for gradual improvements

### Documentation
- Comprehensive README with architecture diagrams
- Quick start guide with Helm installation
- Technology stack documentation
- Project structure overview
- Development setup instructions
- Configuration examples

### Known Limitations
- GitHub Container Registry push requires additional permissions configuration
- Some services are stub implementations (planned for future releases)
- ML models are currently mock implementations for testing
- Advanced security features planned for future releases

## [0.2.0] - 2025-11-01

### Added

#### Production Deployment
- Production deployment manifests for Talos Kubernetes cluster
- Traefik IngressRoute configuration for external access
- Basic authentication middleware for dashboard protection
- Auth setup script (`auth-setup.sh`) for creating htpasswd credentials
- Production deployment README with step-by-step instructions
- Namespace configuration for sentinel-mesh in production

#### Documentation
- Comprehensive staging deployment guide using OrbStack
- Production deployment guide with Traefik and Cloudflare Tunnel integration
- E2E testing workflow for staging environments
- Production verification checklists
- Cloudflare Tunnel configuration examples
- DNS setup with Terraform examples
- Deployment troubleshooting guide
- Updated CLAUDE.md with detailed deployment procedures

#### CI/CD Improvements
- Auto-merge for GitHub Actions dependency updates
- Dependabot dependency update automation

### Changed
- Updated CLAUDE.md deployment section with two-stage process (staging → production)
- Enhanced documentation structure for better deployment workflow visibility
- Improved production security with basic-auth protection on all dashboards

### Security
- Basic authentication protection for production dashboards
- Secure credential management with Kubernetes secrets
- Integration with existing Cloudflare security features (DDoS, WAF-ready)

## [Unreleased]

### Planned Features
- Real Kafka integration for event streaming
- InfluxDB integration for time-series data
- Elasticsearch integration for log aggregation
- Advanced ML models for anomaly detection
- Multi-cluster support
- Service mesh integration (Istio/Linkerd)
- SLO/SLI tracking
- Automated incident response
- Cost analysis and optimization

---

[0.2.0]: https://github.com/georg-nikola/sentinel-mesh/releases/tag/v0.2.0
[0.1.0]: https://github.com/georg-nikola/sentinel-mesh/releases/tag/v0.1.0
