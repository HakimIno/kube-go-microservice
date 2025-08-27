package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // "-" means don't include in JSON
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Role      string         `json:"role" gorm:"default:'user'"`
	Avatar    string         `json:"avatar"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserCreateRequest struct {
	Username  string `json:"username" binding:"required" example:"johndoe"`
	Email     string `json:"email" binding:"required,email" example:"john@example.com"`
	Password  string `json:"password" binding:"required,min=6" example:"password123"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Role      string `json:"role" example:"user"`
}

type UserUpdateRequest struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Avatar    string `json:"avatar" example:"https://example.com/avatar.jpg"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" example:"oldpassword123"`
	NewPassword     string `json:"new_password" binding:"required,min=6" example:"newpassword123"`
}

type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	FirstName string    `json:"first_name" example:"John"`
	LastName  string    `json:"last_name" example:"Doe"`
	Role      string    `json:"role" example:"user"`
	Avatar    string    `json:"avatar" example:"https://example.com/avatar.jpg"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2025-08-27T08:15:03.003428Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-08-27T08:15:03.003428Z"`
}

// Response Models for Swagger Documentation

type LoginResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Login successful"`
	Data    struct {
		User  UserResponse `json:"user"`
		Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	} `json:"data"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/login"`
	Method    string `json:"method" example:"POST"`
}

type RegisterResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"User created successfully"`
	Data    struct {
		User UserResponse `json:"user"`
	} `json:"data"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/register"`
	Method    string `json:"method" example:"POST"`
}

type RefreshTokenResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Token refreshed successfully"`
	Data    struct {
		User  UserResponse `json:"user"`
		Token string       `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	} `json:"data"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/refresh"`
	Method    string `json:"method" example:"POST"`
}

type GetUserResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"User retrieved successfully"`
	Data    struct {
		User UserResponse `json:"user"`
	} `json:"data"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/1"`
	Method    string `json:"method" example:"GET"`
}

type UpdateUserResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"User updated successfully"`
	Data    struct {
		User UserResponse `json:"user"`
	} `json:"data"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/1"`
	Method    string `json:"method" example:"PUT"`
}

type DeleteUserResponse struct {
	Success   bool    `json:"success" example:"true"`
	Message   string  `json:"message" example:"User deleted successfully"`
	Data      *string `json:"data" example:"null"`
	Timestamp string  `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string  `json:"path" example:"/api/v1/users/1"`
	Method    string  `json:"method" example:"DELETE"`
}

type ChangePasswordResponse struct {
	Success   bool    `json:"success" example:"true"`
	Message   string  `json:"message" example:"Password changed successfully"`
	Data      *string `json:"data" example:"null"`
	Timestamp string  `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string  `json:"path" example:"/api/v1/users/change-password"`
	Method    string  `json:"method" example:"POST"`
}

type ErrorResponse struct {
	Success bool `json:"success" example:"false"`
	Error   struct {
		Code      string `json:"code" example:"USER_NOT_FOUND"`
		Message   string `json:"message" example:"User not found"`
		Details   string `json:"details" example:"User with ID 999 not found"`
		Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03.832523143Z"`
		RequestID string `json:"request_id" example:"59744195-9e8d-4e07-b17e-39dac2ae2b48"`
	} `json:"error"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/999"`
	Method    string `json:"method" example:"GET"`
}

type ValidationErrorResponse struct {
	Success bool `json:"success" example:"false"`
	Error   struct {
		Code      string `json:"code" example:"VALIDATION_FAILED"`
		Message   string `json:"message" example:"Validation failed"`
		Details   string `json:"details" example:"Invalid request data format"`
		Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03.832523143Z"`
		RequestID string `json:"request_id" example:"59744195-9e8d-4e07-b17e-39dac2ae2b48"`
	} `json:"error"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/register"`
	Method    string `json:"method" example:"POST"`
}
