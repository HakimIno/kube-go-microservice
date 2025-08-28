package auth

import (
	"fmt"
	"time"

	"kube/internal/middleware"
	apperrors "kube/pkg/errors"
	"kube/pkg/models"
	"kube/pkg/services"
	"kube/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	*services.BaseService
	jwtSecret string
}

func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{
		BaseService: services.NewBaseService(db),
		jwtSecret:   jwtSecret,
	}
}

func (s *Service) Login(req *models.UserLoginRequest) (*models.UserResponse, string, error) {
	var user models.User
	if err := s.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, "", apperrors.New(apperrors.ErrCodeInvalidCredentials, "Invalid credentials", "Email or password is incorrect")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", apperrors.New(apperrors.ErrCodeInvalidCredentials, "Invalid credentials", "Email or password is incorrect")
	}

	if !user.IsActive {
		return nil, "", apperrors.New(apperrors.ErrCodeAccountDeactivated, "Account deactivated", "Your account has been deactivated")
	}

	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, "", apperrors.Wrap(err, apperrors.ErrCodeInternalError, "Token generation failed", err.Error())
	}

	return s.toUserResponse(&user), tokenString, nil
}

func (s *Service) RefreshToken(userID uint) (*models.UserResponse, string, error) {
	var user models.User
	if err := s.GetDB().First(&user, userID).Error; err != nil {
		return nil, "", apperrors.Wrap(err, apperrors.ErrCodeUserNotFound, "User not found", "User not found")
	}

	if !user.IsActive {
		return nil, "", apperrors.New(apperrors.ErrCodeAccountDeactivated, "Account deactivated", "Your account has been deactivated")
	}

	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, "", apperrors.Wrap(err, apperrors.ErrCodeInternalError, "Token generation failed", err.Error())
	}

	return s.toUserResponse(&user), tokenString, nil
}

func (s *Service) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User
	if err := s.GetDB().First(&user, userID).Error; err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeUserNotFound, "User not found", "User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return apperrors.New(apperrors.ErrCodeInvalidCredentials, "Invalid current password", "Current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeInternalError, "Password hashing failed", err.Error())
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.GetDB().Save(&user).Error; err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to update password", err.Error())
	}

	return nil
}

func (s *Service) GenerateQRCode(req *models.QRCodeRequest) (*models.QRCodeResponse, error) {
	sessionID := fmt.Sprintf("qr_%s", uuid.New().String()[:16])

	// สร้าง QR code ที่มีแค่เส้นเดียว - URL สำหรับ mobile app
	qrData := fmt.Sprintf("kube://qr-login?session_id=%s", sessionID)

	// สร้าง QR code พร้อม logo ตรงกลาง
	logoPath := "assets/logo_app.png"
	qrCodeImage, err := utils.GenerateQRCodeWithLogo(qrData, logoPath, 256)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeInternalError, "Failed to generate QR code", err.Error())
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	session := &models.QRLoginSession{
		ID:         sessionID,
		Status:     "pending",
		QRCodeData: qrData,
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.GetDB().Create(session).Error; err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to create QR login session", err.Error())
	}

	return &models.QRCodeResponse{
		SessionID:   sessionID,
		QRCodeImage: qrCodeImage,
		ExpiresAt:   expiresAt,
	}, nil
}

// QRScan - mobile app ส่ง app token เมื่อสแกน QR code
func (s *Service) QRScan(req *models.QRScanRequest) error {
	var session models.QRLoginSession
	if err := s.GetDB().Where("id = ? AND status = 'pending' AND expires_at > ?", req.SessionID, time.Now()).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperrors.New(apperrors.ErrCodeInvalidCredentials, "Invalid or expired session", "QR login session not found or expired")
		}
		return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Database error", err.Error())
	}

	// ตรวจสอบ app token
	user, err := s.validateAppToken(req.AppToken)
	if err != nil {
		return err
	}

	// อัปเดต session เป็น "scanned"
	session.UserID = &user.ID
	session.Status = "scanned"
	session.UpdatedAt = time.Now()

	if err := s.GetDB().Save(&session).Error; err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to update session", err.Error())
	}

	return nil
}

func (s *Service) GetQRLoginStatus(sessionID string) (*models.QRLoginStatusResponse, error) {
	var session models.QRLoginSession
	if err := s.GetDB().Preload("User").Where("id = ?", sessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.New(apperrors.ErrCodeInvalidCredentials, "Session not found", "QR login session not found")
		}
		return nil, apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Database error", err.Error())
	}

	response := &models.QRLoginStatusResponse{
		SessionID: session.ID,
		Status:    session.Status,
		Message:   s.getStatusMessage(session.Status),
	}

	if session.Status == "confirmed" && session.User != nil {
		response.User = s.toUserResponse(session.User)

		claims := &middleware.Claims{
			UserID: session.User.ID,
			Email:  session.User.Email,
			Role:   session.User.Role,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(s.jwtSecret))
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.ErrCodeInternalError, "Token generation failed", err.Error())
		}
		response.Token = tokenString
		response.Message = "Login successful"
	}

	if time.Now().After(session.ExpiresAt) && session.Status == "pending" {
		session.Status = "expired"
		session.UpdatedAt = time.Now()
		s.GetDB().Save(&session)
		response.Status = "expired"
		response.Message = "Session expired"
	}

	return response, nil
}

// validateAppToken ตรวจสอบ app token และดึงข้อมูล user
func (s *Service) validateAppToken(appToken string) (*models.User, error) {
	// Parse และ validate JWT token
	token, err := jwt.ParseWithClaims(appToken, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, apperrors.New(apperrors.ErrCodeInvalidCredentials, "Invalid app token", "App token is invalid or expired")
	}

	claims, ok := token.Claims.(*middleware.Claims)
	if !ok {
		return nil, apperrors.New(apperrors.ErrCodeInvalidCredentials, "Invalid token format", "Token claims format is invalid")
	}

	// ดึงข้อมูล user จาก database
	var user models.User
	if err := s.GetDB().Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.New(apperrors.ErrCodeInvalidCredentials, "User not found", "User associated with token not found")
		}
		return nil, apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Database error", err.Error())
	}

	if !user.IsActive {
		return nil, apperrors.New(apperrors.ErrCodeAccountDeactivated, "Account deactivated", "Your account has been deactivated")
	}

	return &user, nil
}

func (s *Service) CleanupExpiredSessions() error {
	result := s.GetDB().Where("expires_at < ? AND status = 'pending'", time.Now()).Delete(&models.QRLoginSession{})
	if result.Error != nil {
		return apperrors.Wrap(result.Error, apperrors.ErrCodeDatabaseError, "Failed to cleanup expired sessions", result.Error.Error())
	}
	return nil
}

func (s *Service) getStatusMessage(status string) string {
	switch status {
	case "pending":
		return "Waiting for QR code scan"
	case "scanned":
		return "QR code scanned, waiting for confirmation"
	case "confirmed":
		return "Login successful"
	case "expired":
		return "Session expired"
	default:
		return "Unknown status"
	}
}

func (s *Service) toUserResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
