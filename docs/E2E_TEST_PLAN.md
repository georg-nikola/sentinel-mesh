# Sentinel Mesh E2E Test Plan

This document outlines all functionalities that need to be tested end-to-end before releases.

## Test Environment Requirements

- Kubernetes cluster (minikube, kind, OrbStack, or similar)
- kubectl configured and connected
- Helm 3.x installed
- Docker installed
- Network access for image pulling

## Functionality Checklist

### 1. Service Health & Availability

#### 1.1 API Service
- [ ] Service starts successfully
- [ ] Health endpoint responds: `GET /health` returns 200
- [ ] Readiness endpoint responds: `GET /ready` returns 200
- [ ] Metrics endpoint responds: `GET /metrics` returns Prometheus format
- [ ] Service version information is correct
- [ ] Service handles graceful shutdown (SIGTERM)
- [ ] Service restarts after crash (K8s liveness probe)

#### 1.2 Collector Service
- [ ] Service starts successfully
- [ ] Health endpoint responds: `GET /health` returns 200
- [ ] Readiness endpoint responds: `GET /ready` returns 200
- [ ] Metrics endpoint responds: `GET /metrics` returns Prometheus format
- [ ] Service collects Kubernetes metrics
- [ ] Service handles graceful shutdown

#### 1.3 Processor Service
- [ ] Service starts successfully
- [ ] Health endpoint responds: `GET /health` returns 200
- [ ] Readiness endpoint responds: `GET /ready` returns 200
- [ ] Metrics endpoint responds: `GET /metrics` returns Prometheus format
- [ ] Service handles graceful shutdown

#### 1.4 Analyzer Service
- [ ] Service starts successfully
- [ ] Health endpoint responds: `GET /health` returns 200
- [ ] Readiness endpoint responds: `GET /ready` returns 200
- [ ] Metrics endpoint responds: `GET /metrics` returns Prometheus format
- [ ] Service handles graceful shutdown

#### 1.5 Alerting Service
- [ ] Service starts successfully
- [ ] Health endpoint responds: `GET /health` returns 200
- [ ] Readiness endpoint responds: `GET /ready` returns 200
- [ ] Metrics endpoint responds: `GET /metrics` returns Prometheus format
- [ ] Service handles graceful shutdown

#### 1.6 ML Service (Python)
- [ ] Service starts successfully
- [ ] Health endpoint responds: `GET /health` returns 200
- [ ] Readiness endpoint responds: `GET /ready` returns 200
- [ ] Metrics endpoint responds: `GET /metrics` returns Prometheus format
- [ ] Anomaly detection endpoint: `GET /api/v1/anomalies` returns valid JSON
- [ ] Predictions endpoint: `GET /api/v1/predictions` returns valid JSON
- [ ] Service handles graceful shutdown

### 2. Frontend Dashboard

#### 2.1 Application Loading
- [ ] Frontend loads without errors
- [ ] No console errors on page load
- [ ] All assets (CSS, JS) load successfully
- [ ] Favicon and meta tags are correct

#### 2.2 Navigation
- [ ] Dashboard route loads: `/`
- [ ] Metrics route loads: `/metrics`
- [ ] Security route loads: `/security`
- [ ] Nodes route loads: `/nodes`
- [ ] Services route loads: `/services`
- [ ] Active route is highlighted in navigation
- [ ] Navigation works via sidebar clicks
- [ ] Direct URL navigation works

#### 2.3 Dashboard View
- [ ] CPU usage chart renders
- [ ] Memory usage chart renders
- [ ] Charts update with real data (not stuck at 0)
- [ ] Charts show last 20 data points
- [ ] Old data points are removed after reaching 20
- [ ] Data refreshes every 30 seconds
- [ ] Chart animations work smoothly

#### 2.4 Metrics View
- [ ] Request Rate chart renders
- [ ] CPU Usage chart renders
- [ ] Memory Usage chart renders
- [ ] Response Time chart renders
- [ ] All 4 charts update with real data
- [ ] Charts auto-refresh every 30 seconds
- [ ] Data points are trimmed to last 20

#### 2.5 Security View
- [ ] Security score displays
- [ ] Alert list renders
- [ ] Alert severity badges display correctly (critical/warning/info)
- [ ] ML anomalies are fetched and displayed
- [ ] Dismiss button removes alerts
- [ ] Stale alerts (>1 hour) are auto-removed
- [ ] Resolved ML anomalies are removed from list
- [ ] New ML anomalies appear automatically
- [ ] Timestamps show relative time ("5 minutes ago")

#### 2.6 Nodes View
- [ ] Node list displays
- [ ] Node status indicators work
- [ ] Node metrics are shown

#### 2.7 Services View
- [ ] Service list displays
- [ ] Service status indicators work
- [ ] Service health is accurate

#### 2.8 Notifications System
- [ ] Notification bell icon shows count
- [ ] Clicking bell opens notification panel
- [ ] Notifications display with correct severity colors
- [ ] Dismiss button removes individual notifications
- [ ] Notifications are visible in both light and dark modes
- [ ] Panel closes when clicking outside

#### 2.9 Dark Mode
- [ ] Dark mode toggle works
- [ ] Theme preference persists on page reload
- [ ] All components are visible in dark mode
- [ ] Charts are readable in dark mode
- [ ] Notifications have good contrast in dark mode
- [ ] Alert badges are visible in dark mode

### 3. API Integration

#### 3.1 Frontend → API Communication
- [ ] Frontend can reach API service
- [ ] CORS headers are configured correctly
- [ ] API returns valid JSON responses
- [ ] Error handling works for failed requests
- [ ] Timeout handling works (5s timeout)

#### 3.2 Frontend → ML Service Communication
- [ ] Frontend can reach ML service
- [ ] Anomalies endpoint returns data
- [ ] Predictions endpoint returns data
- [ ] Connection timeout handled gracefully
- [ ] Failed requests don't break UI

#### 3.3 Service-to-Service Communication
- [ ] API can communicate with Collector
- [ ] API can communicate with Processor
- [ ] API can communicate with Analyzer
- [ ] API can communicate with Alerting
- [ ] Services can reach ML service

### 4. Data Flow & Processing

#### 4.1 Metrics Collection
- [ ] CPU metrics are collected
- [ ] Memory metrics are collected
- [ ] Request rate metrics are collected
- [ ] Response time metrics are collected
- [ ] Metrics are accurate (within reasonable range)
- [ ] Metrics update in real-time

#### 4.2 Anomaly Detection
- [ ] ML service detects CPU spikes
- [ ] ML service detects memory leaks
- [ ] Anomalies are categorized by severity
- [ ] Anomalies include timestamp
- [ ] Anomalies include resource information
- [ ] Anomalies appear in Security view

#### 4.3 Predictions
- [ ] Next hour predictions are generated
- [ ] Next day predictions are generated
- [ ] Predictions include CPU, memory, and pod count
- [ ] Predictions are reasonable values

### 5. Kubernetes Integration

#### 5.1 Deployment
- [ ] All pods start successfully
- [ ] Pods reach Running state
- [ ] No CrashLoopBackOff errors
- [ ] Resource limits are respected
- [ ] Resource requests are met

#### 5.2 Services
- [ ] All Kubernetes services are created
- [ ] Services have correct selectors
- [ ] Services expose correct ports
- [ ] ClusterIP services are accessible internally
- [ ] NodePort/LoadBalancer services accessible externally

#### 5.3 Health Probes
- [ ] Liveness probes pass for all pods
- [ ] Readiness probes pass for all pods
- [ ] Failed probes trigger pod restart
- [ ] Startup probes (if configured) work

#### 5.4 ConfigMaps & Secrets
- [ ] ConfigMaps are mounted correctly
- [ ] Environment variables are set
- [ ] Secrets are mounted securely
- [ ] Configuration changes trigger pod updates

#### 5.5 Persistence (if applicable)
- [ ] Persistent volumes are created
- [ ] Data persists across pod restarts
- [ ] Volume mounts are correct

### 6. Helm Chart Functionality

#### 6.1 Installation
- [ ] `helm install` completes successfully
- [ ] All resources are created
- [ ] No validation errors
- [ ] Release is marked as deployed

#### 6.2 Upgrade
- [ ] `helm upgrade` works without errors
- [ ] Rolling updates happen gracefully
- [ ] No downtime during upgrade
- [ ] Old pods are terminated properly

#### 6.3 Rollback
- [ ] `helm rollback` works correctly
- [ ] Previous version is restored
- [ ] Services remain available

#### 6.4 Uninstall
- [ ] `helm uninstall` removes all resources
- [ ] No orphaned resources remain
- [ ] Namespace can be deleted cleanly

#### 6.5 Values Customization
- [ ] Custom replica counts work
- [ ] Custom resource limits work
- [ ] Custom environment variables work
- [ ] Custom image tags work
- [ ] Ingress configuration works

### 7. Security

#### 7.1 Container Security
- [ ] No critical vulnerabilities in images (Trivy scan)
- [ ] Images run as non-root user (where possible)
- [ ] Read-only root filesystem (where applicable)
- [ ] No hardcoded secrets in images

#### 7.2 Network Security
- [ ] Services are only accessible where needed
- [ ] No unnecessary port exposure
- [ ] TLS/SSL configured (if enabled)

#### 7.3 RBAC & Permissions
- [ ] Service accounts have minimal permissions
- [ ] RBAC roles are scoped appropriately
- [ ] No cluster-admin permissions used unnecessarily

### 8. Performance & Scalability

#### 8.1 Resource Usage
- [ ] Services stay within memory limits
- [ ] Services stay within CPU limits
- [ ] No memory leaks observed
- [ ] No CPU throttling under normal load

#### 8.2 Horizontal Scaling
- [ ] Multiple replicas can run simultaneously
- [ ] Load is distributed across replicas
- [ ] Services handle replica failures gracefully

#### 8.3 Response Times
- [ ] API endpoints respond within acceptable time (<1s)
- [ ] Frontend loads within acceptable time (<3s)
- [ ] Charts update without lag
- [ ] No UI freezing during data updates

### 9. Error Handling & Recovery

#### 9.1 Service Failures
- [ ] Services restart after crash
- [ ] Dependent services handle upstream failures
- [ ] Circuit breakers prevent cascade failures
- [ ] Error messages are logged properly

#### 9.2 Network Failures
- [ ] Services retry failed requests
- [ ] Timeout handling works correctly
- [ ] UI shows appropriate error messages
- [ ] Services recover when network is restored

#### 9.3 Invalid Data
- [ ] Services validate input data
- [ ] Invalid requests return appropriate errors
- [ ] Services don't crash on bad data
- [ ] Error responses follow standard format

### 10. Monitoring & Observability

#### 10.1 Metrics Export
- [ ] Prometheus metrics are exposed
- [ ] Metrics format is valid
- [ ] Metrics include service-specific data
- [ ] Metrics can be scraped by Prometheus

#### 10.2 Logging
- [ ] All services log to stdout/stderr
- [ ] Log format is consistent
- [ ] Log levels are appropriate
- [ ] Sensitive data is not logged

#### 10.3 Tracing (future)
- [ ] Distributed tracing headers are propagated
- [ ] Traces can be collected
- [ ] Trace data is meaningful

## E2E Test Scenarios

### Scenario 1: Fresh Installation
1. Install Sentinel Mesh via Helm
2. Verify all pods are running
3. Port-forward frontend service
4. Access dashboard in browser
5. Verify all views load correctly
6. Verify charts display data
7. Verify ML anomalies appear in Security view
8. Toggle dark mode and verify visibility
9. Dismiss alerts and verify they disappear
10. Wait 30s and verify data refreshes

**Expected Result**: All functionality works without errors.

### Scenario 2: Service Restart
1. Delete API pod: `kubectl delete pod -l app=api`
2. Wait for pod to restart
3. Verify pod enters Running state
4. Verify frontend continues to work
5. Verify data updates resume

**Expected Result**: Automatic recovery with minimal disruption.

### Scenario 3: High Load
1. Generate high CPU load in cluster
2. Verify Collector detects increased usage
3. Verify charts show spike
4. Verify ML service detects anomaly
5. Verify alert appears in Security view

**Expected Result**: System accurately detects and reports anomalies.

### Scenario 4: Configuration Update
1. Update Helm values (e.g., replica count)
2. Run `helm upgrade`
3. Verify rolling update completes
4. Verify new configuration is applied
5. Verify no service disruption

**Expected Result**: Seamless configuration updates.

### Scenario 5: Multi-User Access
1. Open dashboard in multiple browsers
2. Verify all instances update independently
3. Verify no conflicts or race conditions
4. Dismiss alert in one browser
5. Verify other browsers are not affected

**Expected Result**: Independent user sessions.

### Scenario 6: Long-Running Stability
1. Install Sentinel Mesh
2. Let run for 24 hours
3. Check for memory leaks
4. Check for crashed pods
5. Verify continuous data collection
6. Verify charts still updating

**Expected Result**: Stable operation over extended period.

### Scenario 7: Network Partition
1. Block network access to ML service
2. Verify frontend handles timeout gracefully
3. Verify Security view still loads (with old data)
4. Restore network access
5. Verify ML anomalies resume updating

**Expected Result**: Graceful degradation and recovery.

## Automated E2E Test Implementation

### Test Framework Options

**Option 1: Playwright (Recommended for Frontend)**
- Full browser automation
- TypeScript support
- Screenshot/video recording
- Network interception

**Option 2: Cypress**
- Great developer experience
- Time-travel debugging
- Easy to write tests

**Option 3: k6 (for API/Load Testing)**
- Performance testing
- Load generation
- Metrics collection

**Option 4: Bash + curl + kubectl**
- Simple integration tests
- No additional dependencies
- Easy to run in CI/CD

### Recommended Test Structure

```
tests/
├── e2e/
│   ├── frontend/              # Playwright/Cypress tests
│   │   ├── dashboard.spec.ts
│   │   ├── metrics.spec.ts
│   │   ├── security.spec.ts
│   │   ├── navigation.spec.ts
│   │   └── dark-mode.spec.ts
│   ├── api/                   # API integration tests
│   │   ├── health.test.sh
│   │   ├── metrics.test.sh
│   │   └── integration.test.sh
│   ├── k8s/                   # Kubernetes tests
│   │   ├── deployment.test.sh
│   │   ├── scaling.test.sh
│   │   └── resilience.test.sh
│   └── load/                  # Load tests
│       ├── stress-test.js
│       └── endurance-test.js
├── fixtures/                  # Test data
└── helpers/                   # Test utilities
```

## Pre-Release Checklist

Before creating a new release:

- [ ] All unit tests pass
- [ ] All integration tests pass
- [ ] All E2E scenarios pass
- [ ] Security scan shows no critical vulnerabilities
- [ ] Performance benchmarks meet targets
- [ ] Documentation is updated
- [ ] CHANGELOG is updated
- [ ] Breaking changes are documented
- [ ] Migration guide provided (if needed)
- [ ] Docker images build successfully
- [ ] Images are tagged with version
- [ ] Images are pushed to registry
- [ ] Helm chart version is updated
- [ ] Release notes are prepared

## Test Execution Commands

```bash docs-drift:skip
# Run all unit tests
make test-unit

# Run integration tests
make test-integration

# Run E2E tests
make test-e2e

# Run full test suite
make test-all

# Run in CI mode (with coverage)
make test-ci

# Run specific test suite
make test-frontend
make test-backend
make test-ml

# Load testing
make test-load

# Security scan
make test-security
```

## Continuous Testing

### On Every Commit (CI)
- Unit tests
- Linting
- Type checking
- Security scanning

### On Pull Request
- Integration tests
- E2E smoke tests (critical paths)
- Performance regression tests

### On Release Tag
- Full E2E test suite
- Load testing
- Security audit
- Image building and pushing
- Helm chart publishing

## Test Data & Fixtures

### Mock Data Requirements
- Sample Kubernetes metrics
- Sample anomaly data
- Sample prediction results
- Various alert scenarios
- Different user preferences

### Test Environment Setup
```bash docs-drift:skip
# Create test namespace
kubectl create namespace sentinel-mesh-test

# Install with test values
helm install sentinel-mesh-test ./deployments/helm/sentinel-mesh \
  --namespace sentinel-mesh-test \
  --values ./tests/fixtures/test-values.yaml

# Run E2E tests
npm run test:e2e

# Cleanup
kubectl delete namespace sentinel-mesh-test
```

## Success Criteria

A release is ready when:

1. **Functionality**: All items in this checklist pass ✅
2. **Performance**: Response times < SLA targets
3. **Reliability**: No crashes or errors in 24h test
4. **Security**: No critical/high vulnerabilities
5. **Documentation**: All features documented
6. **User Experience**: Smooth, bug-free UI operation

---

**Note**: This is a living document. Update as new features are added or requirements change.
