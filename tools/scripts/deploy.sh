#!/bin/bash

# Production deployment script for Realty Core
# Handles building, testing, and deploying the application

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
COMPOSE_FILE="docker-compose.production.yml"
ENV_FILE=".env.production"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

show_help() {
    cat << EOF
Production Deployment Script for Realty Core

Usage: $0 [COMMAND] [OPTIONS]

COMMANDS:
    build       Build Docker images
    test        Run tests before deployment
    deploy      Deploy to production
    stop        Stop all services
    restart     Restart all services
    logs        Show logs from services
    status      Show status of services
    backup      Backup database
    restore     Restore database from backup
    update      Update to latest version
    health      Check health of all services

OPTIONS:
    -e, --env-file FILE     Environment file (default: $ENV_FILE)
    -f, --compose-file FILE Compose file (default: $COMPOSE_FILE)
    -v, --verbose          Verbose output
    -h, --help             Show this help message

EXAMPLES:
    $0 build                    # Build images
    $0 test                     # Run tests
    $0 deploy                   # Full deployment
    $0 deploy -v                # Verbose deployment
    $0 logs app                 # Show app logs
    $0 backup                   # Backup database
    $0 restart app              # Restart only app service
    $0 status                   # Show all service status

NOTES:
    - Make sure Docker and Docker Compose are installed
    - Configure $ENV_FILE before deployment
    - Run tests before deploying to production
    - Database backups are stored in backups/ directory
EOF
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi
    
    # Check Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose is not installed"
        exit 1
    fi
    
    # Check if Docker is running
    if ! docker info >/dev/null 2>&1; then
        log_error "Docker is not running"
        exit 1
    fi
    
    log_success "Prerequisites check passed"
}

# Check environment file
check_env_file() {
    if [[ ! -f "$ENV_FILE" ]]; then
        log_error "Environment file not found: $ENV_FILE"
        log_info "Copy .env.production.example to $ENV_FILE and configure it"
        exit 1
    fi
    
    # Check for required variables
    required_vars=("POSTGRES_PASSWORD" "DATABASE_URL")
    for var in "${required_vars[@]}"; do
        if ! grep -q "^${var}=" "$ENV_FILE"; then
            log_error "Required environment variable $var not found in $ENV_FILE"
            exit 1
        fi
    done
    
    log_success "Environment file validation passed"
}

# Build Docker images
build_images() {
    log_info "Building Docker images..."
    
    cd "$PROJECT_DIR"
    
    # Build the main application image
    docker build -t realty-core:latest .
    
    # Also build distroless version
    docker build -f Dockerfile.distroless -t realty-core:latest-distroless .
    
    log_success "Images built successfully"
}

# Run tests
run_tests() {
    log_info "Running tests..."
    
    cd "$PROJECT_DIR"
    
    # Run Go tests
    log_info "Running Go tests..."
    go test ./... -v
    
    # Run integration tests if they exist
    if [[ -f "tests/integration_test.go" ]]; then
        log_info "Running integration tests..."
        go test ./tests/... -v -tags=integration
    fi
    
    log_success "All tests passed"
}

# Deploy services
deploy_services() {
    log_info "Deploying services..."
    
    cd "$PROJECT_DIR"
    
    # Pull latest images
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" pull
    
    # Deploy with zero downtime
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --remove-orphans
    
    # Wait for services to be healthy
    log_info "Waiting for services to be healthy..."
    sleep 30
    
    # Check health
    check_health
    
    log_success "Deployment completed"
}

# Stop services
stop_services() {
    log_info "Stopping services..."
    
    cd "$PROJECT_DIR"
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" down
    
    log_success "Services stopped"
}

# Restart services
restart_services() {
    local service="${1:-}"
    
    cd "$PROJECT_DIR"
    
    if [[ -n "$service" ]]; then
        log_info "Restarting service: $service"
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" restart "$service"
    else
        log_info "Restarting all services..."
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" restart
    fi
    
    log_success "Services restarted"
}

# Show logs
show_logs() {
    local service="${1:-}"
    local follow="${2:-false}"
    
    cd "$PROJECT_DIR"
    
    if [[ "$follow" == "true" ]]; then
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" logs -f $service
    else
        docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" logs --tail=100 $service
    fi
}

# Show service status
show_status() {
    log_info "Service Status:"
    
    cd "$PROJECT_DIR"
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" ps
    
    echo ""
    log_info "Container Resource Usage:"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"
}

# Check health of services
check_health() {
    log_info "Checking service health..."
    
    cd "$PROJECT_DIR"
    
    # Check app health
    if docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec -T app /inmobiliaria -health-check; then
        log_success "Application is healthy"
    else
        log_error "Application health check failed"
        return 1
    fi
    
    # Check database health
    if docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec -T postgres pg_isready -U inmobiliaria_user -d inmobiliaria_db; then
        log_success "Database is healthy"
    else
        log_error "Database health check failed"
        return 1
    fi
    
    log_success "All services are healthy"
}

# Backup database
backup_database() {
    log_info "Creating database backup..."
    
    cd "$PROJECT_DIR"
    
    # Create backups directory
    mkdir -p backups
    
    # Generate backup filename with timestamp
    backup_file="backups/realty_core_$(date +%Y%m%d_%H%M%S).sql"
    
    # Create backup
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec -T postgres \
        pg_dump -U inmobiliaria_user inmobiliaria_db > "$backup_file"
    
    # Compress backup
    gzip "$backup_file"
    
    log_success "Database backup created: ${backup_file}.gz"
    
    # Clean old backups (keep last 7 days)
    find backups/ -name "*.sql.gz" -mtime +7 -delete
    log_info "Old backups cleaned up"
}

# Restore database
restore_database() {
    local backup_file="$1"
    
    if [[ -z "$backup_file" ]]; then
        log_error "Backup file not specified"
        echo "Usage: $0 restore <backup_file>"
        exit 1
    fi
    
    if [[ ! -f "$backup_file" ]]; then
        log_error "Backup file not found: $backup_file"
        exit 1
    fi
    
    log_warning "This will replace the current database. Are you sure? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        log_info "Restore cancelled"
        exit 0
    fi
    
    log_info "Restoring database from: $backup_file"
    
    cd "$PROJECT_DIR"
    
    # Stop app to prevent connections
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" stop app
    
    # Restore database
    if [[ "$backup_file" =~ \.gz$ ]]; then
        zcat "$backup_file" | docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec -T postgres \
            psql -U inmobiliaria_user -d inmobiliaria_db
    else
        cat "$backup_file" | docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" exec -T postgres \
            psql -U inmobiliaria_user -d inmobiliaria_db
    fi
    
    # Start app
    docker-compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" start app
    
    log_success "Database restored successfully"
}

# Update application
update_application() {
    log_info "Updating application..."
    
    # Create backup before update
    backup_database
    
    # Pull latest changes (if using git)
    if [[ -d ".git" ]]; then
        log_info "Pulling latest changes..."
        git pull origin main
    fi
    
    # Build new images
    build_images
    
    # Run tests
    run_tests
    
    # Deploy
    deploy_services
    
    log_success "Application updated successfully"
}

# Main script logic
main() {
    local command="$1"
    local verbose=false
    
    # Parse options
    shift
    while [[ $# -gt 0 ]]; do
        case $1 in
            -e|--env-file)
                ENV_FILE="$2"
                shift 2
                ;;
            -f|--compose-file)
                COMPOSE_FILE="$2"
                shift 2
                ;;
            -v|--verbose)
                verbose=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                break
                ;;
        esac
    done
    
    # Set verbose mode
    if [[ "$verbose" == true ]]; then
        set -x
    fi
    
    # Check prerequisites for most commands
    case "$command" in
        help|--help|-h)
            show_help
            exit 0
            ;;
        *)
            check_prerequisites
            check_env_file
            ;;
    esac
    
    # Execute command
    case "$command" in
        build)
            build_images
            ;;
        test)
            run_tests
            ;;
        deploy)
            build_images
            run_tests
            deploy_services
            ;;
        stop)
            stop_services
            ;;
        restart)
            restart_services "$1"
            ;;
        logs)
            show_logs "$1" "$2"
            ;;
        status)
            show_status
            ;;
        health)
            check_health
            ;;
        backup)
            backup_database
            ;;
        restore)
            restore_database "$1"
            ;;
        update)
            update_application
            ;;
        *)
            log_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Check if script is run directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi