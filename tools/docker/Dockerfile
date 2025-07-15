# Multi-stage Docker build for production-optimized Go application
# Stage 1: Build stage
FROM golang:1.24-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
# - Static linking for smaller binary
# - Strip debug information
# - Disable CGO for better portability
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o inmobiliaria \
    ./cmd/server

# Stage 2: Runtime stage
FROM scratch

# Copy timezone data from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates for HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create non-root user for security
COPY --from=builder /etc/passwd /etc/passwd

# Copy the built binary
COPY --from=builder /app/inmobiliaria /inmobiliaria

# Set binary as executable
# Note: This is implicit with the COPY command in scratch

# Expose application port
EXPOSE 8080

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=10s --start-period=60s --retries=3 \
    CMD ["/inmobiliaria", "-health-check"]

# Set metadata labels
LABEL maintainer="Juan Quizhpi <juan@example.com>" \
      description="Realty Core - Sistema de gestión inmobiliaria" \
      version="1.9.0" \
      org.opencontainers.image.title="realty-core" \
      org.opencontainers.image.description="Sistema de gestión inmobiliaria para Ecuador" \
      org.opencontainers.image.version="1.9.0" \
      org.opencontainers.image.vendor="Realty Core" \
      org.opencontainers.image.licenses="MIT"

# Define environment variables with default values
ENV PORT=8080 \
    LOG_LEVEL=info \
    ENV=production \
    GIN_MODE=release

# Use non-root user (if we had created one)
# USER nobody

# Define startup command
ENTRYPOINT ["/inmobiliaria"]

# Document expected volume mounts
VOLUME ["/data"]

# Document expected environment variables
# Required:
# - DATABASE_URL: PostgreSQL connection string
# Optional:
# - PORT: Application port (default: 8080)
# - LOG_LEVEL: Logging level (default: info)
# - ENV: Environment (default: production)
# - MAX_UPLOAD_SIZE: Max file upload size (default: 10MB)
# - CACHE_ENABLED: Enable LRU cache (default: true)
# - CACHE_SIZE_MB: Cache size in MB (default: 50)