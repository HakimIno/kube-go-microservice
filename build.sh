#!/bin/bash

# Build all services
echo "Building all services..."

# Create output directories
mkdir -p output/bin

# Build user-service
echo "Building user-service..."
go build -o output/bin/user-service ./cmd/user-service

# Build video-upload-service (if main.go exists and has content)
if [ -s "cmd/video-upload-service/main.go" ]; then
    echo "Building video-upload-service..."
    if go build -o output/bin/video-upload-service ./cmd/video-upload-service 2>/dev/null; then
        echo "✓ video-upload-service built successfully"
    else
        echo "✗ Failed to build video-upload-service (main.go may be empty or invalid)"
    fi
fi

# Build metadata-service (if main.go exists and has content)
if [ -s "cmd/metadata-service/main.go" ]; then
    echo "Building metadata-service..."
    if go build -o output/bin/metadata-service ./cmd/metadata-service 2>/dev/null; then
        echo "✓ metadata-service built successfully"
    else
        echo "✗ Failed to build metadata-service (main.go may be empty or invalid)"
    fi
fi

echo "Build completed! Binaries are in output/bin/"