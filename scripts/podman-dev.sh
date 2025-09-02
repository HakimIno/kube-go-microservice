#!/bin/bash

# Podman Development Environment Script
# This script sets up and runs the development environment using Podman

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

# Check if podman-compose is installed
check_podman_compose() {
    if ! command -v podman-compose &> /dev/null; then
        print_warning "podman-compose is not installed. Installing..."
        pip3 install podman-compose
    fi
    print_success "podman-compose is available"
}

# Build and run services
run_dev_environment() {
    print_status "Starting development environment with Podman..."
    
    # Change to the correct directory
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    cd "$SCRIPT_DIR/../deployments/podman-compose"
    
    # Stop any existing containers
    print_status "Stopping existing containers..."
    podman-compose -f docker-compose.dev.yml down || true
    
    # Build and start services
    print_status "Building and starting services..."
    podman-compose -f docker-compose.dev.yml up --build -d
    
    print_success "Development environment started!"
    print_status "Services available at:"
    print_status "  - User Service: http://localhost:8081"
    print_status "  - PostgreSQL: localhost:5432"
    print_status "  - Redis: localhost:6379"
    
    # Show logs (individual containers to avoid remote issue)
    print_status "Showing logs (Ctrl+C to stop)..."
    print_status "User Service logs:"
    podman-compose -f docker-compose.dev.yml logs -f user-service &
    USER_LOG_PID=$!
    
    print_status "PostgreSQL logs:"
    podman-compose -f docker-compose.dev.yml logs -f postgres &
    POSTGRES_LOG_PID=$!
    
    print_status "Redis logs:"
    podman-compose -f docker-compose.dev.yml logs -f redis &
    REDIS_LOG_PID=$!
    
    # Wait for any process to exit (macOS compatible)
    wait
    
    # Kill all background processes
    kill $USER_LOG_PID $POSTGRES_LOG_PID $REDIS_LOG_PID 2>/dev/null || true
}

# Cleanup function
cleanup() {
    print_status "Cleaning up..."
    # Get the script directory and navigate to the correct path
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PODMAN_COMPOSE_DIR="$SCRIPT_DIR/../deployments/podman-compose"
    
    if [ -d "$PODMAN_COMPOSE_DIR" ]; then
        cd "$PODMAN_COMPOSE_DIR"
        podman-compose -f docker-compose.dev.yml down
    else
        print_warning "Podman compose directory not found: $PODMAN_COMPOSE_DIR"
    fi
    
    # Kill background processes if they exist
    if [ ! -z "$USER_LOG_PID" ]; then
        kill $USER_LOG_PID 2>/dev/null || true
    fi
    if [ ! -z "$POSTGRES_LOG_PID" ]; then
        kill $POSTGRES_LOG_PID 2>/dev/null || true
    fi
    if [ ! -z "$REDIS_LOG_PID" ]; then
        kill $REDIS_LOG_PID 2>/dev/null || true
    fi
    
    print_success "Cleanup completed"
}

# Main execution
main() {
    print_status "Starting Podman development environment setup..."
    
    check_podman
    check_podman_compose
    
    # Set up signal handlers for cleanup
    trap cleanup EXIT INT TERM
    
    run_dev_environment
}

# Run main function
main "$@"
