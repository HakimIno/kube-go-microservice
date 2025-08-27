#!/bin/bash

# Generate Service Script
# Usage: ./scripts/generate-service.sh <service-name> <port>

if [ $# -ne 2 ]; then
    echo "Usage: $0 <service-name> <port>"
    echo "Example: $0 video-service 8082"
    exit 1
fi

SERVICE_NAME=$1
PORT=$2
SERVICE_DIR="services/$SERVICE_NAME"
CMD_DIR="cmd/$SERVICE_NAME"

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

echo "Generating service: $SERVICE_NAME on port $PORT"
echo "Model name: $MODEL_NAME"

# Create directories
mkdir -p "$SERVICE_DIR"
mkdir -p "$CMD_DIR"

# Create service files
cat > "$SERVICE_DIR/service.go" << EOF
package ${SERVICE_NAME//-/_}

import (
	"time"

	apperrors "kube/pkg/errors"
	"kube/pkg/models"
	"kube/pkg/services"

	"gorm.io/gorm"
)

type Service struct {
	*services.BaseService
	// Add your specific dependencies here
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		BaseService: services.NewBaseService(db),
		// Initialize your specific dependencies
	}
}

// Example CRUD methods following user service pattern

func (s *Service) Create${MODEL_NAME}(req *models.${MODEL_NAME}CreateRequest) (*models.${MODEL_NAME}Response, error) {
	var item *models.${MODEL_NAME}

	err := s.WithTransaction(func(tx *gorm.DB) error {
		// Check if item already exists (customize based on your needs)
		var existingItem models.${MODEL_NAME}
		if err := tx.Where("name = ?", req.Name).First(&existingItem).Error; err == nil {
			return apperrors.New(apperrors.ErrCodeDuplicateRecord, "Item already exists", "Name already registered")
		}

		item = &models.${MODEL_NAME}{
			Name:        req.Name,
			Description: req.Description,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := tx.Create(item).Error; err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to create item", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.to${MODEL_NAME}Response(item), nil
}

func (s *Service) Get${MODEL_NAME}ByID(id uint) (*models.${MODEL_NAME}Response, error) {
	var item models.${MODEL_NAME}
	if err := s.GetDB().First(&item, id).Error; err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeRecordNotFound, "Item not found", "${MODEL_NAME} with ID "+string(rune(id))+" not found")
	}
	return s.to${MODEL_NAME}Response(&item), nil
}

func (s *Service) Update${MODEL_NAME}(id uint, req *models.${MODEL_NAME}UpdateRequest) (*models.${MODEL_NAME}Response, error) {
	var item *models.${MODEL_NAME}

	err := s.WithTransaction(func(tx *gorm.DB) error {
		if err := tx.First(&item, id).Error; err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeRecordNotFound, "Item not found", "${MODEL_NAME} with ID "+string(rune(id))+" not found")
		}

		item.Name = req.Name
		item.Description = req.Description
		item.UpdatedAt = time.Now()

		if err := tx.Save(item).Error; err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to update item", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.to${MODEL_NAME}Response(item), nil
}

func (s *Service) Delete${MODEL_NAME}(id uint) error {
	if err := s.GetDB().Delete(&models.${MODEL_NAME}{}, id).Error; err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to delete item", err.Error())
	}
	return nil
}

func (s *Service) to${MODEL_NAME}Response(item *models.${MODEL_NAME}) *models.${MODEL_NAME}Response {
	return &models.${MODEL_NAME}Response{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		IsActive:    item.IsActive,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}
EOF

# Create handler files
cat > "$SERVICE_DIR/handler.go" << EOF
package ${SERVICE_NAME//-/_}

import (
	"kube/pkg/errors"
	"kube/pkg/handlers"
	"kube/pkg/models"

	"github.com/cloudwego/hertz/pkg/app"
)

type Handler struct {
	*handlers.BaseHandler
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		BaseHandler: handlers.NewBaseHandler(),
		service:     service,
	}
}

// Create${MODEL_NAME} godoc
// @Summary Create a new ${MODEL_NAME}
// @Description Create a new ${MODEL_NAME} with the provided information
// @Tags ${MODEL_NAME}
// @Accept json
// @Produce json
// @Param ${MODEL_NAME} body models.${MODEL_NAME}CreateRequest true "${MODEL_NAME} creation data"
// @Success 201 {object} map[string]interface{} "${MODEL_NAME} created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data or ${MODEL_NAME} already exists"
// @Router /api/v1/${SERVICE_NAME//-/_} [post]
func (h *Handler) Create${MODEL_NAME}(c *app.RequestContext) {
	var req models.${MODEL_NAME}CreateRequest
	if err := c.BindJSON(&req); err != nil {
		h.SendValidationError(c, "Invalid request data format")
		return
	}

	item, err := h.service.Create${MODEL_NAME}(&req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 201, item, "${MODEL_NAME} created successfully")
}

// Get${MODEL_NAME} godoc
// @Summary Get ${MODEL_NAME} by ID
// @Description Retrieve ${MODEL_NAME} information by ID
// @Tags ${MODEL_NAME}
// @Accept json
// @Produce json
// @Param id path int true "${MODEL_NAME} ID"
// @Success 200 {object} map[string]interface{} "${MODEL_NAME} information"
// @Failure 400 {object} map[string]interface{} "Invalid ${MODEL_NAME} ID"
// @Failure 404 {object} map[string]interface{} "${MODEL_NAME} not found"
// @Router /api/v1/${SERVICE_NAME//-/_}/{id} [get]
func (h *Handler) Get${MODEL_NAME}(c *app.RequestContext) {
	id, err := h.GetParamUint(c, "id")
	if err != nil {
		h.SendValidationError(c, "Invalid ${MODEL_NAME} ID format")
		return
	}

	item, err := h.service.Get${MODEL_NAME}ByID(uint(id))
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, item, "${MODEL_NAME} retrieved successfully")
}

// Update${MODEL_NAME} godoc
// @Summary Update ${MODEL_NAME} information
// @Description Update ${MODEL_NAME} profile information
// @Tags ${MODEL_NAME}
// @Accept json
// @Produce json
// @Param id path int true "${MODEL_NAME} ID"
// @Param ${MODEL_NAME} body models.${MODEL_NAME}UpdateRequest true "${MODEL_NAME} update data"
// @Success 200 {object} map[string]interface{} "${MODEL_NAME} updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data or ${MODEL_NAME} ID"
// @Failure 404 {object} map[string]interface{} "${MODEL_NAME} not found"
// @Router /api/v1/${SERVICE_NAME//-/_}/{id} [put]
func (h *Handler) Update${MODEL_NAME}(c *app.RequestContext) {
	id, err := h.GetParamUint(c, "id")
	if err != nil {
		h.SendValidationError(c, "Invalid ${MODEL_NAME} ID format")
		return
	}

	var req models.${MODEL_NAME}UpdateRequest
	if err := c.BindJSON(&req); err != nil {
		h.SendValidationError(c, "Invalid request data format")
		return
	}

	item, err := h.service.Update${MODEL_NAME}(uint(id), &req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, item, "${MODEL_NAME} updated successfully")
}

// Delete${MODEL_NAME} godoc
// @Summary Delete ${MODEL_NAME}
// @Description Delete a ${MODEL_NAME} by ID
// @Tags ${MODEL_NAME}
// @Accept json
// @Produce json
// @Param id path int true "${MODEL_NAME} ID"
// @Success 200 {object} map[string]interface{} "${MODEL_NAME} deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ${MODEL_NAME} ID or deletion failed"
// @Router /api/v1/${SERVICE_NAME//-/_}/{id} [delete]
func (h *Handler) Delete${MODEL_NAME}(c *app.RequestContext) {
	id, err := h.GetParamUint(c, "id")
	if err != nil {
		h.SendValidationError(c, "Invalid ${MODEL_NAME} ID format")
		return
	}

	if err := h.service.Delete${MODEL_NAME}(uint(id)); err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, nil, "${MODEL_NAME} deleted successfully")
}
EOF

# Create routes files
cat > "$SERVICE_DIR/routes.go" << EOF
package ${SERVICE_NAME//-/_}

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterRoutes(h *server.Hertz, service *Service) {
	handler := NewHandler(service)

	// API routes with versioning
	api := h.Group("/api/v1/${SERVICE_NAME//-/_}")
	{
		api.POST("", func(ctx context.Context, c *app.RequestContext) {
			handler.Create${MODEL_NAME}(c)
		})
		api.GET("/:id", func(ctx context.Context, c *app.RequestContext) {
			handler.Get${MODEL_NAME}(c)
		})
		api.PUT("/:id", func(ctx context.Context, c *app.RequestContext) {
			handler.Update${MODEL_NAME}(c)
		})
		api.DELETE("/:id", func(ctx context.Context, c *app.RequestContext) {
			handler.Delete${MODEL_NAME}(c)
		})
	}
}
EOF

# Create main.go
cat > "$CMD_DIR/main.go" << EOF
package main

import (
	"log"

	_ "kube/docs" // This is generated by swag init
	"kube/internal/config"
	"kube/internal/database"
	"kube/pkg/models"
	"kube/pkg/server"
	"kube/services/${SERVICE_NAME}"
	"time"
)

// @title $SERVICE_NAME API
// @version 1.0
// @description This is a $SERVICE_NAME service API built with Hertz framework.

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

	// Auto migrate models (customize based on your needs)
	if err := db.AutoMigrate(&models.${MODEL_NAME}{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	${SERVICE_NAME//-/_}Service := ${SERVICE_NAME//-/_}.NewService(db)

	serverConfig := server.ServerConfig{
		Port:         "$PORT",
		ServiceName:  "$SERVICE_NAME",
		SwaggerURL:   "http://localhost:$PORT",
		RateLimit:    100,
		RateDuration: time.Minute,
	}

	srv := server.NewServer(serverConfig)
	${SERVICE_NAME//-/_}.RegisterRoutes(srv.Hertz, ${SERVICE_NAME//-/_}Service)
	srv.Start()
}
EOF

# Create run script
cat > "scripts/run-$SERVICE_NAME.sh" << EOF
#!/bin/bash

# Run $SERVICE_NAME
cd "\$(dirname "\$0")/.."

echo "Starting $SERVICE_NAME..."

# Set environment variables (customize as needed)
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
go run cmd/$SERVICE_NAME/main.go
EOF

chmod +x "scripts/run-$SERVICE_NAME.sh"

# Create model files
MODEL_FILE="pkg/models/${SERVICE_NAME//-/_}.go"
cat > "$MODEL_FILE" << EOF
package models

import (
	"time"
)

// ${MODEL_NAME} represents the ${MODEL_NAME} entity
type ${MODEL_NAME} struct {
	ID          uint      \`json:"id" gorm:"primaryKey"\`
	Name        string    \`json:"name" gorm:"not null;unique"\`
	Description string    \`json:"description"\`
	IsActive    bool      \`json:"is_active" gorm:"default:true"\`
	CreatedAt   time.Time \`json:"created_at"\`
	UpdatedAt   time.Time \`json:"updated_at"\`
}

// ${MODEL_NAME}CreateRequest represents the request payload for creating a ${MODEL_NAME}
type ${MODEL_NAME}CreateRequest struct {
	Name        string \`json:"name" binding:"required,min=1,max=255"\`
	Description string \`json:"description" binding:"max=1000"\`
}

// ${MODEL_NAME}UpdateRequest represents the request payload for updating a ${MODEL_NAME}
type ${MODEL_NAME}UpdateRequest struct {
	Name        string \`json:"name" binding:"omitempty,min=1,max=255"\`
	Description string \`json:"description" binding:"omitempty,max=1000"\`
}

// ${MODEL_NAME}Response represents the response payload for ${MODEL_NAME} data
type ${MODEL_NAME}Response struct {
	ID          uint      \`json:"id"\`
	Name        string    \`json:"name"\`
	Description string    \`json:"description"\`
	IsActive    bool      \`json:"is_active"\`
	CreatedAt   time.Time \`json:"created_at"\`
	UpdatedAt   time.Time \`json:"updated_at"\`
}
EOF

echo "Service $SERVICE_NAME generated successfully!"
echo "Files created:"
echo "  - $MODEL_FILE"
echo "  - $SERVICE_DIR/service.go"
echo "  - $SERVICE_DIR/handler.go"
echo "  - $SERVICE_DIR/routes.go"
echo "  - $CMD_DIR/main.go"
echo "  - scripts/run-$SERVICE_NAME.sh"
echo ""
echo "To build the service:"
echo "  ./scripts/build-service.sh $SERVICE_NAME"
echo ""
echo "To run the service:"
echo "  ./scripts/run-$SERVICE_NAME.sh"
echo ""
echo "To test the API:"
echo "  POST http://localhost:$PORT/api/v1/${SERVICE_NAME//-/_}"
echo "  GET  http://localhost:$PORT/api/v1/${SERVICE_NAME//-/_}/{id}"
