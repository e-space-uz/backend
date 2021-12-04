package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type EntityFilesI interface {
	Create(ctx context.Context, req *models.CreateEntityFiles) (string, error)
	Get(ctx context.Context, id string) (*models.EntityFiles, error)
	GetAll(ctx context.Context, page, limit uint32, search string) ([]*models.EntityFiles, uint32, error)
	Update(ctx context.Context, req *models.EntityFiles) error
	Delete(ctx context.Context, id string) error
}
