#!/bin/bash

# Development script for running services with hot reload

echo "🚀 Starting development environment..."

# Function to start development services
start_dev() {
    echo "Starting services with hot reload..."
    cd deployments/podman-compose
    podman-compose -f docker-compose.dev.yml up --build
}

# Function to stop services
stop_dev() {
    echo "Stopping development services..."
    cd deployments/podman-compose
    podman-compose -f docker-compose.dev.yml down
}

# Function to restart a specific service
restart_service() {
    if [ -z "$1" ]; then
        echo "Usage: $0 restart <service-name>"
        echo "Available services: maps-service, user-service"
        exit 1
    fi
    
    echo "Restarting $1..."
    cd deployments/podman-compose
    podman-compose -f docker-compose.dev.yml restart $1
}

# Function to view logs
logs() {
    if [ -z "$1" ]; then
        echo "Usage: $0 logs <service-name>"
        echo "Available services: maps-service, user-service, postgres, redis"
        exit 1
    fi
    
    cd deployments/podman-compose
    podman-compose -f docker-compose.dev.yml logs -f $1
}

# Main script logic
case "$1" in
    start)
        start_dev
        ;;
    stop)
        stop_dev
        ;;
    restart)
        restart_service $2
        ;;
    logs)
        logs $2
        ;;
    *)
        echo "Usage: $0 {start|stop|restart <service>|logs <service>}"
        echo ""
        echo "Commands:"
        echo "  start                 - Start all development services with hot reload"
        echo "  stop                  - Stop all development services"
        echo "  restart <service>     - Restart a specific service"
        echo "  logs <service>        - View logs for a specific service"
        echo ""
        echo "Available services: maps-service, user-service, postgres, redis"
        exit 1
        ;;
esac