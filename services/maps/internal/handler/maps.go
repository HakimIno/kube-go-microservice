package handler

import (
	"context"
	"strconv"

	"kube/services/maps/internal/model"
	"kube/services/maps/internal/service"
	apperrors "kube/shared/pkg/errors"

	"github.com/cloudwego/hertz/pkg/app"
)

type MapsHandler struct {
	service *service.MapsService
}

func NewMapsHandler(service *service.MapsService) *MapsHandler {
	return &MapsHandler{service: service}
}

func (h *MapsHandler) SearchPlaces(ctx context.Context, c *app.RequestContext) {
	var req model.SearchPlacesRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	response, err := h.service.SearchPlaces(&req)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

func (h *MapsHandler) SearchPlacesByQuery(ctx context.Context, c *app.RequestContext) {
	query := c.Query("q")
	if query == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Query parameter 'q' is required"))
		return
	}

	// Parse optional parameters
	latitude := 0.0
	longitude := 0.0
	radius := 1000 // default radius in meters
	limit := 20    // default limit

	if latStr := c.Query("lat"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			latitude = lat
		}
	}

	if lngStr := c.Query("lng"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			longitude = lng
		}
	}

	if radiusStr := c.Query("radius"); radiusStr != "" {
		if r, err := strconv.Atoi(radiusStr); err == nil && r > 0 {
			radius = r
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	_ = c.Query("category") // category parameter not used in current implementation

	response, err := h.service.SearchPlacesByQuery(query, latitude, longitude, radius, limit)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

func (h *MapsHandler) GetPlaceDetails(ctx context.Context, c *app.RequestContext) {
	placeID := c.Param("id")
	if placeID == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Place ID is required"))
		return
	}

	details, err := h.service.GetPlaceDetails(placeID)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    details,
	})
}

func (h *MapsHandler) AddFavoritePlace(ctx context.Context, c *app.RequestContext) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperrors.HandleError(c, apperrors.NewUnauthorizedError("User ID not found"))
		return
	}

	var req model.AddFavoriteRequest
	if err := c.BindAndValidate(&req); err != nil {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Invalid request data"))
		return
	}

	err := h.service.AddFavoritePlace(userID, &req)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(201, map[string]interface{}{
		"success": true,
		"message": "Place added to favorites successfully",
	})
}

func (h *MapsHandler) GetFavoritePlaces(ctx context.Context, c *app.RequestContext) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperrors.HandleError(c, apperrors.NewUnauthorizedError("User ID not found"))
		return
	}

	// Get pagination parameters
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 20 // default limit
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

	response, err := h.service.GetFavoritePlaces(userID, limit, offset)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

func (h *MapsHandler) RemoveFavoritePlace(ctx context.Context, c *app.RequestContext) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperrors.HandleError(c, apperrors.NewUnauthorizedError("User ID not found"))
		return
	}

	placeID := c.Param("id")
	if placeID == "" {
		apperrors.HandleError(c, apperrors.NewBadRequestError("Place ID is required"))
		return
	}

	err := h.service.RemoveFavoritePlace(userID, placeID)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"message": "Place removed from favorites successfully",
	})
}

func (h *MapsHandler) GetSearchHistory(ctx context.Context, c *app.RequestContext) {
	userID := c.GetString("user_id")
	if userID == "" {
		apperrors.HandleError(c, apperrors.NewUnauthorizedError("User ID not found"))
		return
	}

	// Get pagination parameters
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	limit := 20 // default limit
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

	histories, err := h.service.GetSearchHistory(userID, limit, offset)
	if err != nil {
		apperrors.HandleError(c, err)
		return
	}

	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    histories,
		"pagination": map[string]interface{}{
			"limit":  limit,
			"offset": offset,
		},
	})
}
