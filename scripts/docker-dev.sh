#!/bin/bash

# Docker development script
# This script runs the development environment using Docker Compose

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

print_status "🐳 Starting Docker development environment..."

# Function to cleanup on exit
cleanup() {
    print_status "🛑 Stopping Docker development environment..."
    docker compose -f deployments/docker-compose/docker-compose.dev.yml down
    exit 0
}

# Set up trap for cleanup
trap cleanup SIGINT SIGTERM

# Check if Docker is installed
check_docker_installed() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed!"
        echo ""
        print_status "📝 Installation instructions:"
        echo "  🍎 macOS: Download from https://docs.docker.com/desktop/install/mac-install/"
        echo "  🐧 Linux: Run 'curl -fsSL https://get.docker.com -o get-docker.sh && sh get-docker.sh'"
        echo "  🪟 Windows: Download from https://docs.docker.com/desktop/install/windows-install/"
        echo ""
        print_warning "💡 Alternative: You can use Podman instead by running 'make podman-dev'"
        exit 1
    fi
}

# Check if Docker is running
check_docker_running() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running!"
        echo ""
        print_status "📝 Attempting to start Docker..."
        
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # Try to start Docker Desktop on macOS
            if [ -d "/Applications/Docker.app" ]; then
                print_status "🍎 Starting Docker Desktop..."
                open /Applications/Docker.app
                
                # Wait for Docker to start (up to 60 seconds)
                print_status "⏳ Waiting for Docker to start..."
                for i in {1..60}; do
                    if docker info > /dev/null 2>&1; then
                        print_success "✅ Docker started successfully!"
                        return 0
                    fi
                    sleep 1
                    if [ $((i % 10)) -eq 0 ]; then
                        print_status "⏳ Still waiting... ($i/60 seconds)"
                    fi
                done
                
                print_error "❌ Docker failed to start within 60 seconds"
                print_status "📝 Manual steps:"
                echo "  1. Manually open Docker Desktop application"
                echo "  2. Wait for Docker to fully start"
                echo "  3. Run 'make docker-dev' again"
            else
                print_error "Docker Desktop not found at /Applications/Docker.app"
                print_status "📝 Please install Docker Desktop from:"
                echo "  🔗 https://docs.docker.com/desktop/install/mac-install/"
            fi
        elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
            print_status "🐧 Attempting to start Docker service..."
            if command -v systemctl &> /dev/null; then
                sudo systemctl start docker
            elif command -v service &> /dev/null; then
                sudo service docker start
            fi
            sleep 5
            if docker info > /dev/null 2>&1; then
                print_success "✅ Docker started successfully!"
                return 0
            fi
        elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
            print_status "🪟 Please start Docker Desktop manually on Windows"
        fi
        
        echo ""
        print_warning "💡 Alternative: You can use Podman instead by running 'make podman-dev'"
        exit 1
    fi
    print_success "Docker is running!"
}

# Check Docker and Docker Compose
check_docker_installed
check_docker_running

# Build and start development environment
print_status "🔨 Building development environment..."
docker compose -f deployments/docker-compose/docker-compose.dev.yml up --build -d

print_success "✅ Docker development environment is running!"
print_status "🌐 Services available at:"
print_status "  📖 Swagger UI: http://localhost:8081/swagger/index.html"
print_status "  🔗 API Base: http://localhost:8081"
print_status "  🐘 PostgreSQL: localhost:5432"
print_status "  🔴 Redis: localhost:6379"
echo ""
print_status "📋 Showing service logs (Ctrl+C to stop)..."

# Show logs for all services
docker compose -f deployments/docker-compose/docker-compose.dev.yml logs -f
