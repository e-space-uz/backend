package repo

import "github.com/e-space-uz/backend/models"

type RegionStorageI interface {
	Get(id string) (*models.Region, error)
	GetAll(page, limit, soato uint32, name string) ([]*models.Region, uint32, error)
	GetAllByCity(cityID, name string) ([]*models.Region, uint32, error)
}
