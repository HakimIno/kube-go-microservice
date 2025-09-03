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
	@echo "🐳 Container Commands:"
	@echo "  make dev docker    - Start development environment with Docker"
	@echo "  make dev podman    - Start development environment with Podman"
	@echo "  make prod docker   - Start production environment with Docker"
	@echo "  make prod podman   - Start production environment with Podman"
	@echo "  make build docker  - Build container image with Docker"
	@echo "  make build podman  - Build container image with Podman"
	@echo ""
	@echo "🔧 Service Generation:"
	@echo "  make generate-service SERVICE=<name> PORT=<port>"
	@echo "  make build-generate SERVICE=<name> PORT=<port> RUNTIME=<docker|podman>"
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
	@echo "  1. For Docker: make dev docker"
	@echo "  2. For Podman: make dev podman"
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

# Unified container commands
docker-dev:
	@echo "Starting Docker development environment..."
	./scripts/dev.sh docker

podman-dev:
	@echo "Starting Podman development environment..."
	./scripts/dev.sh podman

docker-prod:
	@echo "Starting Docker production environment..."
	./scripts/prod.sh docker start

podman-prod:
	@echo "Starting Podman production environment..."
	./scripts/prod.sh podman start

docker-stop:
	@echo "Stopping Docker production environment..."
	./scripts/prod.sh docker stop

podman-stop:
	@echo "Stopping Podman production environment..."
	./scripts/prod.sh podman stop

docker-logs:
	@echo "Showing Docker production logs..."
	./scripts/prod.sh docker logs

podman-logs:
	@echo "Showing Podman production logs..."
	./scripts/prod.sh podman logs

docker-restart:
	@echo "Restarting Docker production environment..."
	./scripts/prod.sh docker restart

podman-restart:
	@echo "Restarting Podman production environment..."
	./scripts/prod.sh podman restart

# Build commands
build-docker:
	@echo "Building Docker container image..."
	./scripts/build.sh user-service docker

build-podman:
	@echo "Building Podman container image..."
	./scripts/build.sh user-service podman

# Generate and build new services
generate-service:
	@echo "Usage: make generate-service SERVICE=<name> PORT=<port>"
	@echo "Example: make generate-service SERVICE=video-service PORT=8082"
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: SERVICE parameter is required"; \
		echo "Example: make generate-service SERVICE=video-service PORT=8082"; \
		exit 1; \
	fi
	@echo "Generating service: $(SERVICE) on port $(PORT:-=8081)"
	./scripts/generate-service.sh "$(SERVICE)" "$(PORT:-=8081)"

# Build with generation
build-generate:
	@echo "Usage: make build-generate SERVICE=<name> PORT=<port> [RUNTIME=docker|podman]"
	@echo "Example: make build-generate SERVICE=video-service PORT=8082 RUNTIME=docker"
	@if [ -z "$(SERVICE)" ]; then \
		echo "Error: SERVICE parameter is required"; \
		echo "Example: make build-generate SERVICE=video-service PORT=8082 RUNTIME=docker"; \
		exit 1; \
	fi
	@echo "Generating and building service: $(SERVICE) on port $(PORT:-=8081) with $(RUNTIME:-=docker)"
	./scripts/build.sh "$(SERVICE)" "$(RUNTIME:-=docker)" --generate --port "$(PORT:-=8081)"

# Legacy aliases for backward compatibility
docker-check: docker-dev
	@echo "Note: docker-check is deprecated, use 'make docker-dev' instead"

docker-install-mac:
	@echo "Opening Docker Desktop download page for macOS..."
	@open "https://docs.docker.com/desktop/install/mac-install/"

docker-build: build-docker

docker-run: docker-dev

docker-clean:
	@echo "Cleaning Docker containers and images..."
	docker compose -f deployments/docker-compose/docker-compose.yml down -v
	docker compose -f deployments/docker-compose/docker-compose.dev.yml down -v
	docker system prune -f

podman-clean:
	@echo "Cleaning Podman containers and images..."
	cd deployments/podman-compose && podman-compose -f docker-compose.yml down -v
	cd deployments/podman-compose && podman-compose -f docker-compose.dev.yml down -v
	podman system prune -f
