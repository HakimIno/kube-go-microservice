package handler

import (
	"context"
	"kube/biz/service"
	"kube/pkg/errors"
	"kube/pkg/handlers"
	"kube/pkg/models"

	"github.com/cloudwego/hertz/pkg/app"
)

type AuthHandler struct {
	*handlers.BaseHandler
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		BaseHandler: handlers.NewBaseHandler(),
		service:     service,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.UserLoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse "Login successful"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data"
// @Failure 401 {object} models.ErrorResponse "Invalid credentials"
// @Failure 403 {object} models.ErrorResponse "Account deactivated"
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(ctx context.Context, c *app.RequestContext) {
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
// @Description Generate a new JWT token for the authenticated user. Requires valid Bearer token in Authorization header.
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.RefreshTokenResponse "Token refreshed successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Bearer token required or invalid"
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(ctx context.Context, c *app.RequestContext) {
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

// ChangePassword godoc
// @Summary Change user password
// @Description Change the password of the currently authenticated user. Requires valid Bearer token in Authorization header.
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body models.ChangePasswordRequest true "Password change data"
// @Success 200 {object} models.ChangePasswordResponse "Password changed successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Bearer token required or invalid current password"
// @Router /api/v1/auth/change-password [post]
func (h *AuthHandler) ChangePassword(ctx context.Context, c *app.RequestContext) {
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

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate session (client should remove token)
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.LogoutRequest false "Optional logout data"
// @Success 200 {object} models.LogoutResponse "Logged out successfully"
// @Failure 401 {object} models.ErrorResponse "Unauthorized - Bearer token required or invalid"
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(ctx context.Context, c *app.RequestContext) {
	userID, exists := c.Get("user_id")
	if !exists {
		errors.SendError(c, errors.New(errors.ErrCodeUnauthorized, "User not authenticated", "User ID not found in context"))
		return
	}

	var req models.LogoutRequest
	if err := c.BindJSON(&req); err != nil {
		// If binding fails, use empty request (optional data)
		req = models.LogoutRequest{}
	}

	if err := h.service.Logout(userID.(uint), &req); err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, nil, "Logged out successfully")
}

// GenerateQRCode godoc
// @Summary Generate QR code for login
// @Description Generate a QR code that can be scanned for login authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.QRCodeRequest false "Device information (optional)"
// @Success 200 {object} models.GenerateQRResponse "QR code generated successfully"
// @Failure 500 {object} models.ErrorResponse "Failed to generate QR code"
// @Router /api/v1/auth/qr/generate [post]
func (h *AuthHandler) GenerateQRCode(ctx context.Context, c *app.RequestContext) {
	var req models.QRCodeRequest
	if err := c.BindJSON(&req); err != nil {
		// If binding fails, use empty request (device_info is optional)
		req = models.QRCodeRequest{}
	}

	qrCode, err := h.service.GenerateQRCode(&req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, qrCode, "QR code generated successfully")
}

// QRConfirm godoc
// @Summary Mobile app approve QR code scan
// @Description Mobile app sends app token when scanning QR code and approves login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.QRConfirmRequest true "QR scan approval data from mobile app"
// @Success 200 {object} models.QRConfirmResponse "QR scan and login confirmed"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data"
// @Failure 401 {object} models.ErrorResponse "Invalid token or session"
// @Router /api/v1/auth/qr/confirm [post]
func (h *AuthHandler) QRConfirm(ctx context.Context, c *app.RequestContext) {
	var req models.QRConfirmRequest
	if err := c.BindJSON(&req); err != nil {
		h.SendValidationError(c, "Invalid request data format")
		return
	}

	// ตรวจสอบว่ามี app_token
	if req.AppToken == "" {
		h.SendValidationError(c, "app_token is required")
		return
	}

	err := h.service.QRConfirm(&req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, nil, "QR scan and login confirmed")
}

// QRReject godoc
// @Summary Mobile app reject QR code scan
// @Description Mobile app sends app token when scanning QR code and rejects login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.QRRejectRequest true "QR scan rejection data from mobile app"
// @Success 200 {object} models.QRRejectResponse "QR login rejected"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data"
// @Failure 401 {object} models.ErrorResponse "Invalid token or session"
// @Router /api/v1/auth/qr/reject [post]
func (h *AuthHandler) QRReject(ctx context.Context, c *app.RequestContext) {
	var req models.QRRejectRequest
	if err := c.BindJSON(&req); err != nil {
		h.SendValidationError(c, "Invalid request data format")
		return
	}

	// ตรวจสอบว่ามี app_token
	if req.AppToken == "" {
		h.SendValidationError(c, "app_token is required")
		return
	}

	err := h.service.QRReject(&req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, nil, "QR login rejected")
}

// GetQRLoginStatus godoc
// @Summary Get QR login status
// @Description Check the status of a QR login session (for polling after QR code scan)
// @Tags auth
// @Accept json
// @Produce json
// @Param session_id query string true "QR login session ID"
// @Success 200 {object} models.QRLoginStatusResponseWrapper "Session status retrieved"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid session ID"
// @Failure 404 {object} models.ErrorResponse "Session not found"
// @Router /api/v1/auth/qr/status [get]
func (h *AuthHandler) GetQRLoginStatus(ctx context.Context, c *app.RequestContext) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		h.SendValidationError(c, "session_id is required")
		return
	}

	status, err := h.service.GetQRLoginStatus(sessionID)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, status, "Session status retrieved")
}
