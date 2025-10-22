# Multi-stage build for Go services
# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for service name
ARG SERVICE
ARG VERSION=dev
ARG BUILD_TIME
ARG GIT_COMMIT

# Build the service
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT} -w -s" \
    -o /build/service \
    ./cmd/${SERVICE}

# Stage 2: Runtime
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 sentinel && \
    adduser -D -u 1000 -G sentinel sentinel

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/service /app/service

# Copy configuration files
COPY --from=builder /build/configs /app/configs

# Change ownership
RUN chown -R sentinel:sentinel /app

# Switch to non-root user
USER sentinel

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/service", "healthcheck"] || exit 1

# Expose default port
EXPOSE 8080

# Run the service
ENTRYPOINT ["/app/service"]
CMD ["--config", "/app/configs/dev/config.yaml"]
