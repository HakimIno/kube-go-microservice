package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"kube/internal/middleware"
	"kube/pkg/errors"
	"kube/pkg/models"
	"kube/pkg/services"
	"kube/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	*services.BaseService
	jwtSecret string
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		BaseService: services.NewBaseService(db),
		jwtSecret:   jwtSecret,
	}
}

func (s *AuthService) generateSecureSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	sessionID := base64.URLEncoding.EncodeToString(bytes)

	if len(sessionID) > 32 {
		sessionID = sessionID[:32]
	}

	return fmt.Sprintf("kube%s", sessionID), nil
}

func (s *AuthService) Login(req *models.UserLoginRequest) (*models.UserResponse, string, error) {
	var user models.User
	if err := s.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, "", errors.New(errors.ErrCodeInvalidCredentials, "Invalid credentials", "Email or password is incorrect")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", errors.New(errors.ErrCodeInvalidCredentials, "Invalid credentials", "Email or password is incorrect")
	}

	if !user.IsActive {
		return nil, "", errors.New(errors.ErrCodeAccountDeactivated, "Account deactivated", "Your account has been deactivated")
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
		return nil, "", errors.Wrap(err, errors.ErrCodeInternalError, "Token generation failed", err.Error())
	}

	return s.toUserResponse(&user), tokenString, nil
}

func (s *AuthService) RefreshToken(userID uint) (*models.UserResponse, string, error) {
	var user models.User
	if err := s.GetDB().First(&user, userID).Error; err != nil {
		return nil, "", errors.Wrap(err, errors.ErrCodeUserNotFound, "User not found", "User not found")
	}

	if !user.IsActive {
		return nil, "", errors.New(errors.ErrCodeAccountDeactivated, "Account deactivated", "Your account has been deactivated")
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
		return nil, "", errors.Wrap(err, errors.ErrCodeInternalError, "Token generation failed", err.Error())
	}

	return s.toUserResponse(&user), tokenString, nil
}

func (s *AuthService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User
	if err := s.GetDB().First(&user, userID).Error; err != nil {
		return errors.Wrap(err, errors.ErrCodeUserNotFound, "User not found", "User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New(errors.ErrCodeInvalidCredentials, "Invalid current password", "Current password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, errors.ErrCodeInternalError, "Password hashing failed", err.Error())
	}

	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.GetDB().Save(&user).Error; err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "Failed to update password", err.Error())
	}

	return nil
}

func (s *AuthService) Logout(userID uint, req *models.LogoutRequest) error {
	// สำหรับ JWT tokens ปกติแล้วไม่จำเป็นต้องทำอะไร server-side
	// เพราะ tokens จะหมดอายุเองตามเวลาที่กำหนด
	// แต่ถ้าต้องการ server-side logout สามารถ implement token blacklist ได้ที่นี่

	// ในอนาคตสามารถเพิ่ม logic นี้ได้:
	// - บันทึก token ลงใน blacklist
	// - ลบ refresh token จาก database
	// - อัปเดต user session status

	// ตอนนี้แค่ return success
	return nil
}

func (s *AuthService) GenerateQRCode(req *models.QRCodeRequest) (*models.QRCodeResponse, error) {
	sessionID, err := s.generateSecureSessionID()
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternalError, "Failed to generate session ID", err.Error())
	}

	qrData := fmt.Sprintf("kube://qr-login?session_id=%s", sessionID)

	logoPath := "assets/logo.jpg"
	qrCodeImage, err := utils.GenerateQRCodeWithLogo(qrData, logoPath, 256)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternalError, "Failed to generate QR code", err.Error())
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
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseError, "Failed to create QR login session", err.Error())
	}

	return &models.QRCodeResponse{
		SessionID:   sessionID,
		QRCodeImage: qrCodeImage,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *AuthService) QRConfirm(req *models.QRConfirmRequest) error {
	var session models.QRLoginSession
	if err := s.GetDB().Where("id = ? AND status = 'pending' AND expires_at > ?", req.SessionID, time.Now()).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ErrCodeInvalidCredentials, "Invalid or expired session", "QR login session not found or expired")
		}
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "Database error", err.Error())
	}

	// ตรวจสอบ app token
	user, err := s.validateAppToken(req.AppToken)
	if err != nil {
		return err
	}

	// อัปเดต session เป็น "confirmed" เมื่อ mobile app approve
	session.UserID = &user.ID
	session.Status = "confirmed"
	session.UpdatedAt = time.Now()

	if err := s.GetDB().Save(&session).Error; err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "Failed to update session", err.Error())
	}

	return nil
}

func (s *AuthService) QRReject(req *models.QRRejectRequest) error {
	var session models.QRLoginSession
	if err := s.GetDB().Where("id = ? AND status = 'pending' AND expires_at > ?", req.SessionID, time.Now()).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ErrCodeInvalidCredentials, "Invalid or expired session", "QR login session not found or expired")
		}
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "Database error", err.Error())
	}

	// ตรวจสอบ app token
	user, err := s.validateAppToken(req.AppToken)
	if err != nil {
		return err
	}

	// อัปเดต session เป็น "rejected" เมื่อ mobile app reject
	session.UserID = &user.ID
	session.Status = "rejected"
	session.UpdatedAt = time.Now()

	if err := s.GetDB().Save(&session).Error; err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "Failed to update session", err.Error())
	}

	return nil
}

func (s *AuthService) GetQRLoginStatus(sessionID string) (*models.QRLoginStatusResponse, error) {
	var session models.QRLoginSession
	if err := s.GetDB().Where("id = ?", sessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrCodeInvalidCredentials, "Session not found", "QR login session not found")
		}
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseError, "Database error", err.Error())
	}

	response := &models.QRLoginStatusResponse{
		SessionID: session.ID,
		Status:    session.Status,
		Message:   s.getStatusMessage(session.Status),
	}

	if session.Status == "confirmed" {
		// สร้าง JWT token สำหรับ user ที่ confirmed
		// ต้องโหลด user data เฉพาะตอนนี้
		var user models.User
		if err := s.GetDB().Where("id = ?", session.UserID).First(&user).Error; err != nil {
			return nil, errors.Wrap(err, errors.ErrCodeDatabaseError, "Failed to load user data", err.Error())
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
			return nil, errors.Wrap(err, errors.ErrCodeInternalError, "Token generation failed", err.Error())
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

func (s *AuthService) validateAppToken(appToken string) (*models.User, error) {
	// Parse และ validate JWT token
	token, err := jwt.ParseWithClaims(appToken, &middleware.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New(errors.ErrCodeInvalidCredentials, "Invalid app token", "App token is invalid or expired")
	}

	claims, ok := token.Claims.(*middleware.Claims)
	if !ok {
		return nil, errors.New(errors.ErrCodeInvalidCredentials, "Invalid token format", "Token claims format is invalid")
	}

	// ดึงข้อมูล user จาก database
	var user models.User
	if err := s.GetDB().Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrCodeInvalidCredentials, "User not found", "User associated with token not found")
		}
		return nil, errors.Wrap(err, errors.ErrCodeDatabaseError, "Database error", err.Error())
	}

	if !user.IsActive {
		return nil, errors.New(errors.ErrCodeAccountDeactivated, "Account deactivated", "Your account has been deactivated")
	}

	return &user, nil
}

func (s *AuthService) CleanupExpiredSessions() error {
	// ลบ expired sessions และ rejected sessions ที่หมดอายุแล้ว
	result := s.GetDB().Where("expires_at < ? AND (status = 'pending' OR status = 'rejected')", time.Now()).Delete(&models.QRLoginSession{})
	if result.Error != nil {
		return errors.Wrap(result.Error, errors.ErrCodeDatabaseError, "Failed to cleanup expired sessions", result.Error.Error())
	}
	return nil
}

func (s *AuthService) getStatusMessage(status string) string {
	switch status {
	case "pending":
		return "Waiting for QR code scan"
	case "confirmed":
		return "Login successful"
	case "rejected":
		return "Login rejected by user"
	case "expired":
		return "Session expired"
	default:
		return "Unknown status"
	}
}

func (s *AuthService) toUserResponse(user *models.User) *models.UserResponse {
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
