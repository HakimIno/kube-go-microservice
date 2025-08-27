.PHONY: build run dev clean swagger docker-build docker-run docker-dev

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
