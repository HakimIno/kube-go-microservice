# Kube - Microservices Framework

A modern microservices framework built with Go and Hertz, featuring containerization support for both Docker and Podman.

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker or Podman
- PostgreSQL
- Redis

### Development Environment

#### Option 1: Using Docker (Recommended)
```bash
# Start development environment with Docker
make dev docker

# Or directly with the script
./scripts/dev.sh docker
```

#### Option 2: Using Podman
```bash
# Start development environment with Podman
make dev podman

# Or directly with the script
./scripts/dev.sh podman
```

#### Option 3: Local Development
```bash
# Run locally with hot reload
make dev
```

### Production Environment

#### Docker
```bash
# Start production environment
make prod docker

# Stop production environment
make prod stop docker

# View logs
make prod logs docker

# Restart services
make prod restart docker
```

#### Podman
```bash
# Start production environment
make prod podman

# Stop production environment
make prod stop podman

# View logs
make prod logs podman

# Restart services
make prod restart podman
```

## 🛠️ Available Commands

### Development
- `make dev` - Run locally with hot reload
- `make dev docker` - Start Docker development environment
- `make dev podman` - Start Podman development environment

### Production
- `make prod docker` - Start Docker production environment
- `make prod podman` - Start Podman production environment
- `make prod stop [docker|podman]` - Stop production environment
- `make prod logs [docker|podman]` - View production logs
- `make prod restart [docker|podman]` - Restart production environment

### Building
- `make build` - Build Go binary
- `make build docker` - Build Docker container image
- `make build podman` - Build Podman container image

### Utilities
- `make swagger` - Generate API documentation
- `make clean` - Clean build artifacts
- `make deps` - Install dependencies
- `make setup` - Setup development environment

## 📁 Scripts

### Unified Scripts

#### `./scripts/dev.sh` - Development Environment
```bash
# Usage: ./scripts/dev.sh [docker|podman] [--build]
./scripts/dev.sh                    # Use Docker (default)
./scripts/dev.sh podman            # Use Podman
./scripts/dev.sh docker --build    # Use Docker with build
./scripts/dev.sh podman --build    # Use Podman with build
```

#### `./scripts/prod.sh` - Production Environment
```bash
# Usage: ./scripts/prod.sh [docker|podman] {start|stop|logs|restart}
./scripts/prod.sh                    # Use Docker and start (default)
./scripts/prod.sh podman            # Use Podman and start
./scripts/prod.sh docker stop       # Use Docker and stop
./scripts/prod.sh podman logs       # Use Podman and show logs
```

#### `./scripts/build.sh` - Build Services
```bash
# Usage: ./scripts/build.sh <service-name> [docker|podman] [--push]
./scripts/build.sh user-service                    # Native Go build
./scripts/build.sh user-service docker            # Docker container build
./scripts/build.sh user-service podman            # Podman container build
./scripts/build.sh user-service docker --push     # Docker build and push
```

### Service Generation
```bash
# Generate a new service
./scripts/generate-service.sh <service-name> <port>
./scripts/generate-service.sh video-service 8082
```

## 🏗️ Project Structure

```
kube/
├── cmd/                    # Service entry points
│   └── user-service/      # User service
├── biz/                   # Business logic
│   ├── handler/           # HTTP handlers
│   ├── router/            # Route definitions
│   └── service/           # Business services
├── internal/              # Internal packages
│   ├── config/            # Configuration
│   ├── database/          # Database connection
│   └── middleware/        # HTTP middleware
├── pkg/                   # Public packages
│   ├── models/            # Data models
│   ├── handlers/          # Base handlers
│   ├── services/          # Base services
│   └── utils/             # Utility functions
├── deployments/           # Deployment configurations
│   ├── docker/            # Docker configurations
│   ├── docker-compose/    # Docker Compose files
│   └── podman-compose/    # Podman Compose files
└── scripts/               # Build and deployment scripts
```

## 🔧 Configuration

Copy the environment file and configure your settings:
```bash
cp env.example .env
# Edit .env with your configuration
```

## 📚 API Documentation

Once the service is running, access the Swagger UI:
- Development: http://localhost:8081/swagger/index.html
- Production: http://localhost:8081/swagger/index.html

## 🐳 Container Support

### Docker
- Full Docker Compose support
- Multi-stage builds
- Development and production configurations

### Podman
- Full Podman Compose support
- Rootless containers
- Compatible with Docker commands

## 🚀 Deployment

### Development
```bash
# Start with Docker
make dev docker

# Start with Podman
make dev podman
```

### Production
```bash
# Start with Docker
make prod docker

# Start with Podman
make prod podman
```

## 📝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
