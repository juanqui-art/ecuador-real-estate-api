#!/bin/bash

# Migration helper script for realty-core
# Uses golang-migrate for professional database migrations

set -e

# Configuration
MIGRATE_BIN="/Users/juanquizhpi/go/bin/migrate"
MIGRATIONS_DIR="migrations"
# Note: Currently using single SQL files (not up/down pairs)
# This works for development, but consider converting to up/down format for production
DATABASE_URL_DEFAULT="postgresql://juanquizhpi@localhost:5433/inmobiliaria_db?sslmode=disable"

# Get database URL from environment or use default
DATABASE_URL=${DATABASE_URL:-$DATABASE_URL_DEFAULT}

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if migrate binary exists
check_migrate() {
    if [ ! -f "$MIGRATE_BIN" ]; then
        log_error "golang-migrate not found at $MIGRATE_BIN"
        log_info "Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
}

# Validate database connection
validate_db() {
    log_info "Validating database connection..."
    if ! $MIGRATE_BIN -path $MIGRATIONS_DIR -database "$DATABASE_URL" version >/dev/null 2>&1; then
        log_error "Cannot connect to database: $DATABASE_URL"
        log_info "Check your DATABASE_URL environment variable or database status"
        exit 1
    fi
    log_success "Database connection validated"
}

# Show current migration version
show_version() {
    log_info "Current migration version:"
    VERSION=$($MIGRATE_BIN -path $MIGRATIONS_DIR -database "$DATABASE_URL" version 2>/dev/null || echo "No migrations applied")
    echo "📊 Version: $VERSION"
}

# Apply all pending migrations
migrate_up() {
    log_info "Applying pending migrations..."
    if $MIGRATE_BIN -path $MIGRATIONS_DIR -database "$DATABASE_URL" up; then
        log_success "Migrations applied successfully"
        show_version
    else
        log_error "Migration failed"
        exit 1
    fi
}

# Rollback one migration
migrate_down() {
    local steps=${1:-1}
    log_warning "Rolling back $steps migration(s)..."
    log_warning "This will modify your database schema!"
    
    if $MIGRATE_BIN -path $MIGRATIONS_DIR -database "$DATABASE_URL" down $steps; then
        log_success "Rollback completed"
        show_version
    else
        log_error "Rollback failed"
        exit 1
    fi
}

# Force set version (dangerous)
force_version() {
    local version=$1
    if [ -z "$version" ]; then
        log_error "Version number required"
        exit 1
    fi
    
    log_warning "Forcing version to $version (dangerous operation!)"
    if $MIGRATE_BIN -path $MIGRATIONS_DIR -database "$DATABASE_URL" force $version; then
        log_success "Version forced to $version"
        show_version
    else
        log_error "Force version failed"
        exit 1
    fi
}

# Create new migration files
create_migration() {
    local name=$1
    if [ -z "$name" ]; then
        log_error "Migration name required"
        echo "Usage: $0 create <migration_name>"
        exit 1
    fi
    
    log_info "Creating new migration: $name"
    if $MIGRATE_BIN create -ext sql -dir $MIGRATIONS_DIR -seq $name; then
        log_success "Migration files created"
        log_info "Don't forget to fill in the up and down SQL files!"
    else
        log_error "Migration creation failed"
        exit 1
    fi
}

# Show help
show_help() {
    echo "🏠 Realty Core Migration Helper"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  up              Apply all pending migrations"
    echo "  down [steps]    Rollback migrations (default: 1)"
    echo "  version         Show current migration version"
    echo "  create <name>   Create new migration files"
    echo "  force <version> Force set migration version (dangerous)"
    echo "  validate        Validate database connection"
    echo ""
    echo "Environment Variables:"
    echo "  DATABASE_URL    PostgreSQL connection string"
    echo "                  Default: $DATABASE_URL_DEFAULT"
    echo ""
    echo "Examples:"
    echo "  $0 up                           # Apply all pending migrations"
    echo "  $0 down 2                       # Rollback 2 migrations"
    echo "  $0 create add_property_features  # Create new migration"
    echo "  $0 force 15                     # Force version to 15"
}

# Main command handler
main() {
    check_migrate
    
    case "${1:-help}" in
        "up")
            validate_db
            migrate_up
            ;;
        "down")
            validate_db
            migrate_down "$2"
            ;;
        "version")
            validate_db
            show_version
            ;;
        "create")
            create_migration "$2"
            ;;
        "force")
            validate_db
            force_version "$2"
            ;;
        "validate")
            validate_db
            ;;
        "help"|"--help"|"-h")
            show_help
            ;;
        *)
            log_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"