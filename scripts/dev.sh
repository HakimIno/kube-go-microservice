#!/bin/bash

# Development script for user-service
# This script automatically generates swagger docs and runs the server

set -e

echo "🚀 Starting development mode..."

# Function to cleanup on exit
cleanup() {
    echo "🛑 Stopping development server..."
    if [ ! -z "$PID" ]; then
        kill $PID 2>/dev/null || true
    fi
    exit 0
}

# Set up trap for cleanup
trap cleanup SIGINT SIGTERM

# Generate swagger docs
echo "📝 Generating swagger documentation..."
swag init -g cmd/user-service/main.go

# Function to watch for changes and regenerate swagger
watch_and_generate() {
    echo "👀 Watching for changes in Go files..."
    
    # Use fswatch if available, otherwise use inotifywait
    if command -v fswatch > /dev/null; then
        fswatch -o . -e ".*" -i ".*\.go$" | while read f; do
            echo "🔄 File changed, regenerating swagger docs..."
            swag init -g cmd/user-service/main.go
        done
    elif command -v inotifywait > /dev/null; then
        while inotifywait -r -e modify,create,delete . --exclude '\.git|tmp|docs'; do
            echo "🔄 File changed, regenerating swagger docs..."
            swag init -g cmd/user-service/main.go
        done
    else
        echo "⚠️  No file watcher found. Please install fswatch or inotifywait for auto-regeneration."
        echo "   For macOS: brew install fswatch"
        echo "   For Ubuntu: sudo apt-get install inotify-tools"
    fi
}

# Start file watcher in background
watch_and_generate &
WATCHER_PID=$!

# Run the server
echo "🏃 Starting user-service server..."
go run cmd/user-service/main.go &
PID=$!

# Wait for server to start
sleep 2

echo "✅ Development server is running!"
echo "📖 Swagger UI: http://localhost:8081/swagger/index.html"
echo "🔗 API Base: http://localhost:8081"
echo ""
echo "Press Ctrl+C to stop"

# Wait for background processes
wait $PID $WATCHER_PID
