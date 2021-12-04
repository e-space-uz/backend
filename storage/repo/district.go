package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type DistrictI interface {
	Get(ctx context.Context, id string) (*models.District, error)
	GetAll(ctx context.Context, page, limit, soato uint32, name string) ([]*models.District, uint32, error)
	GetAllByCityRegion(ctx context.Context, regionID, cityID, name string) ([]*models.District, uint32, error)
}
