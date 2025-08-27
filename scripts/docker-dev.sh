#!/bin/bash

# Docker development script
# This script runs the development environment using Docker Compose

set -e

echo "🐳 Starting Docker development environment..."

# Function to cleanup on exit
cleanup() {
    echo "🛑 Stopping Docker development environment..."
    docker compose -f deployments/docker-compose/docker-compose.dev.yml down
    exit 0
}

# Set up trap for cleanup
trap cleanup SIGINT SIGTERM

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Build and start development environment
echo "🔨 Building development environment..."
docker compose -f deployments/docker-compose/docker-compose.dev.yml up --build

echo "✅ Docker development environment is running!"
echo "📖 Swagger UI: http://localhost:8081/swagger/index.html"
echo "🔗 API Base: http://localhost:8081"
echo ""
echo "Press Ctrl+C to stop"

# Wait for user to stop
wait
