package user

import (
	"time"

	"kube/internal/middleware"
	apperrors "kube/pkg/errors"
	"kube/pkg/models"
	"kube/pkg/services"

	"github.com/golang-jwt/jwt/v5"
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

func (s *Service) CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error) {
	var user *models.User

	err := s.WithTransaction(func(tx *gorm.DB) error {
		var existingUser models.User
		if err := tx.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
			return apperrors.New(apperrors.ErrCodeUserAlreadyExists, "User already exists", "Email or username already registered")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeInternalError, "Password hashing failed", err.Error())
		}

		// Set default role if not provided
		role := req.Role
		if role == "" {
			role = "user" // Default role
		}

		user = &models.User{
			Username:  req.Username,
			Email:     req.Email,
			Password:  string(hashedPassword),
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Role:      role,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := tx.Create(user).Error; err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to create user", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *Service) GetUserByID(id uint) (*models.UserResponse, error) {
	var user models.User
	if err := s.GetDB().First(&user, id).Error; err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeUserNotFound, "User not found", "User with ID "+string(rune(id))+" not found")
	}
	return s.toUserResponse(&user), nil
}

func (s *Service) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	var user *models.User

	err := s.WithTransaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeUserNotFound, "User not found", "User with ID "+string(rune(id))+" not found")
		}

		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.Avatar = req.Avatar
		user.UpdatedAt = time.Now()

		if err := tx.Save(user).Error; err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to update user", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *Service) DeleteUser(id uint) error {
	if err := s.GetDB().Delete(&models.User{}, id).Error; err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to delete user", err.Error())
	}
	return nil
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

	// Generate JWT token with role
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24 hours
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

// RefreshToken generates a new token for the user
func (s *Service) RefreshToken(userID uint) (*models.UserResponse, string, error) {
	var user models.User
	if err := s.GetDB().First(&user, userID).Error; err != nil {
		return nil, "", apperrors.Wrap(err, apperrors.ErrCodeUserNotFound, "User not found", "User not found")
	}

	if !user.IsActive {
		return nil, "", apperrors.New(apperrors.ErrCodeAccountDeactivated, "Account deactivated", "Your account has been deactivated")
	}

	// Generate new JWT token
	claims := &middleware.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 24 hours
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

// ChangePassword allows users to change their password
func (s *Service) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	var user models.User
	if err := s.GetDB().First(&user, userID).Error; err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeUserNotFound, "User not found", "User not found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return apperrors.New(apperrors.ErrCodeInvalidCredentials, "Invalid current password", "Current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeInternalError, "Password hashing failed", err.Error())
	}

	// Update password
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	if err := s.GetDB().Save(&user).Error; err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to update password", err.Error())
	}

	return nil
}

// GetCurrentUser returns the current user's information
func (s *Service) GetCurrentUser(userID uint) (*models.UserResponse, error) {
	return s.GetUserByID(userID)
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
