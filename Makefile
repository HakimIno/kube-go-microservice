.PHONY: help build run dev clean swagger docker-check docker-install-mac docker-build docker-run docker-dev docker-prod docker-stop docker-clean docker-logs docker-restart podman-dev podman-prod podman-clean

# Default target
help:
	@echo "🚀 Available commands:"
	@echo ""
	@echo "📦 Build & Run:"
	@echo "  make build         - Build the user service"
	@echo "  make run           - Build and run the service"
	@echo "  make dev           - Run in development mode with hot reload"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  make docker-check      - Check Docker installation and status"
	@echo "  make docker-install-mac - Open Docker Desktop download page (macOS)"
	@echo "  make docker-dev        - Start Docker development environment"
	@echo "  make docker-prod       - Start Docker production environment"
	@echo "  make docker-stop       - Stop Docker containers"
	@echo "  make docker-clean      - Clean Docker containers and images"
	@echo "  make docker-logs       - Show Docker development logs"
	@echo "  make docker-restart    - Restart Docker development environment"
	@echo ""
	@echo "🟦 Podman:"
	@echo "  make podman-dev        - Start Podman development environment"
	@echo "  make podman-prod       - Start Podman production environment"
	@echo "  make podman-clean      - Clean Podman containers and images"
	@echo ""
	@echo "📚 Documentation:"
	@echo "  make swagger       - Generate swagger documentation"
	@echo ""
	@echo "🧹 Utilities:"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make deps          - Install dependencies"
	@echo "  make setup         - Setup development environment"
	@echo ""
	@echo "💡 Quick start:"
	@echo "  1. For Docker: make docker-check && make docker-dev"
	@echo "  2. For Podman: make podman-dev"
	@echo "  3. For local:  make dev"

# Build the application
build: swagger
	@echo "Building user-service..."
	go build -o output/bin/user-service cmd/user-service/main.go

# Run the application
run: build
	@echo "Running user-service..."
	./output/bin/user-service

# Development mode - auto rebuild and run
dev: swagger
	@echo "Starting development mode..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Installing air..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

# Generate swagger documentation
swagger:
	@echo "Generating swagger documentation..."
	swag init -g cmd/user-service/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf output/bin/*
	rm -rf docs/docs.go docs/swagger.json docs/swagger.yaml

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@latest

# Setup development environment
setup: deps swagger
	@echo "Development environment setup complete!"

# Watch for changes and auto-rebuild (alternative to air)
watch:
	@echo "Watching for changes..."
	@if command -v fswatch > /dev/null; then \
		fswatch -o . | xargs -n1 -I{} make swagger; \
	else \
		echo "fswatch not found. Please install it or use 'make dev' instead."; \
	fi

# Docker commands
docker-check:
	@echo "Checking Docker installation and status..."
	@if command -v docker > /dev/null; then \
		echo "✅ Docker is installed: $$(docker --version)"; \
		if docker info > /dev/null 2>&1; then \
			echo "✅ Docker is running"; \
		else \
			echo "❌ Docker is not running. Please start Docker first."; \
			echo "💡 Alternative: Run 'make podman-dev' to use Podman instead"; \
		fi; \
	else \
		echo "❌ Docker is not installed"; \
		echo "📝 Install from: https://docs.docker.com/get-docker/"; \
		echo "💡 Alternative: Run 'make podman-dev' to use Podman instead"; \
	fi

docker-install-mac:
	@echo "Opening Docker Desktop download page for macOS..."
	@open "https://docs.docker.com/desktop/install/mac-install/"

docker-build:
	@echo "Building Docker image..."
	docker build -f deployments/docker/Dockerfile.user-service -t user-service .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8081:8081 user-service

docker-dev:
	@echo "Starting Docker development environment..."
	./scripts/docker-dev.sh

docker-prod:
	@echo "Starting production Docker environment..."
	docker compose -f deployments/docker-compose/docker-compose.yml up --build

docker-stop:
	@echo "Stopping Docker containers..."
	docker compose -f deployments/docker-compose/docker-compose.yml down
	docker compose -f deployments/docker-compose/docker-compose.dev.yml down

docker-clean:
	@echo "Cleaning Docker containers and images..."
	docker compose -f deployments/docker-compose/docker-compose.yml down -v
	docker compose -f deployments/docker-compose/docker-compose.dev.yml down -v
	docker system prune -f

docker-logs:
	@echo "Showing Docker development environment logs..."
	docker compose -f deployments/docker-compose/docker-compose.dev.yml logs -f

docker-restart:
	@echo "Restarting Docker development environment..."
	docker compose -f deployments/docker-compose/docker-compose.dev.yml restart

# Podman commands
podman-dev:
	@echo "Starting Podman development environment..."
	./scripts/podman-dev.sh

podman-prod:
	@echo "Starting production Podman environment..."
	./scripts/podman-prod.sh

podman-clean:
	@echo "Cleaning Podman containers and images..."
	cd deployments/podman-compose && podman-compose -f docker-compose.yml down -v
	cd deployments/podman-compose && podman-compose -f docker-compose.dev.yml down -v
	podman system prune -f
