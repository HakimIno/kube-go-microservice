package user

import (
	"kube/pkg/errors"
	"kube/pkg/handlers"
	"kube/pkg/models"

	"github.com/cloudwego/hertz/pkg/app"
)

type Handler struct {
	*handlers.BaseHandler
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		BaseHandler: handlers.NewBaseHandler(),
		service:     service,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserCreateRequest true "User registration data"
// @Success 201 {object} models.RegisterResponse "User created successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data"
// @Failure 409 {object} models.ErrorResponse "User already exists"
// @Router /api/v1/users/register [post]
func (h *Handler) Register(c *app.RequestContext) {
	var req models.UserCreateRequest
	if err := c.BindJSON(&req); err != nil {
		h.SendValidationError(c, "Invalid request data format")
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 201, user, "User created successfully")
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body models.UserLoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data"
// @Failure 401 {object} models.ErrorResponse "Invalid credentials"
// @Failure 403 {object} models.ErrorResponse "Account deactivated"
// @Router /api/v1/users/login [post]
func (h *Handler) Login(c *app.RequestContext) {
	var req models.UserLoginRequest
	if err := c.BindJSON(&req); err != nil {
		h.SendValidationError(c, "Invalid request data format")
		return
	}

	user, token, err := h.service.Login(&req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}
	h.SendSuccess(c, 200, response, "Login successful")
}

// RefreshToken godoc
// @Summary Refresh authentication token
// @Description Generate a new JWT token for the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.RefreshTokenResponse "Token refreshed successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /api/v1/users/refresh [post]
func (h *Handler) RefreshToken(c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.SendError(c, errors.New(errors.ErrCodeUnauthorized, "User not authenticated", "User ID not found in context"))
		return
	}

	user, token, err := h.service.RefreshToken(userID.(uint))
	if err != nil {
		errors.SendError(c, err)
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": token,
	}
	h.SendSuccess(c, 200, response, "Token refreshed successfully")
}

// GetCurrentUser godoc
// @Summary Get current user information
// @Description Retrieve information of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.GetUserResponse "Current user information"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Router /api/v1/users/me [get]
func (h *Handler) GetCurrentUser(c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.SendError(c, errors.New(errors.ErrCodeUnauthorized, "User not authenticated", "User ID not found in context"))
		return
	}

	user, err := h.service.GetCurrentUser(userID.(uint))
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, user, "Current user retrieved successfully")
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the password of the currently authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body models.ChangePasswordRequest true "Password change data"
// @Success 200 {object} models.ChangePasswordResponse "Password changed successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data"
// @Failure 401 {object} models.ErrorResponse "Unauthorized or invalid current password"
// @Router /api/v1/users/change-password [post]
func (h *Handler) ChangePassword(c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.SendError(c, errors.New(errors.ErrCodeUnauthorized, "User not authenticated", "User ID not found in context"))
		return
	}

	var req models.ChangePasswordRequest
	if err := c.BindJSON(&req); err != nil {
		h.SendValidationError(c, "Invalid request data format")
		return
	}

	if err := h.service.ChangePassword(userID.(uint), &req); err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, nil, "Password changed successfully")
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Success 200 {object} models.GetUserResponse "User information"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Router /api/v1/users/{id} [get]
func (h *Handler) GetUser(c *app.RequestContext) {
	id, err := h.GetParamUint(c, "id")
	if err != nil {
		h.SendValidationError(c, "Invalid user ID format")
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, user, "User retrieved successfully")
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user profile information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Param user body models.UserUpdateRequest true "User update data"
// @Success 200 {object} models.UpdateUserResponse "User updated successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data or user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Router /api/v1/users/{id} [put]
func (h *Handler) UpdateUser(c *app.RequestContext) {
	id, err := h.GetParamUint(c, "id")
	if err != nil {
		h.SendValidationError(c, "Invalid user ID format")
		return
	}

	var req models.UserUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		h.SendValidationError(c, "Invalid request data format")
		return
	}

	user, err := h.service.UpdateUser(uint(id), &req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, user, "User updated successfully")
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user account by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID" minimum(1)
// @Success 200 {object} models.DeleteUserResponse "User deleted successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid user ID"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Router /api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(c *app.RequestContext) {
	id, err := h.GetParamUint(c, "id")
	if err != nil {
		h.SendValidationError(c, "Invalid user ID format")
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, nil, "User deleted successfully")
}
