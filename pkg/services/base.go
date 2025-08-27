package services

import (
	"time"

	"gorm.io/gorm"
)

// BaseService provides common service utilities
type BaseService struct {
	db *gorm.DB
}

// NewBaseService creates a new base service
func NewBaseService(db *gorm.DB) *BaseService {
	return &BaseService{db: db}
}

// GetDB returns the database instance
func (s *BaseService) GetDB() *gorm.DB {
	return s.db
}

// BeginTransaction starts a new database transaction
func (s *BaseService) BeginTransaction() *gorm.DB {
	return s.db.Begin()
}

// CommitTransaction commits a database transaction
func (s *BaseService) CommitTransaction(tx *gorm.DB) error {
	return tx.Commit().Error
}

// RollbackTransaction rolls back a database transaction
func (s *BaseService) RollbackTransaction(tx *gorm.DB) error {
	return tx.Rollback().Error
}

// WithTransaction executes a function within a database transaction
func (s *BaseService) WithTransaction(fn func(*gorm.DB) error) error {
	tx := s.BeginTransaction()
	if err := fn(tx); err != nil {
		s.RollbackTransaction(tx)
		return err
	}
	return s.CommitTransaction(tx)
}

// SetTimestamps sets created_at and updated_at timestamps
func SetTimestamps(model interface{}) {
	// Use reflection to set timestamps if the model has these fields
	// This is a simplified version - in practice you might want to use reflection
	// or implement this in each specific model
	_ = time.Now() // Placeholder for future implementation
}

// SetUpdatedAt sets only the updated_at timestamp
func SetUpdatedAt(model interface{}) {
	// Implementation would depend on the model structure
	_ = time.Now() // Placeholder for future implementation
}
