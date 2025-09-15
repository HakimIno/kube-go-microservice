package models

import (
	"time"
)

type QRLoginSession struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	UserID     *uint     `json:"user_id" gorm:"index"`
	Status     string    `json:"status" gorm:"type:varchar(20);default:'pending'"` // pending, confirmed, rejected, expired
	QRCodeData string    `json:"qr_code_data" gorm:"type:text"`
	ExpiresAt  time.Time `json:"expires_at" gorm:"index"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type QRCodeRequest struct {
	DeviceInfo string `json:"device_info,omitempty" example:"Mobile App v1.0"`
}

type QRCodeResponse struct {
	SessionID   string    `json:"session_id" example:"qr_abc123def456"`
	QRCodeImage string    `json:"qr_code_image" example:"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAA..."` // Base64 encoded PNG
	ExpiresAt   time.Time `json:"expires_at" example:"2025-08-27T08:20:03Z"`
}

// QR Confirm Request - สำหรับ mobile app ส่งมาเมื่อสแกน QR code และ approve
type QRConfirmRequest struct {
	SessionID string `json:"session_id" binding:"required" example:"qr_abc123def456"`
	AppToken  string `json:"app_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// QR Reject Request - สำหรับ mobile app ส่งมาเมื่อสแกน QR code และ reject
type QRRejectRequest struct {
	SessionID string `json:"session_id" binding:"required" example:"qr_abc123def456"`
	AppToken  string `json:"app_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type QRLoginStatusRequest struct {
	SessionID string `json:"session_id" binding:"required" example:"qr_abc123def456"`
}

type QRLoginStatusResponse struct {
	SessionID string `json:"session_id" example:"qr_abc123def456"`
	Status    string `json:"status" example:"confirmed"`
	Token     string `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Message   string `json:"message" example:"Login successful"`
}

// Response Models for Swagger Documentation

type GenerateQRResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"QR code generated successfully"`
	Data    struct {
		QRCode QRCodeResponse `json:"qr_code"`
	} `json:"data"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/qr/generate"`
	Method    string `json:"method" example:"POST"`
}

type QRLoginStatusResponseWrapper struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Session status retrieved"`
	Data    struct {
		Status QRLoginStatusResponse `json:"status"`
	} `json:"data"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/qr/status"`
	Method    string `json:"method" example:"GET"`
}

type QRConfirmResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"QR scan and login confirmed"`
	Data    struct {
		SessionID string `json:"session_id" example:"qr_abc123def456"`
		Status    string `json:"status" example:"confirmed"`
	} `json:"data"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/qr/confirm"`
	Method    string `json:"method" example:"POST"`
}

type QRRejectResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"QR login rejected"`
	Data    struct {
		SessionID string `json:"session_id" example:"qr_abc123def456"`
		Status    string `json:"status" example:"rejected"`
	} `json:"data"`
	Timestamp string `json:"timestamp" example:"2025-08-27T08:15:03Z"`
	Path      string `json:"path" example:"/api/v1/users/qr/reject"`
	Method    string `json:"method" example:"POST"`
}
