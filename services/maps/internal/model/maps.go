package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system (for maps service compatibility)
type User struct {
	ID        string         `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // Hidden from JSON
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Phone     string         `json:"phone"`
	Avatar    string         `json:"avatar"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// QRLoginSession represents a QR code login session (for maps service compatibility)
type QRLoginSession struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	Status    string    `json:"status" gorm:"default:'pending'"` // pending, confirmed, rejected, expired
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Place represents a place from HERE API
type Place struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	Category     string   `json:"category"`
	Distance     int      `json:"distance"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Phone        string   `json:"phone"`
	Website      string   `json:"website"`
	OpeningHours []string `json:"opening_hours"`
	Rating       float64  `json:"rating"`
	PriceLevel   int      `json:"price_level"`
}

// SearchPlacesRequest represents the request to search for places
type SearchPlacesRequest struct {
	Query     string  `json:"query" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    int     `json:"radius"` // in meters
	Limit     int     `json:"limit"`
	Category  string  `json:"category"`
}

// SearchPlacesResponse represents the response for place search
type SearchPlacesResponse struct {
	Places []Place `json:"places"`
	Total  int     `json:"total"`
	Query  string  `json:"query"`
}

// PlaceDetails represents detailed information about a place
type PlaceDetails struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Address      string   `json:"address"`
	Category     string   `json:"category"`
	Latitude     float64  `json:"latitude"`
	Longitude    float64  `json:"longitude"`
	Phone        string   `json:"phone"`
	Website      string   `json:"website"`
	OpeningHours []string `json:"opening_hours"`
	Rating       float64  `json:"rating"`
	PriceLevel   int      `json:"price_level"`
	Description  string   `json:"description"`
	Photos       []string `json:"photos"`
	Reviews      []Review `json:"reviews"`
}

// Review represents a review for a place
type Review struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

// SearchHistory represents a search history entry
type SearchHistory struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Query     string    `json:"query" gorm:"not null"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Radius    int       `json:"radius"`
	Results   int       `json:"results"`
	CreatedAt time.Time `json:"created_at"`
}

// FavoritePlace represents a user's favorite place
type FavoritePlace struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null"`
	PlaceID   string    `json:"place_id" gorm:"not null"`
	Title     string    `json:"title" gorm:"not null"`
	Address   string    `json:"address" gorm:"not null"`
	Latitude  float64   `json:"latitude" gorm:"not null"`
	Longitude float64   `json:"longitude" gorm:"not null"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// AddFavoriteRequest represents the request to add a favorite place
type AddFavoriteRequest struct {
	PlaceID   string  `json:"place_id" binding:"required"`
	Title     string  `json:"title" binding:"required"`
	Address   string  `json:"address" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Category  string  `json:"category"`
}

// GetFavoritesResponse represents the response for getting favorites
type GetFavoritesResponse struct {
	Favorites []FavoritePlace `json:"favorites"`
	Total     int             `json:"total"`
}
