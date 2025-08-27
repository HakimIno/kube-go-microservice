package user

import (
	"strconv"

	"user-service/pkg/errors"
	"user-service/pkg/models"

	"github.com/cloudwego/hertz/pkg/app"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserCreateRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data or user already exists"
// @Router /api/v1/users/register [post]
func (h *Handler) Register(c *app.RequestContext) {
	var req models.UserCreateRequest
	if err := c.BindJSON(&req); err != nil {
		errors.SendValidationError(c, "Invalid request data format")
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	errors.SendSuccess(c, 201, user, "User created successfully")
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body models.UserLoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 401 {object} map[string]interface{} "Invalid credentials or account deactivated"
// @Failure 400 {object} map[string]interface{} "Invalid request data"
// @Router /api/v1/users/login [post]
func (h *Handler) Login(c *app.RequestContext) {
	var req models.UserLoginRequest
	if err := c.BindJSON(&req); err != nil {
		errors.SendValidationError(c, "Invalid request data format")
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
	errors.SendSuccess(c, 200, response, "Login successful")
}

// GetUser godoc
// @Summary Get user by ID
// @Description Retrieve user information by user ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User information"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/users/{id} [get]
func (h *Handler) GetUser(c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.SendValidationError(c, "Invalid user ID format")
		return
	}

	user, err := h.service.GetUserByID(uint(id))
	if err != nil {
		errors.SendError(c, err)
		return
	}

	errors.SendSuccess(c, 200, user, "User retrieved successfully")
}

// UpdateUser godoc
// @Summary Update user information
// @Description Update user profile information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UserUpdateRequest true "User update data"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request data or user ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/users/{id} [put]
func (h *Handler) UpdateUser(c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.SendValidationError(c, "Invalid user ID format")
		return
	}

	var req models.UserUpdateRequest
	if err := c.BindJSON(&req); err != nil {
		errors.SendValidationError(c, "Invalid request data format")
		return
	}

	user, err := h.service.UpdateUser(uint(id), &req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	errors.SendSuccess(c, 200, user, "User updated successfully")
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user account by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid user ID or deletion failed"
// @Router /api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(c *app.RequestContext) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.SendValidationError(c, "Invalid user ID format")
		return
	}

	if err := h.service.DeleteUser(uint(id)); err != nil {
		errors.SendError(c, err)
		return
	}

	errors.SendSuccess(c, 200, nil, "User deleted successfully")
}
