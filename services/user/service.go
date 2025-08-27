package user

import (
	"time"

	apperrors "user-service/pkg/errors"
	"user-service/pkg/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	jwtSecret string
}

func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{db: db, jwtSecret: jwtSecret}
}

func (s *Service) CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error) {
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		return nil, apperrors.New(apperrors.ErrCodeUserAlreadyExists, "User already exists", "Email or username already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeInternalError, "Password hashing failed", err.Error())
	}

	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to create user", err.Error())
	}

	return s.toUserResponse(user), nil
}

func (s *Service) GetUserByID(id uint) (*models.UserResponse, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeUserNotFound, "User not found", "User with ID "+string(rune(id))+" not found")
	}
	return s.toUserResponse(&user), nil
}

func (s *Service) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeUserNotFound, "User not found", "User with ID "+string(rune(id))+" not found")
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Avatar = req.Avatar
	user.UpdatedAt = time.Now()

	if err := s.db.Save(&user).Error; err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to update user", err.Error())
	}

	return s.toUserResponse(&user), nil
}

func (s *Service) DeleteUser(id uint) error {
	if err := s.db.Delete(&models.User{}, id).Error; err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeDatabaseError, "Failed to delete user", err.Error())
	}
	return nil
}

func (s *Service) Login(req *models.UserLoginRequest) (*models.UserResponse, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, "", apperrors.New(apperrors.ErrCodeInvalidCredentials, "Invalid credentials", "Email or password is incorrect")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", apperrors.New(apperrors.ErrCodeInvalidCredentials, "Invalid credentials", "Email or password is incorrect")
	}

	if !user.IsActive {
		return nil, "", apperrors.New(apperrors.ErrCodeAccountDeactivated, "Account deactivated", "Your account has been deactivated")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, "", err
	}

	return s.toUserResponse(&user), tokenString, nil
}

func (s *Service) toUserResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.Avatar,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
