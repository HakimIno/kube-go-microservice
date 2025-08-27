#!/bin/bash

# Build single service
# Usage: ./scripts/build-service.sh <service-name>
# Example: ./scripts/build-service.sh user-service

if [ $# -eq 0 ]; then
    echo "Usage: $0 <service-name>"
    echo "Available services:"
    ls cmd/ | grep -v "^$"
    exit 1
fi

SERVICE_NAME=$1
SERVICE_PATH="cmd/${SERVICE_NAME}"

if [ ! -d "$SERVICE_PATH" ]; then
    echo "Error: Service '${SERVICE_NAME}' not found in cmd/"
    echo "Available services:"
    ls cmd/ | grep -v "^$"
    exit 1
fi

if [ ! -s "${SERVICE_PATH}/main.go" ]; then
    echo "Error: ${SERVICE_PATH}/main.go is empty or doesn't exist"
    exit 1
fi

echo "Building ${SERVICE_NAME}..."

# Create output directory
mkdir -p output/bin

# Build the service
if go build -o "output/bin/${SERVICE_NAME}" "./${SERVICE_PATH}"; then
    echo "✓ ${SERVICE_NAME} built successfully!"
    echo "Binary location: output/bin/${SERVICE_NAME}"
else
    echo "✗ Failed to build ${SERVICE_NAME}"
    exit 1
fi
