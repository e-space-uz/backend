package repo

import "github.com/e-space-uz/backend/models"

type CityStorageI interface {
	Get(id string) (*models.City, error)
	GetAll(page, limit, soato uint32, name string) ([]*models.City, uint32, error)
}
