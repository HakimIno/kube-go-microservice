package handler

import (
	"context"
	"kube/biz/service"
	"kube/pkg/errors"
	"kube/pkg/handlers"
	"kube/pkg/models"

	"github.com/cloudwego/hertz/pkg/app"
)

type MapsHandler struct {
	*handlers.BaseHandler
	service *service.MapsService
}

func NewMapsHandler(service *service.MapsService) *MapsHandler {
	return &MapsHandler{
		BaseHandler: handlers.NewBaseHandler(),
		service:     service,
	}
}

// SearchPlaces godoc
// @Summary Search for places using HERE API
// @Description Search for places near a specific location using HERE API with caching to reduce costs
// @Tags maps
// @Accept json
// @Produce json
// @Param request body models.PlaceSearchRequest true "Search parameters"
// @Success 200 {object} models.PlaceSearchResponse "Places found successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request data"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/maps/search [post]
func (h *MapsHandler) SearchPlaces(ctx context.Context, c *app.RequestContext) {
	var req models.PlaceSearchRequest
	if err := c.BindJSON(&req); err != nil {
		h.SendValidationError(c, "Invalid request data format")
		return
	}

	// Validate required fields
	if req.Query == "" {
		h.SendValidationError(c, "Query is required")
		return
	}

	// Set defaults
	if req.Limit == 0 {
		req.Limit = 5
	}
	if req.Language == "" {
		req.Language = "th"
	}

	result, err := h.service.SearchPlaces(&req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, result, "Places found successfully")
}

// SearchPlacesByQuery godoc
// @Summary Search for places by query string
// @Description Search for places using query string with default location (Bangkok)
// @Tags maps
// @Accept json
// @Produce json
// @Param q query string true "Search query" example:"วัดพระแก้ว กรุงเทพมหานคร"
// @Param lat query number false "Latitude" example:"13.7563"
// @Param lng query number false "Longitude" example:"100.5018"
// @Param limit query int false "Number of results" example:"5"
// @Param lang query string false "Language code" example:"th"
// @Success 200 {object} models.PlaceSearchResponse "Places found successfully"
// @Failure 400 {object} models.ValidationErrorResponse "Invalid request parameters"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /api/v1/maps/search/query [get]
func (h *MapsHandler) SearchPlacesByQuery(ctx context.Context, c *app.RequestContext) {
	query := string(c.Query("q"))
	if query == "" {
		h.SendValidationError(c, "Query parameter 'q' is required")
		return
	}

	// Parse other parameters with defaults
	req := &models.PlaceSearchRequest{
		Query:    query,
		Lat:      13.7563, // Default to Bangkok
		Lng:      100.5018,
		Limit:    5,
		Language: "th",
	}

	// Override with query parameters if provided
	if latStr := string(c.Query("lat")); latStr != "" {
		if lat, err := h.GetQueryFloat(c, "lat"); err == nil {
			req.Lat = lat
		}
	}

	if lngStr := string(c.Query("lng")); lngStr != "" {
		if lng, err := h.GetQueryFloat(c, "lng"); err == nil {
			req.Lng = lng
		}
	}

	if limitStr := string(c.Query("limit")); limitStr != "" {
		if limit, err := h.GetQueryInt(c, "limit"); err == nil {
			req.Limit = limit
		}
	}

	if lang := string(c.Query("lang")); lang != "" {
		req.Language = lang
	}

	result, err := h.service.SearchPlaces(req)
	if err != nil {
		errors.SendError(c, err)
		return
	}

	h.SendSuccess(c, 200, result, "Places found successfully")
}
