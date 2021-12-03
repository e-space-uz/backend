package repo

import (
	"github.com/e-space-uz/backend/models"
)

type PropertyI interface {
	Create(req *models.CreateUpdateProperty) (string, error)
	Get(id string) (*models.Property, error)
	GetAll(page, limit uint32, name string) ([]*models.Property, uint32, error)
	Update(req *models.CreateUpdateProperty) error
	Delete(id string) error
	PropertyExists(id string) (bool, error)
}
