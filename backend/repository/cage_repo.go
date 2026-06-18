package repository

import (
	"ocean-ranch/database"
	"ocean-ranch/models"
)

type CageRepository struct{}

func NewCageRepository() *CageRepository {
	return &CageRepository{}
}

func (r *CageRepository) FindAll() ([]models.Cage, error) {
	var cages []models.Cage
	result := database.DB.Find(&cages)
	return cages, result.Error
}

func (r *CageRepository) FindByID(id int) (*models.Cage, error) {
	var cage models.Cage
	result := database.DB.First(&cage, id)
	return &cage, result.Error
}

func (r *CageRepository) Update(cage *models.Cage) error {
	return database.DB.Save(cage).Error
}
