#!/bin/bash

# Unified Production Environment Script
# Supports both Docker and Podman
# Usage: ./scripts/prod.sh [docker|podman] {start|stop|logs|restart}

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

# Default container runtime
CONTAINER_RUNTIME="docker"
ACTION="start"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        docker|podman)
            CONTAINER_RUNTIME="$1"
            shift
            ;;
        start|stop|logs|restart)
            ACTION="$1"
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [docker|podman] {start|stop|logs|restart}"
            echo ""
            echo "Options:"
            echo "  docker|podman  Container runtime to use (default: docker)"
            echo "  start          Start production environment (default)"
            echo "  stop           Stop production environment"
            echo "  logs           Show logs"
            echo "  restart        Restart production environment"
            echo "  -h, --help     Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                    # Use Docker and start (default)"
            echo "  $0 podman            # Use Podman and start"
            echo "  $0 docker stop       # Use Docker and stop"
            echo "  $0 podman logs       # Use Podman and show logs"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

print_status "🚀 Production environment script using $CONTAINER_RUNTIME..."

# Check if container runtime is installed
check_runtime() {
    if [ "$CONTAINER_RUNTIME" = "docker" ]; then
        if ! command -v docker &> /dev/null; then
            print_error "Docker is not installed!"
            exit 1
        fi
        if ! docker info > /dev/null 2>&1; then
            print_error "Docker is not running!"
            exit 1
        fi
        print_success "Docker is available: $(docker --version)"
    else
        if ! command -v podman &> /dev/null; then
            print_error "Podman is not installed. Please install Podman first."
            print_status "Installation guide: https://podman.io/getting-started/installation"
            exit 1
        fi
        print_success "Podman is available: $(podman --version)"
        
        # Check if podman-compose is installed
        if ! command -v podman-compose &> /dev/null; then
            print_warning "podman-compose is not installed. Installing..."
            pip3 install podman-compose
        fi
        print_success "podman-compose is available"
    fi
}

# Start production environment
start_prod_environment() {
    print_status "🚀 Starting production environment with $CONTAINER_RUNTIME..."
    
    if [ "$CONTAINER_RUNTIME" = "docker" ]; then
        cd deployments/docker-compose
        docker compose -f docker-compose.prod.yml up --build -d
        print_success "✅ Production environment started with Docker!"
    else
        cd deployments/podman-compose
        podman-compose -f docker-compose.yml up --build -d
        print_success "✅ Production environment started with Podman!"
    fi
    
    print_status "🌐 Services available at:"
    print_status "  - User Service: http://localhost:8081"
    print_status "  - PostgreSQL: localhost:5432"
    print_status "  - Redis: localhost:6379"
    
    # Show status
    print_status "📊 Container status:"
    if [ "$CONTAINER_RUNTIME" = "docker" ]; then
        docker compose -f docker-compose.prod.yml ps
    else
        podman-compose -f docker-compose.yml ps
    fi
}

# Stop production environment
stop_prod_environment() {
    print_status "🛑 Stopping production environment..."
    
    if [ "$CONTAINER_RUNTIME" = "docker" ]; then
        cd deployments/docker-compose
        docker compose -f docker-compose.prod.yml down
    else
        cd deployments/podman-compose
        podman-compose -f docker-compose.yml down
    fi
    
    print_success "✅ Production environment stopped"
}

# Show logs
show_logs() {
    print_status "📋 Showing logs..."
    
    if [ "$CONTAINER_RUNTIME" = "docker" ]; then
        cd deployments/docker-compose
        docker compose -f docker-compose.prod.yml logs -f
    else
        cd deployments/podman-compose
        podman-compose -f docker-compose.yml logs -f
    fi
}

# Restart production environment
restart_prod_environment() {
    print_status "🔄 Restarting production environment..."
    stop_prod_environment
    sleep 2
    start_prod_environment
}

# Main execution
main() {
    check_runtime
    
    case "$ACTION" in
        "start")
            start_prod_environment
            ;;
        "stop")
            stop_prod_environment
            ;;
        "logs")
            show_logs
            ;;
        "restart")
            restart_prod_environment
            ;;
        *)
            print_error "Unknown action: $ACTION"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
