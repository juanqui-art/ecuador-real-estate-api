# Realty Deployment

Setup deployment for: $ARGUMENTS

## Context - Production Deployment
Deploy real estate API with:
- **Docker containers:** Multi-stage builds
- **PostgreSQL:** Production database setup
- **Reverse proxy:** Nginx configuration
- **SSL certificates:** Let's Encrypt
- **Monitoring:** Health checks and logging

## Deployment Patterns:
1. **Docker setup:**
   - Multi-stage Dockerfile for minimal image
   - docker-compose for local development
   - Production container orchestration
   - Health checks and graceful shutdown

2. **Database setup:**
   - PostgreSQL with persistent storage
   - Connection pooling
   - Backup and restore procedures
   - Migration management

3. **Reverse proxy:**
   - Nginx configuration
   - SSL termination
   - Load balancing
   - Static file serving

## Docker configuration:
```dockerfile
# Multi-stage build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o realty-api ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/realty-api .
EXPOSE 8080
CMD ["./realty-api"]
```

## Environment configuration:
```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  api:
    image: realty-api:latest
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/realty
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
    
  db:
    image: postgres:15
    environment:
      - POSTGRES_DB=realty
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

## Common deployment scenarios:
- "production Docker setup"
- "add health checks"
- "configure environment variables"
- "setup SSL certificates"
- "implement blue-green deployment"

## Monitoring and logging:
- Application metrics
- Database monitoring
- Error tracking
- Log aggregation
- Alerting setup

## Output format:
- Docker configuration
- Environment setup
- Deployment scripts
- Monitoring configuration
- Security recommendations