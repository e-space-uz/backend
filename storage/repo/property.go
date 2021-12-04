package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type PropertyI interface {
	Create(ctx context.Context, req *models.CreateUpdateProperty) (string, error)
	Get(ctx context.Context, id string) (*models.Property, error)
	GetAll(ctx context.Context, page, limit uint32, name string) ([]*models.Property, uint32, error)
	Update(ctx context.Context, req *models.CreateUpdateProperty) error
}
