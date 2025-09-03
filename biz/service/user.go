package service

import (
	"time"

	"kube/pkg/errors"
	"kube/pkg/models"
	"kube/pkg/services"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	*services.BaseService
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		BaseService: services.NewBaseService(db),
	}
}

func (s *UserService) CreateUser(req *models.UserCreateRequest) (*models.UserResponse, error) {
	var user *models.User

	err := s.WithTransaction(func(tx *gorm.DB) error {
		var existingUser models.User
		if err := tx.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
			return errors.New(errors.ErrCodeUserAlreadyExists, "User already exists", "Email or username already registered")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Wrap(err, errors.ErrCodeInternalError, "Password hashing failed", err.Error())
		}

		role := req.Role
		if role == "" {
			role = "user"
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
			return errors.Wrap(err, errors.ErrCodeDatabaseError, "Failed to create user", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) GetUserByID(id uint) (*models.UserResponse, error) {
	var user models.User
	if err := s.GetDB().First(&user, id).Error; err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeUserNotFound, "User not found", "User with ID "+string(rune(id))+" not found")
	}
	return s.toUserResponse(&user), nil
}

func (s *UserService) UpdateUser(id uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	var user *models.User

	err := s.WithTransaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, id).Error; err != nil {
			return errors.Wrap(err, errors.ErrCodeUserNotFound, "User not found", "User with ID "+string(rune(id))+" not found")
		}

		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.Avatar = req.Avatar
		user.UpdatedAt = time.Now()

		if err := tx.Save(user).Error; err != nil {
			return errors.Wrap(err, errors.ErrCodeDatabaseError, "Failed to update user", err.Error())
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) DeleteUser(id uint) error {
	if err := s.GetDB().Delete(&models.User{}, id).Error; err != nil {
		return errors.Wrap(err, errors.ErrCodeDatabaseError, "Failed to delete user", err.Error())
	}
	return nil
}

func (s *UserService) GetCurrentUser(userID uint) (*models.UserResponse, error) {
	return s.GetUserByID(userID)
}

func (s *UserService) toUserResponse(user *models.User) *models.UserResponse {
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
