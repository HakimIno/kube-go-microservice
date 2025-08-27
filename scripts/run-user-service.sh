#!/bin/bash

echo "Starting User Service..."

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=video_streaming
export DB_SSLMODE=disable
export REDIS_HOST=localhost
export REDIS_PORT=6379
export JWT_SECRET=your-secret-key
export JWT_EXPIRES_IN=24

# Run the service
go run cmd/user-service/main.go 