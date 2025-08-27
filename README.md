# Video Streaming Microservices

ระบบ Video Streaming แบบ Microservices ที่สร้างด้วย Go และ Hertz Framework

## 🏗️ Architecture

ระบบประกอบด้วย 8 microservices:

1. **User Service** (Port: 8081) - จัดการข้อมูลผู้ใช้
2. **Video Upload Service** (Port: 8082) - อัปโหลดวิดีโอ
3. **Video Processing Service** (Port: 8083) - ประมวลผลวิดีโอ
4. **Metadata Service** (Port: 8084) - จัดการข้อมูลวิดีโอ
5. **Streaming Service** (Port: 8085) - สตรีมมิ่งวิดีโอ
6. **Search Service** (Port: 8086) - ค้นหาวิดีโอ
7. **Recommendation Service** (Port: 8087) - แนะนำวิดีโอ
8. **Engagement Service** (Port: 8088) - จัดการปฏิสัมพันธ์

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL
- Redis

### Running with Docker Compose

```bash
# Clone repository
git clone <repository-url>
cd kube

# Start all services
cd deployments/docker-compose
docker-compose up -d

# Check services
docker-compose ps
```

### Running Locally

```bash
# Start PostgreSQL & Redis
docker-compose -f deployments/docker-compose/docker-compose.yml up postgres redis -d

# Run User Service
./scripts/run-user-service.sh

# Or run manually
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

go run cmd/user-service/main.go
```

## 📁 Project Structure

```
kube/
├── cmd/                           # Entry points สำหรับแต่ละ service
│   ├── user-service/             # ระบบจัดการผู้ใช้
│   ├── video-upload-service/     # ระบบอัปโหลดวิดีโอ
│   ├── video-processing-service/ # ระบบประมวลผลวิดีโอ
│   ├── metadata-service/         # ระบบจัดการข้อมูลวิดีโอ
│   ├── streaming-service/        # ระบบสตรีมมิ่ง
│   ├── search-service/           # ระบบค้นหา
│   ├── recommendation-service/   # ระบบแนะนำวิดีโอ
│   └── engagement-service/       # ระบบจัดการปฏิสัมพันธ์
├── internal/                     # Shared internal packages
│   ├── config/                  # Configuration management
│   ├── database/                # Database connections
│   ├── middleware/              # Shared middleware
│   ├── utils/                   # Shared utilities
│   ├── auth/                    # Authentication & Authorization
│   ├── messaging/               # Message queue (Redis, RabbitMQ)
│   └── storage/                 # File storage (S3, local)
├── pkg/                         # Public packages ที่ service อื่นใช้ได้
│   ├── models/                  # Shared data models
│   ├── constants/               # Shared constants
│   ├── errors/                  # Shared error definitions
│   ├── events/                  # Event definitions
│   └── proto/                   # Protocol buffer definitions
├── services/                    # Business logic ของแต่ละ service
│   ├── user/                    # User service logic
│   ├── video-upload/            # Video upload logic
│   ├── video-processing/        # Video processing logic
│   ├── metadata/                # Metadata management logic
│   ├── streaming/               # Streaming logic
│   ├── search/                  # Search logic
│   ├── recommendation/          # Recommendation logic
│   └── engagement/              # Engagement logic
├── api/                         # API definitions
│   ├── proto/                   # Protocol buffer files
│   └── openapi/                 # OpenAPI specifications
├── deployments/                 # Infrastructure configs
│   ├── docker/                  # Docker files
│   ├── kubernetes/              # K8s manifests
│   └── docker-compose/          # Local development
├── scripts/                     # Build, deploy scripts
├── docs/                        # Documentation
└── tools/                       # Development tools
```

## 🔧 Technology Stack

- **Framework**: Hertz (CloudWeGo)
- **Database**: PostgreSQL
- **Cache**: Redis
- **Message Queue**: Redis/RabbitMQ
- **Container**: Docker
- **Orchestration**: Kubernetes
- **API Documentation**: OpenAPI/Swagger

## 📚 Documentation

- [User Service](./docs/USER_SERVICE.md)
- [API Documentation](./api/openapi/)
- [Deployment Guide](./deployments/)

## 🧪 Testing

### Test User Service

```bash
# Health check
curl http://localhost:8081/health

# Register user
curl -X POST http://localhost:8081/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'

# Login
curl -X POST http://localhost:8081/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## 🚀 Development

### Adding New Service

1. สร้างโฟลเดอร์ใน `cmd/` และ `services/`
2. สร้าง main.go ใน `cmd/[service-name]/`
3. สร้าง business logic ใน `services/[service-name]/`
4. เพิ่ม Dockerfile ใน `deployments/docker/`
5. อัปเดต docker-compose.yml

### Code Style

- ใช้ Go modules
- ใช้ Hertz framework
- ใช้ GORM สำหรับ database
- ใช้ bcrypt สำหรับ password hashing
- ใช้ JWT สำหรับ authentication

## 📄 License

MIT License

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request # kube-go-microservice
