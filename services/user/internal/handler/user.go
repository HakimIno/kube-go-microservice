package handler

import (
	"context"
	"strconv"

	"kube/services/user/internal/model"
	"kube/services/user/internal/service"
	apperrors "kube/shared/pkg/errors"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(ctx context.Context, c *app.RequestContext) {
	var req model.UserCreateRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	user, err := h.service.CreateUser(&req)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(201, map[string]interface{}{
		"success": true,
		"data":    user,
		"message": "User created successfully",
	})
}

func (h *UserHandler) GetCurrentUser(ctx context.Context, c *app.RequestContext) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperrors.HandleError(c, apperrors.NewUnauthorizedError("User ID not found"))
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

func (h *UserHandler) GetUser(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("id")
	if userID == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("User ID is required"))
		return
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    user,
	})
}

func (h *UserHandler) UpdateUser(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("id")
	if userID == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("User ID is required"))
		return
	}

	var req model.UserUpdateRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	user, err := h.service.UpdateUser(userID, &req)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    user,
		"message": "User updated successfully",
	})
}

func (h *UserHandler) DeleteUser(ctx context.Context, c *app.RequestContext) {
	userID := c.Param("id")
	if userID == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("User ID is required"))
		return
	}

	err := h.service.DeleteUser(userID)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) ListUsers(ctx context.Context, c *app.RequestContext) {
	// Get pagination parameters
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 10 // default limit
	offset := 0 // default offset

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	users, err := h.service.ListUsers(limit, offset)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    users,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
		},
	})
}
