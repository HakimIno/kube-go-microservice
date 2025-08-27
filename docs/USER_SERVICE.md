# User Service

User Service เป็น microservice สำหรับจัดการข้อมูลผู้ใช้ในระบบ Video Streaming

## Features

- ✅ User Registration
- ✅ User Login
- ✅ Get User Profile
- ✅ Update User Profile
- ✅ Delete User Account
- ✅ Password Hashing (bcrypt)
- ✅ Database Auto Migration

## API Endpoints

### Health Check
```
GET /health
```

### User Registration
```
POST /api/v1/users/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

### User Login
```
POST /api/v1/users/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

### Get User Profile
```
GET /api/v1/users/{id}
```

### Update User Profile
```
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Smith",
  "avatar": "https://example.com/avatar.jpg"
}
```

### Delete User Account
```
DELETE /api/v1/users/{id}
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_USER | postgres | PostgreSQL username |
| DB_PASSWORD | password | PostgreSQL password |
| DB_NAME | video_streaming | PostgreSQL database name |
| DB_SSLMODE | disable | PostgreSQL SSL mode |
| REDIS_HOST | localhost | Redis host |
| REDIS_PORT | 6379 | Redis port |
| JWT_SECRET | your-secret-key | JWT secret key |
| JWT_EXPIRES_IN | 24 | JWT expiration time (hours) |

## Running Locally

### Prerequisites
- Go 1.23+
- PostgreSQL
- Redis

### Using Script
```bash
./scripts/run-user-service.sh
```

### Using Docker Compose
```bash
cd deployments/docker-compose
docker-compose up user-service
```

### Manual Run
```bash
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

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    avatar VARCHAR(500),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

## Testing

### Using curl

#### Register User
```bash
curl -X POST http://localhost:8081/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

#### Login User
```bash
curl -X POST http://localhost:8081/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

#### Get User Profile
```bash
curl -X GET http://localhost:8081/api/v1/users/1
```

## Port

User Service ทำงานบน port **8081** 