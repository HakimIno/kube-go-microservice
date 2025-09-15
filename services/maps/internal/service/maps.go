package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"kube/services/maps/internal/model"
	"kube/services/maps/internal/repository"
	apperrors "kube/shared/pkg/errors"
	"kube/shared/pkg/services"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MapsService struct {
	*services.BaseService
	repo        *repository.MapsRepository
	hereAPIKey  string
	hereBaseURL string
	redisClient *redis.Client
}

func NewMapsService(db *gorm.DB) *MapsService {
	hereAPIKey := os.Getenv("HERE_API_KEY")
	if hereAPIKey == "" {
		hereAPIKey = "YOUR_HERE_API_KEY" // fallback for development
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &MapsService{
		BaseService: services.NewBaseService(db),
		repo:        repository.NewMapsRepository(db),
		hereAPIKey:  hereAPIKey,
		hereBaseURL: "https://discover.search.hereapi.com/v1/discover",
		redisClient: redisClient,
	}
}

func (s *MapsService) SearchPlaces(req *model.SearchPlacesRequest) (*model.SearchPlacesResponse, error) {
	// Generate cache key from search parameters
	cacheKey := s.generateCacheKey(req)

	// Try to get from cache first
	if cachedResponse, err := s.getFromCache(cacheKey); err == nil {
		if os.Getenv("GO_ENV") == "development" {
			fmt.Printf("Cache HIT for key: %s\n", cacheKey)
		}
		return cachedResponse, nil
	}

	if os.Getenv("GO_ENV") == "development" {
		fmt.Printf("Cache MISS for key: %s\n", cacheKey)
	}

	// Build HERE API URL
	apiURL := s.buildHereAPIURL(req)

	// Make request to HERE API
	places, err := s.fetchPlacesFromHereAPI(apiURL)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to fetch places from HERE API")
	}

	// Create response
	response := &model.SearchPlacesResponse{
		Places: places,
		Total:  len(places),
		Query:  req.Query,
	}

	// Cache the response for 30 minutes
	s.setCache(cacheKey, response, 30*time.Minute)

	// Save search history if user is authenticated
	if req.Latitude != 0 && req.Longitude != 0 {
		// This would typically come from the authenticated user context
		// For now, we'll skip saving search history
	}

	return response, nil
}

func (s *MapsService) SearchPlacesByQuery(query string, latitude, longitude float64, radius, limit int) (*model.SearchPlacesResponse, error) {
	req := &model.SearchPlacesRequest{
		Query:     query,
		Latitude:  latitude,
		Longitude: longitude,
		Radius:    radius,
		Limit:     limit,
	}

	return s.SearchPlaces(req)
}

func (s *MapsService) GetPlaceDetails(placeID string) (*model.PlaceDetails, error) {
	// This would typically fetch detailed information from HERE API
	// For now, return a mock response
	return &model.PlaceDetails{
		ID:           placeID,
		Title:        "Sample Place",
		Address:      "123 Main St, City, Country",
		Category:     "restaurant",
		Latitude:     40.7128,
		Longitude:    -74.0060,
		Phone:        "+1-555-0123",
		Website:      "https://example.com",
		OpeningHours: []string{"Mon-Fri: 9:00-17:00", "Sat-Sun: 10:00-16:00"},
		Rating:       4.5,
		PriceLevel:   2,
		Description:  "A sample place for demonstration",
		Photos:       []string{"https://example.com/photo1.jpg"},
		Reviews:      []model.Review{},
	}, nil
}

func (s *MapsService) AddFavoritePlace(userID string, req *model.AddFavoriteRequest) error {
	// Check if place is already favorited
	isFavorited, err := s.repo.IsPlaceFavorited(userID, req.PlaceID)
	if err != nil {
		return apperrors.NewInternalServerError("Failed to check favorite status")
	}

	if isFavorited {
		return apperrors.NewBadRequestError("Place is already in favorites")
	}

	favorite := &model.FavoritePlace{
		UserID:    userID,
		PlaceID:   req.PlaceID,
		Title:     req.Title,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Category:  req.Category,
	}

	if err := s.repo.CreateFavoritePlace(favorite); err != nil {
		return apperrors.NewInternalServerError("Failed to add favorite place")
	}

	return nil
}

func (s *MapsService) GetFavoritePlaces(userID string, limit, offset int) (*model.GetFavoritesResponse, error) {
	favorites, err := s.repo.GetFavoritePlacesByUserID(userID, limit, offset)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to get favorite places")
	}

	// Convert []*FavoritePlace to []FavoritePlace
	var favoriteList []model.FavoritePlace
	for _, fav := range favorites {
		favoriteList = append(favoriteList, *fav)
	}

	return &model.GetFavoritesResponse{
		Favorites: favoriteList,
		Total:     len(favoriteList),
	}, nil
}

func (s *MapsService) RemoveFavoritePlace(userID, placeID string) error {
	err := s.repo.DeleteFavoritePlaceByUserIDAndPlaceID(userID, placeID)
	if err != nil {
		return apperrors.NewInternalServerError("Failed to remove favorite place")
	}

	return nil
}

func (s *MapsService) GetSearchHistory(userID string, limit, offset int) ([]*model.SearchHistory, error) {
	histories, err := s.repo.GetSearchHistoryByUserID(userID, limit, offset)
	if err != nil {
		return nil, apperrors.NewInternalServerError("Failed to get search history")
	}

	return histories, nil
}

func (s *MapsService) buildHereAPIURL(req *model.SearchPlacesRequest) string {
	baseURL := s.hereBaseURL
	params := url.Values{}

	params.Add("apikey", s.hereAPIKey)
	params.Add("q", req.Query)

	// HERE API requires location parameter - use provided coordinates or default to Bangkok
	latitude := req.Latitude
	longitude := req.Longitude

	// If no coordinates provided, default to Bangkok, Thailand
	if latitude == 0 && longitude == 0 {
		latitude = 13.7563   // Bangkok latitude
		longitude = 100.5018 // Bangkok longitude
	}

	// Add location parameter
	params.Add("at", fmt.Sprintf("%.6f,%.6f", latitude, longitude))

	// Add radius if specified
	if req.Radius > 0 {
		params.Add("in", fmt.Sprintf("circle:%s;r=%d",
			fmt.Sprintf("%.6f,%.6f", latitude, longitude), req.Radius))
	}

	if req.Limit > 0 {
		params.Add("limit", strconv.Itoa(req.Limit))
	}

	if req.Category != "" {
		params.Add("categories", req.Category)
	}

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}

func (s *MapsService) fetchPlacesFromHereAPI(apiURL string) ([]model.Place, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Log request in development only
	if os.Getenv("GO_ENV") == "development" {
		fmt.Printf("Making request to HERE API: %s\n", apiURL)
	}
	
	resp, err := client.Get(apiURL)
	if err != nil {
		if os.Getenv("GO_ENV") == "development" {
			fmt.Printf("Error making request to HERE API: %v\n", err)
		}
		return nil, err
	}
	defer resp.Body.Close()

	if os.Getenv("GO_ENV") == "development" {
		fmt.Printf("HERE API response status: %d\n", resp.StatusCode)
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		if os.Getenv("GO_ENV") == "development" {
			fmt.Printf("HERE API error response: %s\n", string(body))
		}
		return nil, fmt.Errorf("HERE API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var hereResponse struct {
		Items []struct {
			ID      string `json:"id"`
			Title   string `json:"title"`
			Address struct {
				Label string `json:"label"`
			} `json:"address"`
			Categories []struct {
				Name string `json:"name"`
			} `json:"categories"`
			Distance int `json:"distance"`
			Position struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"position"`
			Contacts []struct {
				Phone []struct {
					Value string `json:"value"`
				} `json:"phone"`
				Website []struct {
					Value string `json:"value"`
				} `json:"www"`
			} `json:"contacts"`
			OpeningHours []struct {
				Text []string `json:"text"`
			} `json:"openingHours"`
			FoodTypes []struct {
				Name string `json:"name"`
			} `json:"foodTypes"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &hereResponse); err != nil {
		return nil, err
	}

	var places []model.Place
	for _, item := range hereResponse.Items {
		place := model.Place{
			ID:        item.ID,
			Title:     item.Title,
			Address:   item.Address.Label,
			Distance:  item.Distance,
			Latitude:  item.Position.Lat,
			Longitude: item.Position.Lng,
		}

		// Set category
		if len(item.Categories) > 0 {
			place.Category = item.Categories[0].Name
		}

		// Set phone
		if len(item.Contacts) > 0 && len(item.Contacts[0].Phone) > 0 {
			place.Phone = item.Contacts[0].Phone[0].Value
		}

		// Set website
		if len(item.Contacts) > 0 && len(item.Contacts[0].Website) > 0 {
			place.Website = item.Contacts[0].Website[0].Value
		}

		// Set opening hours
		if len(item.OpeningHours) > 0 {
			place.OpeningHours = item.OpeningHours[0].Text
		}

		places = append(places, place)
	}

	return places, nil
}

// generateCacheKey creates a unique cache key from search parameters
func (s *MapsService) generateCacheKey(req *model.SearchPlacesRequest) string {
	// Create a string representation of the search parameters
	keyData := fmt.Sprintf("%s:%.6f:%.6f:%d:%d:%s",
		req.Query,
		req.Latitude,
		req.Longitude,
		req.Radius,
		req.Limit,
		req.Category,
	)

	// Generate MD5 hash for consistent key length
	hash := md5.Sum([]byte(keyData))
	return fmt.Sprintf("maps:search:%x", hash)
}

// getFromCache retrieves cached search results
func (s *MapsService) getFromCache(key string) (*model.SearchPlacesResponse, error) {
	ctx := context.Background()

	cachedData, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var response model.SearchPlacesResponse
	if err := json.Unmarshal([]byte(cachedData), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// setCache stores search results in cache
func (s *MapsService) setCache(key string, response *model.SearchPlacesResponse, expiration time.Duration) {
	ctx := context.Background()

	data, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error marshaling cache data: %v\n", err)
		return
	}

	if err := s.redisClient.Set(ctx, key, data, expiration).Err(); err != nil {
		fmt.Printf("Error setting cache: %v\n", err)
	}
}
