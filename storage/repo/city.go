package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type CityI interface {
	Get(ctx context.Context, id string) (*models.City, error)
	GetAll(ctx context.Context, page, limit uint32) ([]*models.City, uint32, error)
}
