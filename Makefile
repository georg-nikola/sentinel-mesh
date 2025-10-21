# Sentinel Mesh Makefile
# Provides build, test, and deployment automation

# Project configuration
PROJECT_NAME := sentinel-mesh
REGISTRY := ghcr.io/sentinel-mesh
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse HEAD)

# Go configuration
GO_VERSION := 1.21
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
CGO_ENABLED := 0

# Build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT) -w -s"

# Service names
SERVICES := collector processor analyzer alerting api

# Directories
BUILD_DIR := build
DIST_DIR := dist
DOCS_DIR := docs

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
NC := \033[0m # No Color

.PHONY: help
help: ## Display this help message
	@echo "$(CYAN)Sentinel Mesh Build System$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "$(BLUE)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

## Development Commands

.PHONY: setup
setup: ## Set up development environment
	@echo "$(GREEN)Setting up development environment...$(NC)"
	@echo "$(YELLOW)Installing Go dependencies...$(NC)"
	go mod download
	go mod tidy
	@echo "$(YELLOW)Installing Python dependencies...$(NC)"
	cd ml && pip install -r requirements.txt
	@echo "$(YELLOW)Installing Node.js dependencies...$(NC)"
	cd web && npm install
	@echo "$(GREEN)Development environment setup complete!$(NC)"

.PHONY: dev-up
dev-up: ## Start development dependencies (Kafka, Redis, etc.)
	@echo "$(GREEN)Starting development dependencies...$(NC)"
	docker-compose -f deployments/docker/docker-compose.dev.yml up -d
	@echo "$(GREEN)Development dependencies started!$(NC)"

.PHONY: dev-down
dev-down: ## Stop development dependencies
	@echo "$(YELLOW)Stopping development dependencies...$(NC)"
	docker-compose -f deployments/docker/docker-compose.dev.yml down
	@echo "$(GREEN)Development dependencies stopped!$(NC)"

## Build Commands

.PHONY: build
build: build-go build-ml build-frontend ## Build all components

.PHONY: build-go
build-go: $(addprefix build-,$(SERVICES)) ## Build all Go services

.PHONY: $(addprefix build-,$(SERVICES))
$(addprefix build-,$(SERVICES)): build-%:
	@echo "$(GREEN)Building $* service...$(NC)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		$(LDFLAGS) \
		-o $(BUILD_DIR)/$* \
		./cmd/$*
	@echo "$(GREEN)Built $* service successfully!$(NC)"

.PHONY: build-ml
build-ml: ## Build ML service
	@echo "$(GREEN)Building ML service...$(NC)"
	cd ml && python -m pip install --upgrade pip
	cd ml && pip install -r requirements.txt
	@echo "$(GREEN)ML service ready!$(NC)"

.PHONY: build-frontend
build-frontend: ## Build frontend application
	@echo "$(GREEN)Building frontend...$(NC)"
	cd web && npm ci
	cd web && npm run build
	@echo "$(GREEN)Frontend built successfully!$(NC)"

## Test Commands

.PHONY: test
test: test-go test-ml test-frontend ## Run all tests

.PHONY: test-go
test-go: ## Run Go tests
	@echo "$(GREEN)Running Go tests...$(NC)"
	go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)Go tests completed!$(NC)"

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	go test -v -tags=integration ./tests/integration/...
	@echo "$(GREEN)Integration tests completed!$(NC)"

.PHONY: test-ml
test-ml: ## Run ML tests
	@echo "$(GREEN)Running ML tests...$(NC)"
	cd ml && python -m pytest -v --cov=. --cov-report=term-missing
	@echo "$(GREEN)ML tests completed!$(NC)"

.PHONY: test-frontend
test-frontend: ## Run frontend tests
	@echo "$(GREEN)Running frontend tests...$(NC)"
	cd web && npm run test
	@echo "$(GREEN)Frontend tests completed!$(NC)"

.PHONY: coverage-go
coverage-go: ## Generate Go test coverage report
	@echo "$(GREEN)Generating Go coverage report...$(NC)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(NC)"

## Code Quality Commands

.PHONY: lint
lint: lint-go lint-ml lint-frontend ## Run all linters

.PHONY: lint-go
lint-go: ## Lint Go code
	@echo "$(GREEN)Linting Go code...$(NC)"
	golangci-lint run --config .golangci.yml
	@echo "$(GREEN)Go linting completed!$(NC)"

.PHONY: lint-ml
lint-ml: ## Lint Python code
	@echo "$(GREEN)Linting Python code...$(NC)"
	cd ml && flake8 .
	cd ml && black --check .
	cd ml && mypy . --ignore-missing-imports
	@echo "$(GREEN)Python linting completed!$(NC)"

.PHONY: lint-frontend
lint-frontend: ## Lint frontend code
	@echo "$(GREEN)Linting frontend code...$(NC)"
	cd web && npm run lint
	cd web && npm run type-check
	@echo "$(GREEN)Frontend linting completed!$(NC)"

.PHONY: format
format: format-go format-ml format-frontend ## Format all code

.PHONY: format-go
format-go: ## Format Go code
	@echo "$(GREEN)Formatting Go code...$(NC)"
	go fmt ./...
	goimports -w .
	@echo "$(GREEN)Go code formatted!$(NC)"

.PHONY: format-ml
format-ml: ## Format Python code
	@echo "$(GREEN)Formatting Python code...$(NC)"
	cd ml && black .
	cd ml && isort .
	@echo "$(GREEN)Python code formatted!$(NC)"

.PHONY: format-frontend
format-frontend: ## Format frontend code
	@echo "$(GREEN)Formatting frontend code...$(NC)"
	cd web && npm run lint -- --fix
	@echo "$(GREEN)Frontend code formatted!$(NC)"

## Docker Commands

.PHONY: docker-build
docker-build: $(addprefix docker-build-,$(SERVICES)) docker-build-ml docker-build-frontend ## Build all Docker images

.PHONY: $(addprefix docker-build-,$(SERVICES))
$(addprefix docker-build-,$(SERVICES)): docker-build-%:
	@echo "$(GREEN)Building Docker image for $*...$(NC)"
	docker build \
		--build-arg SERVICE=$* \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg GIT_COMMIT=$(GIT_COMMIT) \
		-t $(REGISTRY)/$*:$(VERSION) \
		-t $(REGISTRY)/$*:latest \
		-f deployments/docker/Dockerfile.go \
		.
	@echo "$(GREEN)Docker image for $* built successfully!$(NC)"

.PHONY: docker-build-ml
docker-build-ml: ## Build ML service Docker image
	@echo "$(GREEN)Building Docker image for ML service...$(NC)"
	docker build \
		--build-arg VERSION=$(VERSION) \
		-t $(REGISTRY)/ml-service:$(VERSION) \
		-t $(REGISTRY)/ml-service:latest \
		-f deployments/docker/Dockerfile.ml \
		.
	@echo "$(GREEN)Docker image for ML service built successfully!$(NC)"

.PHONY: docker-build-frontend
docker-build-frontend: ## Build frontend Docker image
	@echo "$(GREEN)Building Docker image for frontend...$(NC)"
	docker build \
		--build-arg VERSION=$(VERSION) \
		-t $(REGISTRY)/frontend:$(VERSION) \
		-t $(REGISTRY)/frontend:latest \
		-f deployments/docker/Dockerfile.frontend \
		.
	@echo "$(GREEN)Docker image for frontend built successfully!$(NC)"

.PHONY: docker-push
docker-push: ## Push Docker images to registry
	@echo "$(GREEN)Pushing Docker images...$(NC)"
	@for service in $(SERVICES) ml-service frontend; do \
		echo "$(YELLOW)Pushing $$service...$(NC)"; \
		docker push $(REGISTRY)/$$service:$(VERSION); \
		docker push $(REGISTRY)/$$service:latest; \
	done
	@echo "$(GREEN)All images pushed successfully!$(NC)"

## Kubernetes Commands

.PHONY: k8s-deploy
k8s-deploy: ## Deploy to Kubernetes using Helm
	@echo "$(GREEN)Deploying to Kubernetes...$(NC)"
	helm upgrade --install sentinel-mesh \
		./deployments/helm/sentinel-mesh \
		--namespace sentinel-mesh \
		--create-namespace \
		--set image.tag=$(VERSION)
	@echo "$(GREEN)Deployment completed!$(NC)"

.PHONY: k8s-uninstall
k8s-uninstall: ## Uninstall from Kubernetes
	@echo "$(YELLOW)Uninstalling from Kubernetes...$(NC)"
	helm uninstall sentinel-mesh --namespace sentinel-mesh
	@echo "$(GREEN)Uninstall completed!$(NC)"

.PHONY: k8s-status
k8s-status: ## Check Kubernetes deployment status
	@echo "$(GREEN)Checking deployment status...$(NC)"
	helm status sentinel-mesh --namespace sentinel-mesh
	kubectl get pods --namespace sentinel-mesh

## Development Utilities

.PHONY: run-collector
run-collector: build-collector ## Run collector service locally
	@echo "$(GREEN)Starting collector service...$(NC)"
	./$(BUILD_DIR)/collector --config ./configs/dev/config.yaml

.PHONY: run-processor
run-processor: build-processor ## Run processor service locally
	@echo "$(GREEN)Starting processor service...$(NC)"
	./$(BUILD_DIR)/processor --config ./configs/dev/config.yaml

.PHONY: run-api
run-api: build-api ## Run API service locally
	@echo "$(GREEN)Starting API service...$(NC)"
	./$(BUILD_DIR)/api --config ./configs/dev/config.yaml

.PHONY: run-ml
run-ml: ## Run ML service locally
	@echo "$(GREEN)Starting ML service...$(NC)"
	cd ml && python main.py

.PHONY: run-frontend
run-frontend: ## Run frontend development server
	@echo "$(GREEN)Starting frontend development server...$(NC)"
	cd web && npm run dev

## Cleanup Commands

.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html
	cd web && rm -rf dist node_modules
	cd ml && find . -name "*.pyc" -delete
	cd ml && find . -name "__pycache__" -delete
	@echo "$(GREEN)Cleanup completed!$(NC)"

.PHONY: clean-docker
clean-docker: ## Clean Docker images and containers
	@echo "$(YELLOW)Cleaning Docker resources...$(NC)"
	docker system prune -f
	@echo "$(GREEN)Docker cleanup completed!$(NC)"

## Documentation

.PHONY: docs
docs: ## Generate documentation
	@echo "$(GREEN)Generating documentation...$(NC)"
	@mkdir -p $(DOCS_DIR)
	go doc -all ./... > $(DOCS_DIR)/go-docs.txt
	@echo "$(GREEN)Documentation generated!$(NC)"

.PHONY: serve-docs
serve-docs: ## Serve documentation locally
	@echo "$(GREEN)Serving documentation...$(NC)"
	cd $(DOCS_DIR) && python -m http.server 8000

## Release Commands

.PHONY: release
release: test docker-build docker-push ## Create a release (test, build, and push images)
	@echo "$(GREEN)Release $(VERSION) completed!$(NC)"

.PHONY: version
version: ## Show version information
	@echo "$(CYAN)Sentinel Mesh Version Information$(NC)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(GO_VERSION)"
	@echo "OS/Arch: $(GOOS)/$(GOARCH)"

## Quick Commands

.PHONY: dev
dev: dev-up build-go run-collector ## Quick development setup (start deps and collector)

.PHONY: all
all: clean setup build test ## Build and test everything

# Default target
.DEFAULT_GOAL := help