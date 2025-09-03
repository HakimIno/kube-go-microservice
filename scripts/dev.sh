#!/bin/bash

# Unified Development Environment Script
# Supports both Docker and Podman
# Usage: ./scripts/dev.sh [docker|podman] [--build]

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
BUILD_FLAG=""

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        docker|podman)
            CONTAINER_RUNTIME="$1"
            shift
            ;;
        --build)
            BUILD_FLAG="--build"
            shift
            ;;
        -h|--help)
            echo "Usage: $0 [docker|podman] [--build]"
            echo ""
            echo "Options:"
            echo "  docker|podman  Container runtime to use (default: docker)"
            echo "  --build        Build images before starting"
            echo "  -h, --help     Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0                    # Use Docker (default)"
            echo "  $0 podman            # Use Podman"
            echo "  $0 docker --build    # Use Docker with build"
            echo "  $0 podman --build    # Use Podman with build"
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            echo "Use -h or --help for usage information"
            exit 1
            ;;
    esac
done

print_status "🚀 Starting development environment with $CONTAINER_RUNTIME..."

# Function to cleanup on exit
cleanup() {
    print_status "🛑 Stopping development environment..."
    if [ "$CONTAINER_RUNTIME" = "docker" ]; then
        docker compose -f deployments/docker-compose/docker-compose.dev.yml down
    else
        cd deployments/podman-compose
        podman-compose -f docker-compose.dev.yml down
    fi
    exit 0
}

# Set up trap for cleanup
trap cleanup SIGINT SIGTERM

# Check if container runtime is installed
check_runtime() {
    if [ "$CONTAINER_RUNTIME" = "docker" ]; then
        if ! command -v docker &> /dev/null; then
            print_error "Docker is not installed!"
            print_status "📝 Installation instructions:"
            echo "  🍎 macOS: Download from https://docs.docker.com/desktop/install/mac-install/"
            echo "  🐧 Linux: Run 'curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh'"
            echo "  🪟 Windows: Download from https://docs.docker.com/desktop/install/windows-install/"
            exit 1
        fi
        
        # Check if Docker is running
        if ! docker info > /dev/null 2>&1; then
            print_error "Docker is not running!"
            if [[ "$OSTYPE" == "darwin"* ]]; then
                if [ -d "/Applications/Docker.app" ]; then
                    print_status "🍎 Starting Docker Desktop..."
                    open /Applications/Docker.app
                    print_status "⏳ Waiting for Docker to start..."
                    for i in {1..60}; do
                        if docker info > /dev/null 2>&1; then
                            print_success "✅ Docker started successfully!"
                            break
                        fi
                        sleep 1
                        if [ $((i % 10)) -eq 0 ]; then
                            print_status "⏳ Still waiting... ($i/60 seconds)"
                        fi
                    done
                fi
            fi
        fi
        print_success "Docker is running!"
    else
        if ! command -v podman &> /dev/null; then
            print_error "Podman is not installed. Please install Podman first."
            print_status "Installation guide: https://podman.io/getting-started/installation"
            exit 1
        fi
        print_success "Podman is installed: $(podman --version)"
        
        # Check if podman-compose is installed
        if ! command -v podman-compose &> /dev/null; then
            print_warning "podman-compose is not installed. Installing..."
            pip3 install podman-compose
        fi
        print_success "podman-compose is available"
    fi
}

# Generate swagger docs
generate_swagger() {
    print_status "📝 Generating swagger documentation..."
    swag init -g cmd/user-service/main.go
}

# Start development environment
start_dev_environment() {
    if [ "$CONTAINER_RUNTIME" = "docker" ]; then
        print_status "🔨 Starting Docker development environment..."
        docker compose -f deployments/docker-compose/docker-compose.dev.yml up -d $BUILD_FLAG
        
        print_success "✅ Docker development environment is running!"
        print_status "🌐 Services available at:"
        print_status "  📖 Swagger UI: http://localhost:8081/swagger/index.html"
        print_status "  🔗 API Base: http://localhost:8081"
        print_status "  🐘 PostgreSQL: localhost:5432"
        print_status "  🔴 Redis: localhost:6379"
        
        print_status "📋 Showing service logs (Ctrl+C to stop)..."
        docker compose -f deployments/docker-compose/docker-compose.dev.yml logs -f
    else
        print_status "🔨 Starting Podman development environment..."
        cd deployments/podman-compose
        
        # Stop any existing containers
        podman-compose -f docker-compose.dev.yml down || true
        
        # Build and start services
        podman-compose -f docker-compose.dev.yml up -d $BUILD_FLAG
        
        print_success "✅ Podman development environment is running!"
        print_status "🌐 Services available at:"
        print_status "  📖 Swagger UI: http://localhost:8081/swagger/index.html"
        print_status "  🔗 API Base: http://localhost:8081"
        print_status "  🐘 PostgreSQL: localhost:5432"
        print_status "  🔴 Redis: localhost:6379"
        
        print_status "📋 Showing service logs (Ctrl+C to stop)..."
        podman-compose -f docker-compose.dev.yml logs -f
    fi
}

# Main execution
main() {
    check_runtime
    generate_swagger
    start_dev_environment
}

# Run main function
main "$@"
