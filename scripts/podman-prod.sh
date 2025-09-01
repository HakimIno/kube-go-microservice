#!/bin/bash

# Podman Production Environment Script
# This script sets up and runs the production environment using Podman

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

# Build and run production services
run_prod_environment() {
    print_status "Starting production environment with Podman..."
    
    cd deployments/podman-compose
    
    # Stop any existing containers
    print_status "Stopping existing containers..."
    podman-compose -f docker-compose.yml down || true
    
    # Build and start services
    print_status "Building and starting production services..."
    podman-compose -f docker-compose.yml up --build -d
    
    print_success "Production environment started!"
    print_status "Services available at:"
    print_status "  - User Service: http://localhost:8081"
    print_status "  - PostgreSQL: localhost:5432"
    print_status "  - Redis: localhost:6379"
    
    # Show status
    print_status "Container status:"
    podman-compose -f docker-compose.yml ps
}

# Stop production environment
stop_prod_environment() {
    print_status "Stopping production environment..."
    cd deployments/podman-compose
    podman-compose -f docker-compose.yml down
    print_success "Production environment stopped"
}

# Show logs
show_logs() {
    print_status "Showing logs..."
    cd deployments/podman-compose
    podman-compose -f docker-compose.yml logs -f
}

# Main execution
main() {
    case "${1:-start}" in
        "start")
            print_status "Starting Podman production environment setup..."
            check_podman
            check_podman_compose
            run_prod_environment
            ;;
        "stop")
            stop_prod_environment
            ;;
        "logs")
            show_logs
            ;;
        "restart")
            stop_prod_environment
            sleep 2
            run_prod_environment
            ;;
        *)
            print_error "Usage: $0 {start|stop|logs|restart}"
            print_status "  start   - Start production environment (default)"
            print_status "  stop    - Stop production environment"
            print_status "  logs    - Show logs"
            print_status "  restart - Restart production environment"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
