# GitHub Container Registry (GHCR) Setup

This document explains how to use Docker images from GitHub Container Registry for Sentinel Mesh.

## Overview

All Sentinel Mesh Docker images are published to GitHub Container Registry (GHCR) at:
```
ghcr.io/georg-nikola/sentinel-mesh/<service>
```

## Available Images

The following services are available as Docker images:

- **API Service**: `ghcr.io/georg-nikola/sentinel-mesh/api`
- **Collector Service**: `ghcr.io/georg-nikola/sentinel-mesh/collector`
- **Processor Service**: `ghcr.io/georg-nikola/sentinel-mesh/processor`
- **Analyzer Service**: `ghcr.io/georg-nikola/sentinel-mesh/analyzer`
- **Alerting Service**: `ghcr.io/georg-nikola/sentinel-mesh/alerting`
- **ML Service**: `ghcr.io/georg-nikola/sentinel-mesh/ml-service`
- **Frontend**: `ghcr.io/georg-nikola/sentinel-mesh/frontend`

## Image Tags

Images are tagged with multiple versions:

- **Semantic versions**: `0.1.0`, `0.1`, `0`
- **Latest**: `latest` (points to the most recent release)
- **SHA-based**: `main-<sha>` (development builds from main branch)

## Pulling Images

### Public Access

All images are public and can be pulled without authentication:

```bash docs-drift:skip
# Pull specific version
docker pull ghcr.io/georg-nikola/sentinel-mesh/api:0.1.0

# Pull latest version
docker pull ghcr.io/georg-nikola/sentinel-mesh/api:latest

# Pull all services
docker pull ghcr.io/georg-nikola/sentinel-mesh/api:latest
docker pull ghcr.io/georg-nikola/sentinel-mesh/collector:latest
docker pull ghcr.io/georg-nikola/sentinel-mesh/processor:latest
docker pull ghcr.io/georg-nikola/sentinel-mesh/analyzer:latest
docker pull ghcr.io/georg-nikola/sentinel-mesh/alerting:latest
docker pull ghcr.io/georg-nikola/sentinel-mesh/ml-service:latest
docker pull ghcr.io/georg-nikola/sentinel-mesh/frontend:latest
```

### With Authentication (for private packages)

If packages are set to private, authenticate with a GitHub Personal Access Token:

```bash docs-drift:skip
# Login to GHCR
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Pull image
docker pull ghcr.io/georg-nikola/sentinel-mesh/api:latest
```

## Using with Docker Compose

Update your `docker-compose.yml` to use GHCR images:

```yaml docs-drift:skip
version: '3.8'

services:
  api:
    image: ghcr.io/georg-nikola/sentinel-mesh/api:latest
    ports:
      - "8080:8080"
    environment:
      - REDIS_URL=redis://redis:6379

  collector:
    image: ghcr.io/georg-nikola/sentinel-mesh/collector:latest
    ports:
      - "8081:8080"

  ml-service:
    image: ghcr.io/georg-nikola/sentinel-mesh/ml-service:latest
    ports:
      - "8000:8000"

  frontend:
    image: ghcr.io/georg-nikola/sentinel-mesh/frontend:latest
    ports:
      - "80:80"
```

## Using with Kubernetes

### Update Helm Values

```yaml docs-drift:skip
# values.yaml
image:
  registry: ghcr.io
  repository: georg-nikola/sentinel-mesh
  tag: "0.1.0"
  pullPolicy: IfNotPresent
```

### Deploy with Helm

```bash docs-drift:skip
helm install sentinel-mesh ./deployments/helm/sentinel-mesh \
  --set image.registry=ghcr.io \
  --set image.repository=georg-nikola/sentinel-mesh \
  --set image.tag=0.1.0
```

### Direct Kubernetes Deployment

```yaml docs-drift:skip
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sentinel-mesh-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
      - name: api
        image: ghcr.io/georg-nikola/sentinel-mesh/api:0.1.0
        ports:
        - containerPort: 8080
```

## CI/CD Integration

### GitHub Actions

```yaml docs-drift:skip
- name: Pull image from GHCR
  run: |
    docker pull ghcr.io/georg-nikola/sentinel-mesh/api:latest
```

### GitLab CI

```yaml docs-drift:skip
deploy:
  script:
    - docker pull ghcr.io/georg-nikola/sentinel-mesh/api:$CI_COMMIT_TAG
    - docker run -d ghcr.io/georg-nikola/sentinel-mesh/api:$CI_COMMIT_TAG
```

## Publishing Images (Maintainers Only)

Images are automatically built and published by GitHub Actions when:

1. **Creating a Release**: Tag and publish a release on GitHub
2. **Manual Trigger**: Use workflow_dispatch with a version tag

```bash docs-drift:skip
# Create and push a tag
git tag v0.1.0
git push origin v0.1.0

# Create a GitHub release (triggers build)
gh release create v0.1.0 --title "Release 0.1.0" --notes "Release notes"
```

The release workflow will:
- Build all 7 service images
- Tag them with semantic versions
- Push to GHCR
- Test the images
- Update Helm charts

## Package Visibility

To make packages public (one-time setup):

1. Go to: https://github.com/georg-nikola?tab=packages
2. Click on each package (e.g., `sentinel-mesh/api`)
3. Go to "Package settings"
4. Under "Danger Zone", click "Change visibility"
5. Select "Public"
6. Confirm by typing the package name

## Image Sizes

All images use multi-stage builds for minimal size:

- **Go services** (API, Collector, etc.): ~20-30 MB (Alpine-based)
- **Python ML service**: ~200-300 MB
- **Frontend**: ~50-100 MB (nginx-based)

## Troubleshooting

### Authentication Issues

```bash docs-drift:skip
# Check if logged in
docker info | grep Username

# Re-login
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

### Image Not Found

```bash docs-drift:skip
# Check available tags
gh api /users/georg-nikola/packages/container/sentinel-mesh%2Fapi/versions

# Or browse packages
open https://github.com/georg-nikola?tab=packages
```

### Rate Limits

GHCR has generous rate limits:
- Anonymous: 1000 requests/hour
- Authenticated: 5000 requests/hour

## Links

- **Packages**: https://github.com/georg-nikola?tab=packages
- **Repository**: https://github.com/georg-nikola/sentinel-mesh
- **Documentation**: https://github.com/georg-nikola/sentinel-mesh/tree/main/docs
