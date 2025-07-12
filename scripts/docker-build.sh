#!/bin/bash

# Docker build script for Realty Core
# Supports both scratch and distroless base images

set -e

# Configuration
IMAGE_NAME="realty-core"
VERSION="1.9.0"
REGISTRY=""

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
Docker Build Script for Realty Core

Usage: $0 [OPTIONS]

OPTIONS:
    -t, --type TYPE          Image type: scratch (default) or distroless
    -v, --version VERSION    Image version (default: $VERSION)
    -r, --registry REGISTRY  Docker registry prefix
    -p, --push              Push image to registry after build
    -l, --latest            Also tag as latest
    -h, --help              Show this help message

EXAMPLES:
    $0                                   # Build scratch-based image
    $0 -t distroless                     # Build distroless-based image
    $0 -t scratch -v 2.0.0 -l           # Build scratch image with custom version and latest tag
    $0 -r myregistry.com -p             # Build and push to registry
    $0 -t distroless -r myregistry.com -p -l  # Build distroless, tag as latest, and push

NOTES:
    - The scratch image is smaller but less debuggable
    - The distroless image is more secure and easier to debug
    - Make sure Docker is running before executing this script
EOF
}

# Default values
IMAGE_TYPE="scratch"
PUSH_IMAGE=false
TAG_LATEST=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--type)
            IMAGE_TYPE="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        -p|--push)
            PUSH_IMAGE=true
            shift
            ;;
        -l|--latest)
            TAG_LATEST=true
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Validate image type
if [[ "$IMAGE_TYPE" != "scratch" && "$IMAGE_TYPE" != "distroless" ]]; then
    log_error "Invalid image type: $IMAGE_TYPE. Must be 'scratch' or 'distroless'"
    exit 1
fi

# Build image name with registry prefix if provided
if [[ -n "$REGISTRY" ]]; then
    FULL_IMAGE_NAME="$REGISTRY/$IMAGE_NAME"
else
    FULL_IMAGE_NAME="$IMAGE_NAME"
fi

# Dockerfile selection
if [[ "$IMAGE_TYPE" == "distroless" ]]; then
    DOCKERFILE="Dockerfile.distroless"
    IMAGE_TAG_SUFFIX="-distroless"
else
    DOCKERFILE="Dockerfile"
    IMAGE_TAG_SUFFIX=""
fi

# Check if Dockerfile exists
if [[ ! -f "$DOCKERFILE" ]]; then
    log_error "Dockerfile not found: $DOCKERFILE"
    exit 1
fi

# Build tags
VERSION_TAG="$FULL_IMAGE_NAME:$VERSION$IMAGE_TAG_SUFFIX"
LATEST_TAG="$FULL_IMAGE_NAME:latest$IMAGE_TAG_SUFFIX"

log_info "Starting Docker build..."
log_info "Image type: $IMAGE_TYPE"
log_info "Dockerfile: $DOCKERFILE"
log_info "Version tag: $VERSION_TAG"

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    log_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Build the image
log_info "Building image with version tag..."
docker build \
    -f "$DOCKERFILE" \
    -t "$VERSION_TAG" \
    --build-arg VERSION="$VERSION" \
    --build-arg BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')" \
    .

if [[ $? -eq 0 ]]; then
    log_success "Successfully built $VERSION_TAG"
else
    log_error "Failed to build image"
    exit 1
fi

# Tag as latest if requested
if [[ "$TAG_LATEST" == true ]]; then
    log_info "Tagging as latest..."
    docker tag "$VERSION_TAG" "$LATEST_TAG"
    log_success "Tagged as $LATEST_TAG"
fi

# Push to registry if requested
if [[ "$PUSH_IMAGE" == true ]]; then
    if [[ -z "$REGISTRY" ]]; then
        log_warning "No registry specified, pushing to Docker Hub"
    fi
    
    log_info "Pushing $VERSION_TAG to registry..."
    docker push "$VERSION_TAG"
    
    if [[ "$TAG_LATEST" == true ]]; then
        log_info "Pushing $LATEST_TAG to registry..."
        docker push "$LATEST_TAG"
    fi
    
    log_success "Successfully pushed to registry"
fi

# Show image information
log_info "Build completed successfully!"
echo ""
echo "Image details:"
docker images "$FULL_IMAGE_NAME" --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.CreatedAt}}"

echo ""
log_success "To run the container:"
echo "docker run -d -p 8080:8080 -e DATABASE_URL='your_db_url' $VERSION_TAG"

echo ""
log_success "To test the health check:"
echo "docker run --rm $VERSION_TAG -health-check"