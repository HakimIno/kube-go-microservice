#!/bin/bash

# Unified Build Script
# Supports both native Go build and container builds with Docker/Podman
# Usage: ./scripts/build.sh <service-name> [docker|podman] [--push] [--generate]

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

# Default values
BUILD_TYPE="native"
CONTAINER_RUNTIME="docker"
PUSH_FLAG=""
GENERATE_FLAG=""
SERVICE_PORT="8081"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        docker|podman)
            BUILD_TYPE="container"
            CONTAINER_RUNTIME="$1"
            shift
            ;;
        --push)
            PUSH_FLAG="--push"
            shift
            ;;
        --generate)
            GENERATE_FLAG="--generate"
            shift
            ;;
        --port)
            SERVICE_PORT="$2"
            shift 2
            ;;
        -h|--help)
            echo "Usage: $0 <service-name> [docker|podman] [--push] [--generate] [--port <port>]"
            echo ""
            echo "Arguments:"
            echo "  service-name    Name of the service to build"
            echo ""
            echo "Options:"
            echo "  docker|podman   Build container image (default: docker)"
            echo "  --push          Push image to registry (container builds only)"
            echo "  --generate      Generate service files before building"
            echo "  --port <port>   Port for the service (default: 8081)"
            echo "  -h, --help      Show this help message"
            echo ""
            echo "Examples:"
            echo "  $0 user-service                    # Native Go build"
            echo "  $0 user-service docker            # Docker container build"
            echo "  $0 user-service podman            # Podman container build"
            echo "  $0 user-service docker --push     # Docker build and push"
            echo "  $0 video-service --generate       # Generate and build video-service"
            echo "  $0 api-service --generate --port 8082  # Generate with custom port"
            exit 0
            ;;
        *)
            if [ -z "$SERVICE_NAME" ]; then
                SERVICE_NAME="$1"
            else
                print_error "Unknown option: $1"
                echo "Use -h or --help for usage information"
                exit 1
            fi
            shift
            ;;
    esac
done

# Check if service name is provided
if [ -z "$SERVICE_NAME" ]; then
    print_error "Service name is required"
    echo "Usage: $0 <service-name> [docker|podman] [--push] [--generate] [--port <port>]"
    echo "Use -h or --help for more information"
    exit 1
fi

SERVICE_PATH="cmd/${SERVICE_NAME}"

# Generate service if requested
generate_service() {
    if [ "$GENERATE_FLAG" = "--generate" ]; then
        print_status "🔧 Generating service: $SERVICE_NAME on port $SERVICE_PORT..."
        
        # Check if generate-service.sh exists
        if [ ! -f "scripts/generate-service.sh" ]; then
            print_error "generate-service.sh not found!"
            exit 1
        fi
        
        # Generate the service
        if ./scripts/generate-service.sh "$SERVICE_NAME" "$SERVICE_PORT"; then
            print_success "✓ Service $SERVICE_NAME generated successfully!"
        else
            print_error "✗ Failed to generate service $SERVICE_NAME"
            exit 1
        fi
    fi
}

# Validate service exists
validate_service() {
    if [ ! -d "$SERVICE_PATH" ]; then
        if [ "$GENERATE_FLAG" = "--generate" ]; then
            print_error "Service '${SERVICE_NAME}' was not generated properly"
            exit 1
        else
            print_error "Service '${SERVICE_NAME}' not found in cmd/"
            echo "Available services:"
            ls cmd/ | grep -v "^$"
            echo ""
            echo "Tip: Use --generate flag to create a new service"
            exit 1
        fi
    fi

    if [ ! -s "${SERVICE_PATH}/main.go" ]; then
        print_error "Error: ${SERVICE_PATH}/main.go is empty or doesn't exist"
        exit 1
    fi
}

# Native Go build
build_native() {
    print_status "🔨 Building ${SERVICE_NAME} with Go..."
    
    # Create output directory
    mkdir -p output/bin
    
    # Build the service
    if go build -o "output/bin/${SERVICE_NAME}" "./${SERVICE_PATH}"; then
        print_success "✓ ${SERVICE_NAME} built successfully!"
        print_status "Binary location: output/bin/${SERVICE_NAME}"
    else
        print_error "✗ Failed to build ${SERVICE_NAME}"
        exit 1
    fi
}

# Container build
build_container() {
    print_status "🐳 Building ${SERVICE_NAME} with ${CONTAINER_RUNTIME}..."
    
    # Check if container runtime is installed
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
            print_error "Podman is not installed!"
            exit 1
        fi
        print_success "Podman is available: $(podman --version)"
    fi
    
    # Check Dockerfile exists
    local dockerfile_path="deployments/docker/Dockerfile.${SERVICE_NAME}"
    if [ ! -f "$dockerfile_path" ]; then
        print_error "Dockerfile not found: $dockerfile_path"
        if [ "$GENERATE_FLAG" = "--generate" ]; then
            print_status "Tip: Dockerfile should be generated with the service"
        fi
        exit 1
    fi
    
    print_status "Using Dockerfile: $dockerfile_path"
    
    # Build the image
    local image_name="${SERVICE_NAME}:latest"
    if [ "$CONTAINER_RUNTIME" = "docker" ]; then
        docker build -f "$dockerfile_path" -t "$image_name" .
    else
        podman build -f "$dockerfile_path" -t "$image_name" .
    fi
    
    print_success "✓ Service $SERVICE_NAME built successfully!"
    print_status "Image: $image_name"
    
    # Push if requested
    if [ "$PUSH_FLAG" = "--push" ]; then
        print_status "📤 Pushing image to registry..."
        if [ "$CONTAINER_RUNTIME" = "docker" ]; then
            docker push "$image_name"
        else
            podman push "$image_name"
        fi
        print_success "✓ Image pushed successfully!"
    fi
}

# Show available services
show_available_services() {
    echo "Available services:"
    ls cmd/ | grep -v "^$" | while read service; do
        echo "  - $service"
    done
}

# Main execution
main() {
    print_status "🚀 Starting build process for: $SERVICE_NAME"
    print_status "Build type: $BUILD_TYPE"
    
    if [ "$BUILD_TYPE" = "container" ]; then
        print_status "Container runtime: $CONTAINER_RUNTIME"
    fi
    
    if [ "$GENERATE_FLAG" = "--generate" ]; then
        print_status "Service generation: enabled"
        print_status "Service port: $SERVICE_PORT"
    fi
    
    # Generate service first if requested
    generate_service
    
    # Validate service exists
    validate_service
    
    # Build the service
    if [ "$BUILD_TYPE" = "native" ]; then
        build_native
    else
        build_container
    fi
    
    print_success "✅ Build completed successfully!"
    
    if [ "$BUILD_TYPE" = "native" ]; then
        print_status "💡 To run the service: ./output/bin/${SERVICE_NAME}"
    else
        print_status "💡 To run the container: ${CONTAINER_RUNTIME} run -p ${SERVICE_PORT}:${SERVICE_PORT} ${SERVICE_NAME}:latest"
    fi
}

# Run main function
main "$@"
