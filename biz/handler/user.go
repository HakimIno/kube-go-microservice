package handler

import (
	"context"
	"kube/biz/service"
	"kube/pkg/errors"
	"kube/pkg/handlers"
	"kube/pkg/models"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserHandler struct {
	*handlers.BaseHandler
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: handlers.NewBaseHandler(),
		service:     service,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with the provided information
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserCreateRequest true "User registration data"
// @Success 201 {object} models.RegisterResponse "User created successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data"
// @Failure 409 {object} models.ErrorResponse "User already exists"
// @Router /api/v1/auth/register [post]
func (h *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
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

// GetCurrentUser godoc
// @Summary Get current user information
// @Description Retrieve information of the currently authenticated user. Requires valid Bearer token in Authorization header.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.GetUserResponse "Current user information"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Bearer token required or invalid"
// @Router /api/v1/users/me [get]
func (h *UserHandler) GetCurrentUser(ctx context.Context, c *app.RequestContext) {
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

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve user information by user ID. Requires valid Bearer token in Authorization header.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID" minimum(1)
// @Success 200 {object} models.GetUserResponse "User information"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid user ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Bearer token required or invalid"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(ctx context.Context, c *app.RequestContext) {
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
// @Description Update user profile information. Requires valid Bearer token in Authorization header.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID" minimum(1)
// @Param user body models.UserUpdateRequest true "User update data"
// @Success 200 {object} models.UpdateUserResponse "User updated successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data or user ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Bearer token required or invalid"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
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
// @Description Delete a user account by ID. Requires valid Bearer token in Authorization header.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID" minimum(1)
// @Success 200 {object} models.DeleteUserResponse "User deleted successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid user ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Bearer token required or invalid"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
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
