package models

import (
	"time"
)

// Channel represents a user's channel
type Channel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	BannerImage string    `json:"banner_image"`
	Avatar      string    `json:"avatar"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ChannelCreateRequest represents the request to create a channel
type ChannelCreateRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	BannerImage string `json:"banner_image"`
	Avatar      string `json:"avatar"`
}

// ChannelUpdateRequest represents the request to update a channel
type ChannelUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BannerImage string `json:"banner_image"`
	Avatar      string `json:"avatar"`
}

// ChannelResponse represents the response for channel data
type ChannelResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	BannerImage string    `json:"banner_image"`
	Avatar      string    `json:"avatar"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
} 