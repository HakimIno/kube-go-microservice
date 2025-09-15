package service

import (
	"errors"

	"kube/services/user/internal/model"
	"kube/services/user/internal/repository"
	apperrors "kube/shared/pkg/errors"
	"kube/shared/pkg/services"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	*services.BaseService
	repo *repository.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		BaseService: services.NewBaseService(db),
		repo:        repository.NewUserRepository(db),
	}
}

func (s *UserService) CreateUser(req *model.UserCreateRequest) (*model.UserResponse, error) {
	// Check if username already exists
	existingUser, _ := s.repo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, apperrors.NewBadRequestError("Username already exists")
	}

	// Check if email already exists
	existingEmail, _ := s.repo.GetByEmail(req.Email)
	if existingEmail != nil {
		return nil, apperrors.NewBadRequestError("Email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to hash password")
	}

	// Create user
	user := &model.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		IsActive:  true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, apperrors.NewInternalServerError("Failed to create user")
	}

	return s.userToResponse(user), nil
}

func (s *UserService) GetUser(id string) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("User not found")
		}
		return nil, apperrors.NewInternalServerError("Failed to get user")
	}

	return s.userToResponse(user), nil
}

func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("User not found")
		}
		return nil, apperrors.NewInternalServerError("Failed to get user")
	}

	return user, nil
}

func (s *UserService) UpdateUser(id string, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NewNotFoundError("User not found")
		}
		return nil, apperrors.NewInternalServerError("Failed to get user")
	}

	// Update fields if provided
	if req.Username != "" {
		// Check if new username already exists
		existingUser, _ := s.repo.GetByUsername(req.Username)
		if existingUser != nil && existingUser.ID != id {
			return nil, apperrors.NewBadRequestError("Username already exists")
		}
		user.Username = req.Username
	}

	if req.Email != "" {
		// Check if new email already exists
		existingEmail, _ := s.repo.GetByEmail(req.Email)
		if existingEmail != nil && existingEmail.ID != id {
			return nil, apperrors.NewBadRequestError("Email already exists")
		}
		user.Email = req.Email
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}

	if err := s.repo.Update(user); err != nil {
		return nil, apperrors.NewInternalServerError("Failed to update user")
	}

	return s.userToResponse(user), nil
}

func (s *UserService) DeleteUser(id string) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError("User not found")
		}
		return apperrors.NewInternalServerError("Failed to get user")
	}

	if err := s.repo.Delete(id); err != nil {
		return apperrors.NewInternalServerError("Failed to delete user")
	}

	return nil
}

func (s *UserService) ListUsers(limit, offset int) ([]*model.UserResponse, error) {
	users, err := s.repo.List(limit, offset)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to list users")
	}

	var responses []*model.UserResponse
	for _, user := range users {
		responses = append(responses, s.userToResponse(user))
	}

	return responses, nil
}

func (s *UserService) ChangePassword(id string, req *model.ChangePasswordRequest) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.NewNotFoundError("User not found")
		}
		return apperrors.NewInternalServerError("Failed to get user")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return apperrors.NewBadRequestError("Current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperrors.NewInternalServerError("Failed to hash password")
	}

	user.Password = string(hashedPassword)
	if err := s.repo.Update(user); err != nil {
		return apperrors.NewInternalServerError("Failed to update password")
	}

	return nil
}

func (s *UserService) userToResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
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
	}
}
