package repository

import (
	"ocean-ranch/database"
	"ocean-ranch/models"
)

type FeedRepository struct{}

func NewFeedRepository() *FeedRepository {
	return &FeedRepository{}
}

func (r *FeedRepository) Create(record *models.FeedRecord) error {
	return database.DB.Create(record).Error
}

func (r *FeedRepository) FindByCageID(cageID int) ([]models.FeedRecord, error) {
	var records []models.FeedRecord
	result := database.DB.Where("cage_id = ?", cageID).Order("created_at DESC").Find(&records)
	return records, result.Error
}

func (r *FeedRepository) FindAll() ([]models.FeedRecord, error) {
	var records []models.FeedRecord
	result := database.DB.Order("created_at DESC").Limit(50).Find(&records)
	return records, result.Error
}
