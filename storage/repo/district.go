package repo

import "github.com/e-space-uz/backend/models"

type DistrictStorageI interface {
	Get(id string) (*models.District, error)
	GetAll(page, limit, soato uint32, name string) ([]*models.District, uint32, error)
	GetAllByCityRegion(regionID, cityID, name string) ([]*models.District, uint32, error)
}
