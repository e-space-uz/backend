package repo

import "github.com/e-space-uz/backend/models"

type GroupPropertyI interface {
	Create(req *models.CreateGroupProperty) (string, error)
	Get(id string) (*models.GroupProperty, error)
	GetAll(page, limit uint32, search string) ([]*models.GetAllGroupProperty, uint32, error)
	Update(req *models.CreateGroupProperty) error
	Delete(id string) error
	GetAllByType(typeOf, step uint32) ([]*models.GroupProperty, uint32, error)
}
