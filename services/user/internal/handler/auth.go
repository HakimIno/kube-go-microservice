package handler

import (
	"context"
	"kube/services/user/internal/model"
	"kube/services/user/internal/service"
	apperrors "kube/shared/pkg/errors"

	"github.com/cloudwego/hertz/pkg/app"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
	var req model.LoginRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	response, err := h.service.Login(&req)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    response,
		"message": "Login successful",
	})
}

func (h *AuthHandler) RefreshToken(ctx context.Context, c *app.RequestContext) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	response, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    response,
		"message": "Token refreshed successfully",
	})
}

func (h *AuthHandler) ChangePassword(ctx context.Context, c *app.RequestContext) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperrors.HandleError(c, apperrors.NewUnauthorizedError("User ID not found"))
		return
	}

	var req model.ChangePasswordRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	err := h.service.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "Password changed successfully",
	})
}

func (h *AuthHandler) Logout(ctx context.Context, c *app.RequestContext) {
	// In a stateless JWT system, logout is typically handled on the client side
	// by removing the token. However, you could implement a token blacklist here
	// if needed for additional security.

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "Logout successful",
	})
}

func (h *AuthHandler) GenerateQRCode(ctx context.Context, c *app.RequestContext) {
	var req model.QRCodeGenerateRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	response, err := h.service.GenerateQRCode(&req)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    response,
		"message": "QR code generated successfully",
	})
}

func (h *AuthHandler) QRConfirm(ctx context.Context, c *app.RequestContext) {
	var req model.QRCodeConfirmRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	response, err := h.service.QRConfirm(&req)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    response,
		"message": "QR login confirmed successfully",
	})
}

func (h *AuthHandler) QRReject(ctx context.Context, c *app.RequestContext) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}

	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	err := h.service.QRReject(req.Token)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "QR login rejected successfully",
	})
}

func (h *AuthHandler) GetQRLoginStatus(ctx context.Context, c *app.RequestContext) {
	token := c.Query("token")
	if token == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Token is required"))
		return
	}

	response, err := h.service.GetQRLoginStatus(token)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}
