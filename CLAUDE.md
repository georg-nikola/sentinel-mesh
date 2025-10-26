# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sentinel Mesh is a cloud-native distributed system monitoring platform built for Kubernetes environments. It combines real-time observability with machine learning-powered security intelligence to provide comprehensive monitoring, anomaly detection, and automated incident response capabilities.

The project uses a microservices architecture with:
- **Go 1.21+** for high-performance backend services
- **Python 3.9+** for ML/analytics services
- **Vue.js 3 + TypeScript** for the frontend dashboard
- **Kubernetes** for orchestration and deployment

## Development Commands

### Backend Services (Go)

```bash
# Install Go dependencies
go mod download
go mod tidy

# Build all services
go build -o bin/api cmd/api/main.go
go build -o bin/collector cmd/collector/main.go
go build -o bin/processor cmd/processor/main.go
go build -o bin/analyzer cmd/analyzer/main.go
go build -o bin/alerting cmd/alerting/main.go

# Run tests
go test ./...
go test -v -race ./...  # With race detection

# Format code
gofmt -w .
go vet ./...

# Run a specific service
go run cmd/api/main.go
PORT=8080 go run cmd/collector/main.go
```

### ML Service (Python)

```bash
# Set up virtual environment
python3 -m venv venv
source venv/bin/activate  # On Windows: venv\Scripts\activate

# Install dependencies
pip install -r ml/requirements.txt

# Run ML service
cd ml
python simple_main.py
# Service runs on http://0.0.0.0:8000

# Format code (required for CI)
black ml/
black --check ml/  # Check without modifying

# Run tests
pytest ml/
pytest ml/ -v  # Verbose output
```

### Frontend (Vue.js)

```bash
# Install dependencies
cd web
npm install

# Development server with hot reload
npm run dev
# Access at http://localhost:5173

# Type checking
npm run type-check

# Linting
npm run lint

# Build for production
npm run build
# Output in web/dist/

# Preview production build
npm run preview
```

### Docker Development

```bash
# Build all service images
docker build -t sentinel-mesh/api -f Dockerfile --target api .
docker build -t sentinel-mesh/collector -f Dockerfile --target collector .
docker build -t sentinel-mesh/ml-service -f ml/Dockerfile ml/

# Run with Docker Compose (local development)
docker-compose up
docker-compose up -d  # Detached mode
docker-compose down   # Stop and remove containers

# View logs
docker-compose logs -f api
docker-compose logs -f ml-service
```

### Kubernetes Deployment

```bash
# Deploy to local Kubernetes (OrbStack, minikube, etc.)
kubectl apply -f deployments/kubernetes/

# Check deployment status
kubectl get pods -n sentinel-mesh
kubectl get services -n sentinel-mesh

# View service logs
kubectl logs -f -n sentinel-mesh deployment/api
kubectl logs -f -n sentinel-mesh deployment/ml-service

# Port forward for local access
kubectl port-forward -n sentinel-mesh svc/api 8080:8080
kubectl port-forward -n sentinel-mesh svc/frontend 3000:80

# Delete deployment
kubectl delete -f deployments/kubernetes/
```

### Helm Charts

```bash
# Lint Helm charts
helm lint deployments/helm/sentinel-mesh

# Template and view generated manifests
helm template sentinel-mesh deployments/helm/sentinel-mesh

# Install with Helm
helm install sentinel-mesh deployments/helm/sentinel-mesh \
  --namespace sentinel-mesh \
  --create-namespace

# Upgrade release
helm upgrade sentinel-mesh deployments/helm/sentinel-mesh

# Uninstall
helm uninstall sentinel-mesh -n sentinel-mesh
```

## Architecture

### Service Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                            │
│                   Vue.js 3 + TypeScript                     │
│              (Dashboard, Metrics, Security UI)              │
└─────────────────────────────────────────────────────────────┘
                            │ HTTP/WebSocket
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                         API Service                         │
│                         (Go/Gin)                            │
│              Health, Metrics, Data Aggregation              │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  Collector   │   │  Processor   │   │   Analyzer   │
│   (Go)       │   │   (Go)       │   │   (Go)       │
│ K8s Metrics  │   │ Data Proc.   │   │ Analysis     │
└──────────────┘   └──────────────┘   └──────────────┘
                            │
                            ▼
                   ┌──────────────┐   ┌──────────────┐
                   │ ML Service   │   │  Alerting    │
                   │  (Python)    │   │   (Go)       │
                   │ Anomaly Det. │   │ Alerts Mgmt  │
                   └──────────────┘   └──────────────┘
```

### Data Flow

1. **Collection**: Collector service gathers metrics from Kubernetes API
2. **Processing**: Processor service normalizes and enriches data
3. **Analysis**: Analyzer service processes data for trends
4. **ML Intelligence**: ML service detects anomalies and predicts resource needs
5. **Alerting**: Alerting service manages notifications
6. **Presentation**: API service serves data to Vue.js frontend
7. **Visualization**: Dashboard displays real-time metrics with Chart.js

### Key Components

#### Backend Services (Go)

All Go services follow a consistent pattern:
- HTTP server with health (`/health`), readiness (`/ready`), and metrics (`/metrics`) endpoints
- Graceful shutdown with signal handling (SIGINT, SIGTERM)
- Configurable port via `PORT` environment variable (default: 8080)
- Version info injected at build time (`version`, `buildTime`, `gitCommit`)
- 15s read/write timeout, 60s idle timeout

**API Service** (`cmd/api/main.go`):
- Primary REST API endpoint
- Aggregates data from other services
- Serves frontend requests

**Collector Service** (`cmd/collector/main.go`):
- Kubernetes metrics collection
- Resource usage monitoring
- Service discovery

**Processor Service** (`cmd/processor/main.go`):
- Data normalization
- Event processing
- Stream handling

**Analyzer Service** (`cmd/analyzer/main.go`):
- Trend analysis
- Pattern recognition
- Statistical processing

**Alerting Service** (`cmd/alerting/main.go`):
- Alert management
- Notification routing
- Escalation logic

#### ML Service (Python)

**Location**: `ml/simple_main.py`

Flask-based service providing:
- `/health` - Health check
- `/ready` - Readiness probe
- `/metrics` - Prometheus metrics
- `/api/v1/anomalies` - Detected anomalies
- `/api/v1/predictions` - Resource predictions

Currently uses mock data for testing; real ML models planned for future releases.

#### Frontend Dashboard (Vue.js)

**Location**: `web/src/`

**Views**:
- `Dashboard.vue` - Overview with CPU/Memory charts (2 charts)
- `Metrics.vue` - Detailed metrics (Request Rate, CPU, Memory, Response Time - 4 charts)
- `Security.vue` - Security alerts and ML anomaly integration
- `Nodes.vue` - Kubernetes node monitoring
- `Services.vue` - Service health tracking

**Key Features**:
- Real-time updates using `setInterval` (30s refresh)
- Chart.js with reactive updates via `shallowRef`
- Dark mode with Tailwind CSS `dark:` variants
- Responsive design for mobile/tablet/desktop
- Live notifications system with severity levels
- Auto-removal of stale alerts (>1 hour old)

**State Management**:
- Vue 3 Composition API with `ref` and `shallowRef`
- No Vuex/Pinia - local component state only
- API calls with axios
- WebSocket connections for real-time updates

**Styling**:
- Tailwind CSS for utility-first styling
- Dark mode via `class` strategy (not media query)
- Custom color schemes for severity indicators

## Project Structure

```
sentinel-mesh/
├── cmd/                          # Service entry points
│   ├── api/main.go              # API service
│   ├── collector/main.go        # Metrics collector
│   ├── processor/main.go        # Data processor
│   ├── analyzer/main.go         # Analysis service
│   └── alerting/main.go         # Alert manager
├── pkg/                          # Shared Go libraries (future)
├── internal/                     # Private Go packages (future)
├── ml/                           # Python ML service
│   ├── simple_main.py           # Flask ML API
│   ├── requirements.txt         # Python dependencies
│   └── Dockerfile               # ML service container
├── web/                          # Vue.js frontend
│   ├── src/
│   │   ├── views/               # Page components
│   │   ├── components/          # Reusable components
│   │   ├── config/              # Configuration
│   │   └── App.vue              # Root component
│   ├── package.json
│   └── vite.config.ts
├── deployments/                  # Deployment configs
│   ├── helm/                    # Helm charts
│   │   └── sentinel-mesh/
│   └── kubernetes/              # Raw K8s manifests
├── .github/                      # GitHub configs
│   ├── workflows/               # CI/CD pipelines
│   │   └── ci.yml              # Main CI workflow
│   ├── dependabot.yml          # Dependency updates
│   ├── qodo.toml               # Qodo PR review config
│   └── QODO_SETUP.md           # Qodo setup guide
├── Dockerfile                    # Multi-stage Go services
├── docker-compose.yml           # Local development
├── go.mod / go.sum              # Go dependencies
├── CHANGELOG.md                 # Release notes
├── README.md                    # Project documentation
└── CLAUDE.md                    # This file
```

## Development Workflow

### Making Changes

1. **Create Feature Branch** (if not admin pushing to main):
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**: Edit code following project conventions

3. **Test Locally**:
   ```bash
   # Go services
   go test ./...

   # ML service
   black --check ml/
   pytest ml/

   # Frontend
   cd web && npm run type-check && npm run lint
   ```

4. **Commit Changes**:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   # Use conventional commits: feat, fix, docs, refactor, test, chore
   ```

5. **Push and Create PR**:
   ```bash
   git push origin feature/your-feature-name
   gh pr create --title "Feature: Your feature" --body "Description..."
   ```

### Conventional Commits

Use these prefixes:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `test:` - Adding/updating tests
- `chore:` - Maintenance tasks
- `perf:` - Performance improvements
- `style:` - Code style changes (formatting)
- `ci:` - CI/CD changes
- `deps:` - Dependency updates

### Code Style Guidelines

**Go**:
- Use `gofmt` for formatting (enforced in CI)
- Follow effective Go patterns
- Handle all errors explicitly
- Use meaningful variable names
- Add comments for exported functions

**Python**:
- Use Black formatter (enforced in CI)
- Follow PEP 8 style guide
- Type hints encouraged but not required
- Docstrings for all functions

**TypeScript/Vue**:
- Use ESLint configuration (in web/.eslintrc.cjs)
- Composition API preferred over Options API
- Props and emits should be typed
- Use `shallowRef` for Chart.js data objects

## Testing Strategy

### Go Services

```bash
# Unit tests
go test ./...

# With coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# With race detection
go test -race ./...

# Specific package
go test ./cmd/api/...
```

### Python ML Service

```bash
# All tests
pytest ml/

# With coverage
pytest ml/ --cov=ml --cov-report=html

# Specific test file
pytest ml/test_service.py -v
```

### Frontend

```bash
cd web

# Type checking
npm run type-check

# Linting
npm run lint

# Fix lint issues
npm run lint -- --fix

# Component tests (if added)
npm run test
```

### Integration Testing

Currently manual testing via:
1. Deploy to local Kubernetes cluster
2. Port-forward services
3. Access frontend at http://localhost:3000
4. Verify real-time updates and API integration

## CI/CD Pipeline

### GitHub Actions Workflow

**File**: `.github/workflows/ci.yml`

**Jobs**:
1. **Test Go Services** - Run Go tests, build binaries
2. **Test ML Services** - Black formatting check, pytest
3. **Test Frontend** - TypeScript type check, ESLint
4. **Test Helm Charts** - Helm lint validation
5. **Security Scan** - Trivy vulnerability scanning
6. **Build Docker Images** - Multi-stage builds for all services

**Triggers**:
- Push to `main` branch
- Pull requests to `main`
- Manual workflow dispatch

**Status Checks Required** (branch protection):
- Test Go Services
- Test ML Services
- Test Frontend
- Test Helm Charts
- Security Scan

### Docker Build Strategy

Multi-stage Dockerfile for Go services:
```dockerfile
# Stage 1: Build
FROM golang:1.21-alpine AS builder
# ... build process ...

# Stage 2: Runtime (per service)
FROM alpine:latest AS api
COPY --from=builder /app/bin/api /usr/local/bin/api
CMD ["api"]
```

Separate Dockerfile for ML service:
```dockerfile
FROM python:3.11-slim
COPY requirements.txt .
RUN pip install -r requirements.txt
CMD ["python", "simple_main.py"]
```

## Configuration

### Environment Variables

**Go Services**:
- `PORT` - HTTP server port (default: 8080)
- `LOG_LEVEL` - Logging level (default: info)
- `ENVIRONMENT` - Environment name (dev/staging/prod)

**ML Service**:
- `FLASK_ENV` - Flask environment (development/production)
- `ML_MODEL_PATH` - Path to ML models (future)

**Frontend** (build-time):
- `VITE_API_URL` - API service URL
- `VITE_ML_SERVICE_URL` - ML service URL

### API Configuration

**File**: `web/src/config/api.ts`

```typescript
export const API_CONFIG = {
  BASE_URL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  ML_SERVICE_URL: import.meta.env.VITE_ML_SERVICE_URL || 'http://localhost:8000',
  TIMEOUT: 5000, // 5 seconds
}
```

## Deployment

### Local Development (OrbStack)

```bash
# Start OrbStack Kubernetes
orb start

# Deploy services
kubectl apply -f deployments/kubernetes/

# Access frontend
kubectl port-forward -n sentinel-mesh svc/frontend 3000:80
# Open http://localhost:3000
```

### Production Deployment

1. **Build and Push Images**:
   ```bash
   # Configure registry (DockerHub, GHCR, etc.)
   docker build -t your-registry/sentinel-mesh-api:v0.1.0 .
   docker push your-registry/sentinel-mesh-api:v0.1.0
   ```

2. **Deploy with Helm**:
   ```bash
   helm upgrade --install sentinel-mesh deployments/helm/sentinel-mesh \
     --namespace sentinel-mesh \
     --create-namespace \
     --set image.tag=v0.1.0 \
     --set ingress.enabled=true \
     --set ingress.hosts[0].host=sentinel.yourdomain.com
   ```

3. **Verify Deployment**:
   ```bash
   kubectl get pods -n sentinel-mesh
   kubectl get ingress -n sentinel-mesh
   ```

## Troubleshooting

### Charts Not Updating in Frontend

**Issue**: Chart.js charts show initial data but don't update with new values.

**Solution**: Ensure using immutable updates with `shallowRef`:
```typescript
// ❌ Wrong - mutating nested objects
chartData.value.labels.push(newLabel)
chartData.value = { ...chartData.value }

// ✅ Correct - creating new objects with new array references
chartData.value = {
  labels: [...chartData.value.labels, newLabel],
  datasets: [{
    ...chartData.value.datasets[0],
    data: [...chartData.value.datasets[0].data, newValue]
  }]
}
```

### Dark Mode Not Working

**Issue**: Dark mode styles not applying correctly.

**Solution**: Ensure using `dark:` variants in Tailwind classes:
```vue
<!-- ❌ Wrong -->
<div class="bg-white text-black">

<!-- ✅ Correct -->
<div class="bg-white dark:bg-gray-800 text-gray-900 dark:text-white">
```

### Python Formatting Failing in CI

**Issue**: Black formatter reports "would reformat" errors.

**Solution**: Run Black locally before committing:
```bash
black ml/
# Or with Docker:
docker run --rm -v "${PWD}:/code" -w /code pyfound/black:latest black ml/
```

### Go Module Issues

**Issue**: Missing dependencies or version conflicts.

**Solution**:
```bash
go mod tidy
go mod download
go clean -modcache  # If cache is corrupted
```

### Kubernetes Pods CrashLooping

**Issue**: Pods failing to start in K8s.

**Debug**:
```bash
kubectl describe pod <pod-name> -n sentinel-mesh
kubectl logs <pod-name> -n sentinel-mesh
kubectl get events -n sentinel-mesh --sort-by='.lastTimestamp'
```

## Important Notes

### Chart.js Reactivity

When using Chart.js with Vue 3:
- Use `shallowRef` for chart data objects
- Create entirely new objects on each update (immutable pattern)
- Never mutate nested arrays/objects directly
- Keep last 20 data points for performance (trim older data)

### ML Service Integration

Current ML service (`ml/simple_main.py`) provides mock data:
- Anomalies are static examples
- Predictions are placeholder values
- Real ML models planned for v0.2.0

Frontend Security view (`web/src/views/Security.vue`) integrates ML anomalies:
- Fetches from `/api/v1/anomalies` every 30s
- Auto-removes stale alerts (>1 hour)
- Syncs with ML service (removes resolved anomalies)

### Branch Protection

Main branch is protected:
- Requires PR review (1 approval)
- Requires all CI checks to pass
- Admin (georg-nikola) can bypass for urgent fixes
- Linear history enforced (no merge commits without fast-forward)

### Dependency Management

Automated via Dependabot:
- Weekly checks on Monday (Go, npm, Python)
- Tuesday for Docker images
- Wednesday for GitHub Actions
- Grouped minor/patch updates to reduce noise

### Code Review Automation

Qodo Merge configured for PR reviews:
- Auto-reviews on PR creation
- Security and performance checks
- Code suggestions and improvements
- Requires manual installation of GitHub App

## Future Enhancements

**Planned for v0.2.0**:
- Real ML models for anomaly detection
- Kafka integration for event streaming
- InfluxDB for time-series storage
- Elasticsearch for log aggregation
- Multi-cluster support
- Service mesh integration (Istio/Linkerd)
- Advanced alerting rules
- SLO/SLI tracking

## Resources

### Documentation
- [Go Documentation](https://go.dev/doc/)
- [Vue.js 3 Guide](https://vuejs.org/guide/)
- [Chart.js Docs](https://www.chartjs.org/docs/)
- [Kubernetes Docs](https://kubernetes.io/docs/)
- [Helm Docs](https://helm.sh/docs/)

### Internal Docs
- [CHANGELOG.md](CHANGELOG.md) - Release history
- [README.md](README.md) - Project overview
- [.github/QODO_SETUP.md](.github/QODO_SETUP.md) - Qodo configuration

### Repository Links
- Issues: https://github.com/georg-nikola/sentinel-mesh/issues
- Releases: https://github.com/georg-nikola/sentinel-mesh/releases
- Actions: https://github.com/georg-nikola/sentinel-mesh/actions

---

**Sentinel Mesh** - *Cloud-native distributed system monitoring with ML-powered security intelligence*
