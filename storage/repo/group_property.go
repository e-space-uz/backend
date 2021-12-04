package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type GroupPropertyI interface {
	Create(ctx context.Context, req *models.CreateGroupProperty) (string, error)
	Get(ctx context.Context, id string) (*models.GroupProperty, error)
	GetAll(ctx context.Context, page, limit uint32, search string) ([]*models.GetAllGroupProperty, uint32, error)
	Update(ctx context.Context, req *models.CreateGroupProperty) error
	Delete(ctx context.Context, id string) error
	GetAllByType(ctx context.Context, typeOf, step uint32) ([]*models.GroupProperty, uint32, error)
}
