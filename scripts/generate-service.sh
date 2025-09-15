#!/bin/bash

# Generate Service Script for New Architecture
# Usage: ./scripts/generate-service-new.sh <service-name> <port>

if [ $# -ne 2 ]; then
    echo "Usage: $0 <service-name> <port>"
    echo "Example: $0 notification-service 8083"
    exit 1
fi

SERVICE_NAME=$1
PORT=$2
SERVICE_DIR="services/$SERVICE_NAME"
CMD_DIR="$SERVICE_DIR/cmd"
INTERNAL_DIR="$SERVICE_DIR/internal"
CONFIG_DIR="$SERVICE_DIR/config"
MIGRATIONS_DIR="$SERVICE_DIR/migrations"

# Convert service name to PascalCase for model names
pascal_case() {
    echo "$1" | awk -F'-' '{
        result = ""
        for(i=1; i<=NF; i++) {
            word = $i
            first_char = toupper(substr(word, 1, 1))
            rest = substr(word, 2)
            result = result first_char rest
        }
        print result
    }'
}

MODEL_NAME=$(pascal_case "$SERVICE_NAME")

echo "🚀 Generating service: $SERVICE_NAME on port $PORT"
echo "📁 Model name: $MODEL_NAME"

# Create directories
mkdir -p "$CMD_DIR"
mkdir -p "$INTERNAL_DIR"/{handler,service,repository,model,router}
mkdir -p "$CONFIG_DIR"
mkdir -p "$MIGRATIONS_DIR"

# Create main.go
cat > "$CMD_DIR/main.go" << EOF
package main

import (
	"kube/services/$SERVICE_NAME/internal/router"
	"kube/shared/internal/config"
	"kube/shared/internal/database"
	"kube/shared/internal/middleware"
	"kube/services/$SERVICE_NAME/internal/model"
	"kube/shared/pkg/server"
	"time"
)

// @title $MODEL_NAME Service API
// @version 1.0
// @description This is a $SERVICE_NAME API built with Hertz framework.

// @contact.name API Support
// @contact.url https://github.com/your-username/kube
// @contact.email support@example.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:$PORT
// @BasePath /
// @schemes http https

func main() {
	cfg := config.Load()
	db := database.Init(cfg.Database)

	// Run migrations
	if err := db.AutoMigrate(&model.$MODEL_NAME{}); err != nil {
		middleware.LogError("Failed to migrate database", err)
		panic("Database migration failed")
	}

	serverConfig := server.ServerConfig{
		Port:         "$PORT",
		ServiceName:  "$SERVICE_NAME",
		SwaggerURL:   "http://localhost:$PORT",
		RateLimit:    100,
		RateDuration: time.Minute,
	}

	srv := server.NewServer(serverConfig)

	// Register routes
	router.RegisterRoutes(srv.Hertz, db)

	srv.Start()
}
EOF

# Create model
cat > "$INTERNAL_DIR/model/$SERVICE_NAME.go" << EOF
package model

import (
	"time"
	"gorm.io/gorm"
)

// $MODEL_NAME represents a $SERVICE_NAME in the system
type $MODEL_NAME struct {
	ID        string         \`json:"id" gorm:"primaryKey"\`
	Name      string         \`json:"name" gorm:"not null"\`
	Status    string         \`json:"status" gorm:"default:'active'"\`
	CreatedAt time.Time      \`json:"created_at"\`
	UpdatedAt time.Time      \`json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`json:"deleted_at" gorm:"index"\`
}

// ${MODEL_NAME}CreateRequest represents the request to create a new $SERVICE_NAME
type ${MODEL_NAME}CreateRequest struct {
	Name   string \`json:"name" binding:"required"\`
	Status string \`json:"status"\`
}

// ${MODEL_NAME}UpdateRequest represents the request to update a $SERVICE_NAME
type ${MODEL_NAME}UpdateRequest struct {
	Name   string \`json:"name"\`
	Status string \`json:"status"\`
}

// ${MODEL_NAME}Response represents the response for $SERVICE_NAME data
type ${MODEL_NAME}Response struct {
	ID        string    \`json:"id"\`
	Name      string    \`json:"name"\`
	Status    string    \`json:"status"\`
	CreatedAt time.Time \`json:"created_at"\`
	UpdatedAt time.Time \`json:"updated_at"\`
}
EOF

# Create repository
cat > "$INTERNAL_DIR/repository/$SERVICE_NAME.go" << EOF
package repository

import (
	"kube/services/$SERVICE_NAME/internal/model"
	"gorm.io/gorm"
)

type ${MODEL_NAME}Repository struct {
	db *gorm.DB
}

func New${MODEL_NAME}Repository(db *gorm.DB) *${MODEL_NAME}Repository {
	return &${MODEL_NAME}Repository{db: db}
}

func (r *${MODEL_NAME}Repository) Create($(echo $SERVICE_NAME | tr '-' '_') *model.$MODEL_NAME) error {
	return r.db.Create($(echo $SERVICE_NAME | tr '-' '_')).Error
}

func (r *${MODEL_NAME}Repository) GetByID(id string) (*model.$MODEL_NAME, error) {
	var $(echo $SERVICE_NAME | tr '-' '_') model.$MODEL_NAME
	err := r.db.Where("id = ?", id).First(&$(echo $SERVICE_NAME | tr '-' '_')).Error
	return &$(echo $SERVICE_NAME | tr '-' '_'), err
}

func (r *${MODEL_NAME}Repository) Update($(echo $SERVICE_NAME | tr '-' '_') *model.$MODEL_NAME) error {
	return r.db.Save($(echo $SERVICE_NAME | tr '-' '_')).Error
}

func (r *${MODEL_NAME}Repository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.$MODEL_NAME{}).Error
}

func (r *${MODEL_NAME}Repository) List(limit, offset int) ([]*model.$MODEL_NAME, error) {
	var $(echo $SERVICE_NAME | tr '-' '_')s []*model.$MODEL_NAME
	err := r.db.Limit(limit).Offset(offset).Find(&$(echo $SERVICE_NAME | tr '-' '_')s).Error
	return $(echo $SERVICE_NAME | tr '-' '_')s, err
}
EOF

# Create service
cat > "$INTERNAL_DIR/service/$SERVICE_NAME.go" << EOF
package service

import (
	"errors"
	"kube/services/$SERVICE_NAME/internal/model"
	"kube/services/$SERVICE_NAME/internal/repository"
	"kube/shared/pkg/errors"
	"kube/shared/pkg/services"
	"gorm.io/gorm"
)

type ${MODEL_NAME}Service struct {
	*services.BaseService
	repo *repository.${MODEL_NAME}Repository
}

func New${MODEL_NAME}Service(db *gorm.DB) *${MODEL_NAME}Service {
	return &${MODEL_NAME}Service{
		BaseService: services.NewBaseService(db),
		repo:        repository.New${MODEL_NAME}Repository(db),
	}
}

func (s *${MODEL_NAME}Service) Create${MODEL_NAME}(req *model.${MODEL_NAME}CreateRequest) (*model.${MODEL_NAME}Response, error) {
	$(echo $SERVICE_NAME | tr '-' '_') := &model.$MODEL_NAME{
		Name:   req.Name,
		Status: req.Status,
	}

	if $(echo $SERVICE_NAME | tr '-' '_').Status == "" {
		$(echo $SERVICE_NAME | tr '-' '_').Status = "active"
	}

	if err := s.repo.Create($(echo $SERVICE_NAME | tr '-' '_')); err != nil {
		return nil, apperrors.NewInternalServerError("Failed to create $SERVICE_NAME")
	}

	return s.$(echo $SERVICE_NAME | tr '-' '_')ToResponse($(echo $SERVICE_NAME | tr '-' '_')), nil
}

func (s *${MODEL_NAME}Service) Get${MODEL_NAME}(id string) (*model.${MODEL_NAME}Response, error) {
	$(echo $SERVICE_NAME | tr '-' '_'), err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("$SERVICE_NAME not found")
		}
		return nil, apperrors.NewInternalServerError("Failed to get $SERVICE_NAME")
	}

	return s.$(echo $SERVICE_NAME | tr '-' '_')ToResponse($(echo $SERVICE_NAME | tr '-' '_')), nil
}

func (s *${MODEL_NAME}Service) Update${MODEL_NAME}(id string, req *model.${MODEL_NAME}UpdateRequest) (*model.${MODEL_NAME}Response, error) {
	$(echo $SERVICE_NAME | tr '-' '_'), err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("$SERVICE_NAME not found")
		}
		return nil, apperrors.NewInternalServerError("Failed to get $SERVICE_NAME")
	}

	if req.Name != "" {
		$(echo $SERVICE_NAME | tr '-' '_').Name = req.Name
	}

	if req.Status != "" {
		$(echo $SERVICE_NAME | tr '-' '_').Status = req.Status
	}

	if err := s.repo.Update($(echo $SERVICE_NAME | tr '-' '_')); err != nil {
		return nil, apperrors.NewInternalServerError("Failed to update $SERVICE_NAME")
	}

	return s.$(echo $SERVICE_NAME | tr '-' '_')ToResponse($(echo $SERVICE_NAME | tr '-' '_')), nil
}

func (s *${MODEL_NAME}Service) Delete${MODEL_NAME}(id string) error {
	$(echo $SERVICE_NAME | tr '-' '_'), err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError("$SERVICE_NAME not found")
		}
		return apperrors.NewInternalServerError("Failed to get $SERVICE_NAME")
	}

	if err := s.repo.Delete(id); err != nil {
		return apperrors.NewInternalServerError("Failed to delete $SERVICE_NAME")
	}

	return nil
}

func (s *${MODEL_NAME}Service) List${MODEL_NAME}s(limit, offset int) ([]*model.${MODEL_NAME}Response, error) {
	$(echo $SERVICE_NAME | tr '-' '_')s, err := s.repo.List(limit, offset)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to list $SERVICE_NAMEs")
	}

	var responses []*model.${MODEL_NAME}Response
	for _, $(echo $SERVICE_NAME | tr '-' '_') := range $(echo $SERVICE_NAME | tr '-' '_')s {
		responses = append(reses, s.$(echo $SERVICE_NAME | tr '-' '_')ToResponse($(echo $SERVICE_NAME | tr '-' '_')))
	}

	return responses, nil
}

func (s *${MODEL_NAME}Service) $(echo $SERVICE_NAME | tr '-' '_')ToResponse($(echo $SERVICE_NAME | tr '-' '_') *model.$MODEL_NAME) *model.${MODEL_NAME}Response {
	return &model.${MODEL_NAME}Response{
		ID:        $(echo $SERVICE_NAME | tr '-' '_').ID,
		Name:      $(echo $SERVICE_NAME | tr '-' '_').Name,
		Status:    $(echo $SERVICE_NAME | tr '-' '_').Status,
		CreatedAt: $(echo $SERVICE_NAME | tr '-' '_').CreatedAt,
		UpdatedAt: $(echo $SERVICE_NAME | tr '-' '_').UpdatedAt,
	}
}
EOF

# Create handler
cat > "$INTERNAL_DIR/handler/$SERVICE_NAME.go" << EOF
package handler

import (
	"strconv"
	"kube/services/$SERVICE_NAME/internal/model"
	"kube/services/$SERVICE_NAME/internal/service"
	"kube/shared/pkg/errors"
	"github.com/cloudwego/hertz/pkg/app"
)

type ${MODEL_NAME}Handler struct {
	service *service.${MODEL_NAME}Service
}

func New${MODEL_NAME}Handler(service *service.${MODEL_NAME}Service) *${MODEL_NAME}Handler {
	return &${MODEL_NAME}Handler{service: service}
}

func (h *${MODEL_NAME}Handler) Create${MODEL_NAME}(ctx context.Context, c *app.RequestContext) {
	var req model.${MODEL_NAME}CreateRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	$(echo $SERVICE_NAME | tr '-' '_'), err := h.service.Create${MODEL_NAME}(&req)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(201, map[string]interface{}{
		"success": true,
		"data":    $(echo $SERVICE_NAME | tr '-' '_'),
		"message": "$SERVICE_NAME created successfully",
	})
}

func (h *${MODEL_NAME}Handler) Get${MODEL_NAME}(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("ID is required"))
		return
	}

	$(echo $SERVICE_NAME | tr '-' '_'), err := h.service.Get${MODEL_NAME}(id)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    $(echo $SERVICE_NAME | tr '-' '_'),
	})
}

func (h *${MODEL_NAME}Handler) Update${MODEL_NAME}(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("ID is required"))
		return
	}

	var req model.${MODEL_NAME}UpdateRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	$(echo $SERVICE_NAME | tr '-' '_'), err := h.service.Update${MODEL_NAME}(id, &req)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    $(echo $SERVICE_NAME | tr '-' '_'),
		"message": "$SERVICE_NAME updated successfully",
	})
}

func (h *${MODEL_NAME}Handler) Delete${MODEL_NAME}(ctx context.Context, c *app.RequestContext) {
	id := c.Param("id")
	if id == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("ID is required"))
		return
	}

	err := h.service.Delete${MODEL_NAME}(id)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "$SERVICE_NAME deleted successfully",
	})
}

func (h *${MODEL_NAME}Handler) List${MODEL_NAME}s(ctx context.Context, c *app.RequestContext) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	$(echo $SERVICE_NAME | tr '-' '_')s, err := h.service.List${MODEL_NAME}s(limit, offset)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    $(echo $SERVICE_NAME | tr '-' '_')s,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
		},
	})
}
EOF

# Create router
cat > "$INTERNAL_DIR/router/routes.go" << EOF
package router

import (
	"kube/services/$SERVICE_NAME/internal/handler"
	"kube/services/$SERVICE_NAME/internal/service"
	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"
)

// RegisterRoutes registers all $SERVICE_NAME service routes
func RegisterRoutes(r *server.Hertz, db *gorm.DB) {
	// Initialize $SERVICE_NAME service
	$(echo $SERVICE_NAME | tr '-' '_')Service := service.New${MODEL_NAME}Service(db)

	// Initialize $SERVICE_NAME handler
	$(echo $SERVICE_NAME | tr '-' '_')Handler := handler.New${MODEL_NAME}Handler($(echo $SERVICE_NAME | tr '-' '_')Service)

	// API v1 group
	v1 := r.Group("/api/v1")

	// $SERVICE_NAME routes
	$(echo $SERVICE_NAME | tr '-' '_')s := v1.Group("/$(echo $SERVICE_NAME | tr '-' '_')s")
	{
		$(echo $SERVICE_NAME | tr '-' '_')s.POST("/", $(echo $SERVICE_NAME | tr '-' '_')Handler.Create${MODEL_NAME})
		$(echo $SERVICE_NAME | tr '-' '_')s.GET("/:id", $(echo $SERVICE_NAME | tr '-' '_')Handler.Get${MODEL_NAME})
		$(echo $SERVICE_NAME | tr '-' '_')s.PUT("/:id", $(echo $SERVICE_NAME | tr '-' '_')Handler.Update${MODEL_NAME})
		$(echo $SERVICE_NAME | tr '-' '_')s.DELETE("/:id", $(echo $SERVICE_NAME | tr '-' '_')Handler.Delete${MODEL_NAME})
		$(echo $SERVICE_NAME | tr '-' '_')s.GET("/", $(echo $SERVICE_NAME | tr '-' '_')Handler.List${MODEL_NAME}s)
	}
}
EOF

# Create config
cat > "$CONFIG_DIR/config.go" << EOF
package config

import (
	"kube/shared/internal/config"
)

type ${MODEL_NAME}Config struct {
	*config.BaseConfig
	// Add service-specific configuration here
}

func Load${MODEL_NAME}Config() *${MODEL_NAME}Config {
	return &${MODEL_NAME}Config{
		BaseConfig: config.Load(),
		// Load service-specific config
	}
}
EOF

# Create migration
cat > "$MIGRATIONS_DIR/001_create_$(echo $SERVICE_NAME | tr '-' '_')s.sql" << EOF
-- Migration: Create $(echo $SERVICE_NAME | tr '-' '_')s table
-- Created: $(date)

CREATE TABLE IF NOT EXISTS $(echo $SERVICE_NAME | tr '-' '_')s (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes
CREATE INDEX idx_$(echo $SERVICE_NAME | tr '-' '_')s_status ON $(echo $SERVICE_NAME | tr '-' '_')s(status);
CREATE INDEX idx_$(echo $SERVICE_NAME | tr '-' '_')s_deleted_at ON $(echo $SERVICE_NAME | tr '-' '_')s(deleted_at);
EOF

# Create Dockerfile
cat > "deployments/docker/Dockerfile.$SERVICE_NAME" << EOF
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the $SERVICE_NAME service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $SERVICE_NAME ./services/$SERVICE_NAME/cmd

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/$SERVICE_NAME .

# Copy migrations
COPY --from=builder /app/services/$SERVICE_NAME/migrations ./migrations

# Expose port
EXPOSE $PORT

# Command to run
CMD ["./$SERVICE_NAME"]
EOF

echo "✅ Service $SERVICE_NAME created successfully!"
echo "📁 Location: services/$SERVICE_NAME/"
echo "🌐 Port: $PORT"
echo "🐳 Dockerfile: deployments/docker/Dockerfile.$SERVICE_NAME"
echo ""
echo "🚀 Next steps:"
echo "1. Update go.mod if needed"
echo "2. Add service-specific business logic"
echo "3. Update docker-compose.yml to include the new service"
echo "4. Run: go run services/$SERVICE_NAME/cmd/main.go"
echo ""
echo "📝 Don't forget to:"
echo "- Update environment variables"
echo "- Add service to monitoring"
echo "- Update API documentation"
