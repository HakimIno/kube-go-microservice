package templates

import (
	"kube/pkg/handlers"

	"github.com/cloudwego/hertz/pkg/app"
)

// HandlerTemplate shows how to create a handler using base handler
type HandlerTemplate struct {
	*handlers.BaseHandler
	// Add your specific dependencies here
	// service *ServiceTemplate
}

// NewHandlerTemplate creates a new handler template
func NewHandlerTemplate() *HandlerTemplate {
	return &HandlerTemplate{
		BaseHandler: handlers.NewBaseHandler(),
		// Initialize your specific dependencies
	}
}

// Example handler method
func (h *HandlerTemplate) ExampleHandler(c *app.RequestContext) {
	// Get parameter using base handler utility
	id, err := h.GetParamUint(c, "id")
	if err != nil {
		h.SendValidationError(c, "Invalid ID format")
		return
	}
	
	// Your handler logic here
	_ = id // Use id in your logic
	
	// Send response using base handler utilities
	h.SendSuccess(c, 200, nil, "Operation successful")
}
