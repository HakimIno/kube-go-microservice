package templates

import (
	"kube/pkg/services"

	"gorm.io/gorm"
)

// ServiceTemplate shows how to create a service using base service
type ServiceTemplate struct {
	*services.BaseService
	// Add your specific dependencies here
	// jwtSecret string
	// otherConfig *Config
}

// NewServiceTemplate creates a new service template
func NewServiceTemplate(db *gorm.DB) *ServiceTemplate {
	return &ServiceTemplate{
		BaseService: services.NewBaseService(db),
		// Initialize your specific dependencies
	}
}

// Example method showing how to use base service
func (s *ServiceTemplate) ExampleMethod() error {
	// Use the database from base service
	db := s.GetDB()
	_ = db // Use db in your logic
	
	// Use transaction helper
	return s.WithTransaction(func(tx *gorm.DB) error {
		// Your transaction logic here
		return nil
	})
}
