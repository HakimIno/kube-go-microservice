package handlers

import (
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// BaseHandler provides common handler utilities
type BaseHandler struct{}

// NewBaseHandler creates a new base handler
func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

// GetParamUint extracts and validates uint parameter from URL
func (h *BaseHandler) GetParamUint(c *app.RequestContext, paramName string) (uint, error) {
	paramStr := c.Param(paramName)
	id, err := strconv.ParseUint(paramStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// SendSuccess sends a successful response
func (h *BaseHandler) SendSuccess(c *app.RequestContext, statusCode int, data interface{}, message string) {
	response := utils.H{
		"success": true,
		"message": message,
	}
	
	if data != nil {
		response["data"] = data
	}
	
	c.JSON(statusCode, response)
}

// SendError sends an error response
func (h *BaseHandler) SendError(c *app.RequestContext, statusCode int, error string, message string) {
	c.JSON(statusCode, utils.H{
		"success": false,
		"error":   error,
		"message": message,
	})
}

// SendValidationError sends a validation error response
func (h *BaseHandler) SendValidationError(c *app.RequestContext, message string) {
	h.SendError(c, consts.StatusBadRequest, "validation_error", message)
}

// SendNotFound sends a not found response
func (h *BaseHandler) SendNotFound(c *app.RequestContext, message string) {
	h.SendError(c, consts.StatusNotFound, "not_found", message)
}

// SendInternalError sends an internal server error response
func (h *BaseHandler) SendInternalError(c *app.RequestContext, message string) {
	h.SendError(c, consts.StatusInternalServerError, "internal_error", message)
}

// SendUnauthorized sends an unauthorized response
func (h *BaseHandler) SendUnauthorized(c *app.RequestContext, message string) {
	h.SendError(c, consts.StatusUnauthorized, "unauthorized", message)
}
