package repository

import (
	"kube/services/maps/internal/model"

	"gorm.io/gorm"
)

type MapsRepository struct {
	db *gorm.DB
}

func NewMapsRepository(db *gorm.DB) *MapsRepository {
	return &MapsRepository{db: db}
}

// SearchHistory methods
func (r *MapsRepository) CreateSearchHistory(history *model.SearchHistory) error {
	return r.db.Create(history).Error
}

func (r *MapsRepository) GetSearchHistoryByUserID(userID string, limit, offset int) ([]*model.SearchHistory, error) {
	var histories []*model.SearchHistory
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&histories).Error
	return histories, err
}

func (r *MapsRepository) DeleteSearchHistory(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.SearchHistory{}).Error
}

// FavoritePlace methods
func (r *MapsRepository) CreateFavoritePlace(favorite *model.FavoritePlace) error {
	return r.db.Create(favorite).Error
}

func (r *MapsRepository) GetFavoritePlacesByUserID(userID string, limit, offset int) ([]*model.FavoritePlace, error) {
	var favorites []*model.FavoritePlace
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&favorites).Error
	return favorites, err
}

func (r *MapsRepository) GetFavoritePlaceByUserIDAndPlaceID(userID, placeID string) (*model.FavoritePlace, error) {
	var favorite model.FavoritePlace
	err := r.db.Where("user_id = ? AND place_id = ?", userID, placeID).First(&favorite).Error
	return &favorite, err
}

func (r *MapsRepository) DeleteFavoritePlace(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.FavoritePlace{}).Error
}

func (r *MapsRepository) DeleteFavoritePlaceByUserIDAndPlaceID(userID, placeID string) error {
	return r.db.Where("user_id = ? AND place_id = ?", userID, placeID).Delete(&model.FavoritePlace{}).Error
}

// Check if place is already favorited
func (r *MapsRepository) IsPlaceFavorited(userID, placeID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.FavoritePlace{}).
		Where("user_id = ? AND place_id = ?", userID, placeID).
		Count(&count).Error
	return count > 0, err
}
