package service

import (
	"encoding/json"
	"fmt"
	"kube/internal/config"
	"kube/internal/storage"
	"kube/pkg/errors"
	"kube/pkg/models"
	"kube/pkg/services"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type MapsService struct {
	*services.BaseService
	hereAPIKey string
	redis      *storage.RedisClient
}

func NewMapsService(db *gorm.DB) *MapsService {
	cfg := config.Load()
	redis := storage.NewRedisClient(cfg.Redis)

	return &MapsService{
		BaseService: services.NewBaseService(db),
		hereAPIKey:  cfg.HEREAPIKey,
		redis:       redis,
	}
}

func (s *MapsService) SearchPlaces(req *models.PlaceSearchRequest) (*models.PlaceSearchResponse, error) {
	// Generate cache key
	cacheKey := s.generateCacheKey(req)

	// Try to get from cache first
	if cachedResult, err := s.getFromCache(cacheKey); err == nil && cachedResult != nil {
		return cachedResult, nil
	}

	// Call HERE API
	hereResponse, err := s.callHereAPI(req)
	if err != nil {
		return nil, err
	}

	// Transform HERE API response to our format
	result := s.transformHereResponse(hereResponse, req)

	// Cache the result for 1 hour
	s.cacheResult(cacheKey, result)

	return result, nil
}

func (s *MapsService) callHereAPI(req *models.PlaceSearchRequest) (*models.HereAPIResponse, error) {
	baseURL := "https://discover.search.hereapi.com/v1/discover"

	params := url.Values{}
	params.Add("q", req.Query)
	params.Add("at", fmt.Sprintf("%.4f,%.4f", req.Lat, req.Lng))
	params.Add("limit", strconv.Itoa(req.Limit))
	params.Add("lang", req.Language)
	params.Add("apiKey", s.hereAPIKey)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeExternalAPIError, "Failed to call HERE API", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(errors.ErrCodeExternalAPIError, "HERE API returned error", fmt.Sprintf("Status code: %d", resp.StatusCode))
	}

	var hereResponse models.HereAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&hereResponse); err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternalError, "Failed to parse HERE API response", err.Error())
	}

	return &hereResponse, nil
}

func (s *MapsService) transformHereResponse(hereResponse *models.HereAPIResponse, req *models.PlaceSearchRequest) *models.PlaceSearchResponse {
	places := make([]models.PlaceItem, len(hereResponse.Items))

	for i, item := range hereResponse.Items {
		places[i] = models.PlaceItem{
			Title:      item.Title,
			ID:         item.ID,
			Language:   item.Language,
			ResultType: item.ResultType,
			Address: models.Address{
				Label:       item.Address.Label,
				CountryCode: item.Address.CountryCode,
				CountryName: item.Address.CountryName,
				County:      item.Address.County,
				City:        item.Address.City,
				District:    item.Address.District,
				Street:      item.Address.Street,
				PostalCode:  item.Address.PostalCode,
			},
			Position: models.Position{
				Lat: item.Position.Lat,
				Lng: item.Position.Lng,
			},
			Distance:   item.Distance,
			Categories: s.transformCategories(item.Categories),
			Contacts:   s.transformContacts(item.Contacts),
		}
	}

	return &models.PlaceSearchResponse{
		Places: places,
		Meta: models.SearchMeta{
			Query:        req.Query,
			ResultsCount: len(places),
			SearchTime:   time.Now(),
			FromCache:    false,
		},
	}
}

func (s *MapsService) transformCategories(categories []models.HereCategory) []models.Category {
	result := make([]models.Category, len(categories))
	for i, cat := range categories {
		result[i] = models.Category{
			ID:      cat.ID,
			Name:    cat.Name,
			Primary: cat.Primary,
		}
	}
	return result
}

func (s *MapsService) transformContacts(contacts []models.HereContact) []models.Contact {
	result := make([]models.Contact, len(contacts))
	for i, contact := range contacts {
		phones := make([]string, len(contact.Phone))
		for j, phone := range contact.Phone {
			phones[j] = phone.Value
		}

		websites := make([]string, len(contact.WWW))
		for j, www := range contact.WWW {
			websites[j] = www.Value
		}

		result[i] = models.Contact{
			Phone: phones,
			WWW:   websites,
		}
	}
	return result
}

func (s *MapsService) generateCacheKey(req *models.PlaceSearchRequest) string {
	return fmt.Sprintf("maps:search:%s:%.4f:%.4f:%d:%s",
		req.Query, req.Lat, req.Lng, req.Limit, req.Language)
}

func (s *MapsService) getFromCache(cacheKey string) (*models.PlaceSearchResponse, error) {
	if s.redis == nil {
		return nil, errors.New(errors.ErrCodeInternalError, "Redis not available", "Redis client not initialized")
	}

	cachedData, err := s.redis.Get(cacheKey)
	if err != nil {
		return nil, err
	}

	var cached models.CachedPlaceSearch
	if err := json.Unmarshal([]byte(cachedData), &cached); err != nil {
		return nil, err
	}

	// Check if cache is expired
	if time.Now().After(cached.ExpiresAt) {
		s.redis.Delete(cacheKey)
		return nil, errors.New(errors.ErrCodeCacheExpired, "Cache expired", "Cached data has expired")
	}

	return &models.PlaceSearchResponse{
		Places: cached.Results,
		Meta: models.SearchMeta{
			Query:        cached.Query,
			ResultsCount: len(cached.Results),
			SearchTime:   cached.CachedAt,
			FromCache:    true,
		},
	}, nil
}

func (s *MapsService) cacheResult(cacheKey string, result *models.PlaceSearchResponse) {
	if s.redis == nil {
		return
	}

	cached := models.CachedPlaceSearch{
		Query:     result.Meta.Query,
		Results:   result.Places,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour), // Cache for 1 hour
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return // Silently fail cache write
	}

	s.redis.SetWithExpiry(cacheKey, string(data), 1*time.Hour)
}
