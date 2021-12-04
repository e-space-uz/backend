package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type CityI interface {
	Get(ctx context.Context, id string) (*models.City, error)
	GetAll(ctx context.Context, page, limit, soato uint32, name string) ([]*models.City, uint32, error)
}
