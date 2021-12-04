package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type RegionI interface {
	Get(ctx context.Context, id string) (*models.Region, error)
	GetAll(ctx context.Context, page, limit uint32, name string) ([]*models.Region, uint32, error)
	GetAllByCity(ctx context.Context, cityID, name string) ([]*models.Region, uint32, error)
}
