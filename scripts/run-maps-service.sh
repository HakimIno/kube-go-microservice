#!/bin/bash

# Run Maps Service Development Script
echo "Starting Maps Service..."

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=video_streaming_dev
export DB_SSLMODE=disable
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=
export REDIS_DB=0
export JWT_SECRET=dev-secret-key-change-in-production
export JWT_EXPIRES_IN=24
export HERE_API_KEY=WJxd-f213J08HqNIGMgM3dfQaRTwbe02D8unFbc_ge4
export ENV=development

# Run the maps service
go run cmd/maps-service/main.go
