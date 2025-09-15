package service

import (
	"errors"
	"time"

	"kube/services/user/internal/model"
	"kube/services/user/internal/repository"
	apperrors "kube/shared/pkg/errors"
	"kube/shared/pkg/services"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	*services.BaseService
	repo      *repository.UserRepository
	jwtSecret string
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		BaseService: services.NewBaseService(db),
		repo:        repository.NewUserRepository(db),
		jwtSecret:   jwtSecret,
	}
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// Get user by username
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUnauthorizedError("Invalid credentials")
		}
		return nil, apperrors.NewInternalServerError("Failed to get user")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperrors.NewUnauthorizedError("Account is deactivated")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, apperrors.NewUnauthorizedError("Invalid credentials")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to generate access token")
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to generate refresh token")
	}

	return &model.LoginResponse{
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*model.LoginResponse, error) {
	// Parse and validate refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, apperrors.NewUnauthorizedError("Invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, apperrors.NewUnauthorizedError("Invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, apperrors.NewUnauthorizedError("Invalid user ID in token")
	}

	// Get user
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewUnauthorizedError("User not found")
		}
		return nil, apperrors.NewInternalServerError("Failed to get user")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperrors.NewUnauthorizedError("Account is deactivated")
	}

	// Generate new tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to generate access token")
	}

	newRefreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to generate refresh token")
	}

	return &model.LoginResponse{
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

func (s *AuthService) GenerateQRCode(req *model.QRCodeGenerateRequest) (*model.QRCodeGenerateResponse, error) {
	// Check if user exists
	user, err := s.repo.GetByID(req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("User not found")
		}
		return nil, apperrors.NewInternalServerError("Failed to get user")
	}

	// Generate token
	token := generateRandomToken()
	expiresAt := time.Now().Add(5 * time.Minute) // QR code expires in 5 minutes

	// Create QR login session
	session := &model.QRLoginSession{
		UserID:    user.ID,
		Token:     token,
		Status:    "pending",
		ExpiresAt: expiresAt,
	}

	if err := s.repo.CreateQRLoginSession(session); err != nil {
		return nil, apperrors.NewInternalServerError("Failed to create QR login session")
	}

	// Generate QR code URL (this would typically be a deep link or web URL)
	qrCodeURL := "https://yourapp.com/qr-login?token=" + token

	return &model.QRCodeGenerateResponse{
		QRCodeURL: qrCodeURL,
		Token:     token,
		ExpiresAt: expiresAt.Format(time.RFC3339),
	}, nil
}

func (s *AuthService) QRConfirm(req *model.QRCodeConfirmRequest) (*model.LoginResponse, error) {
	// Get QR login session
	session, err := s.repo.GetQRLoginSessionByToken(req.Token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("QR login session not found")
		}
		return nil, apperrors.NewInternalServerError("Failed to get QR login session")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return nil, apperrors.NewBadRequestError("QR login session expired")
	}

	// Check if session is already used
	if session.Status != "pending" {
		return nil, apperrors.NewBadRequestError("QR login session already used")
	}

	// Get user
	user, err := s.repo.GetByID(session.UserID)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to get user")
	}

	// Update session status
	session.Status = "confirmed"
	if err := s.repo.UpdateQRLoginSession(session); err != nil {
		return nil, apperrors.NewInternalServerError("Failed to update QR login session")
	}

	// Generate tokens
	accessToken, err := s.generateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to generate access token")
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to generate refresh token")
	}

	return &model.LoginResponse{
		User: model.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Avatar:    user.Avatar,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600, // 1 hour
	}, nil
}

func (s *AuthService) QRReject(token string) error {
	// Get QR login session
	session, err := s.repo.GetQRLoginSessionByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError("QR login session not found")
		}
		return apperrors.NewInternalServerError("Failed to get QR login session")
	}

	// Update session status
	session.Status = "rejected"
	if err := s.repo.UpdateQRLoginSession(session); err != nil {
		return apperrors.NewInternalServerError("Failed to update QR login session")
	}

	return nil
}

func (s *AuthService) GetQRLoginStatus(token string) (*model.QRCodeStatusResponse, error) {
	// Get QR login session
	session, err := s.repo.GetQRLoginSessionByToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("QR login session not found")
		}
		return nil, apperrors.NewInternalServerError("Failed to get QR login session")
	}

	return &model.QRCodeStatusResponse{
		Status:    session.Status,
		ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (s *AuthService) generateAccessToken(userID, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
		"type":     "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) generateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ChangePassword(userID string, currentPassword, newPassword string) error {
	// Get user
	user, err := s.repo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError("User not found")
		}
		return apperrors.NewInternalServerError("Failed to get user")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return apperrors.NewBadRequestError("Current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.NewInternalServerError("Failed to hash password")
	}

	user.Password = string(hashedPassword)
	if err := s.repo.Update(user); err != nil {
		return apperrors.NewInternalServerError("Failed to update password")
	}

	return nil
}

func generateRandomToken() string {
	// This is a simplified version - in production, use a proper random token generator
	return "qr_" + time.Now().Format("20060102150405") + "_" + "random"
}
