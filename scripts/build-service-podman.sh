#!/bin/bash

# Build Service with Podman Script
# This script builds a service using Podman

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if podman is installed
check_podman() {
    if ! command -v podman &> /dev/null; then
        print_error "Podman is not installed. Please install Podman first."
        print_status "Installation guide: https://podman.io/getting-started/installation"
        exit 1
    fi
    print_success "Podman is installed: $(podman --version)"
}

# Build service
build_service() {
    local service_name=$1
    local dockerfile_path="deployments/docker/Dockerfile.${service_name}"
    
    if [ -z "$service_name" ]; then
        print_error "Service name is required"
        print_status "Usage: $0 <service-name>"
        exit 1
    fi
    
    if [ ! -f "$dockerfile_path" ]; then
        print_error "Dockerfile not found: $dockerfile_path"
        exit 1
    fi
    
    print_status "Building service: $service_name"
    print_status "Using Dockerfile: $dockerfile_path"
    
    # Build the image
    podman build -f "$dockerfile_path" -t "${service_name}:latest" .
    
    print_success "Service $service_name built successfully!"
    print_status "Image: ${service_name}:latest"
}

# Main execution
main() {
    print_status "Starting service build with Podman..."
    
    check_podman
    build_service "$1"
}

# Run main function
main "$@"
