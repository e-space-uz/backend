package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type EntityI interface {
	// Read request
	Get(ctx context.Context, id string) (*models.Entity, error)
	GetAll(ctx context.Context, req *models.GetAllEntitiesRequest) ([]*models.GetAllEntities, uint64, error)
	GetAllWithProperties(ctx context.Context, req *models.GetAllEntitiesRequest) ([]*models.GetAllEntities, error)
	// Dashboard
	// Write request
	Create(ctx context.Context, req *models.CreateUpdateEntity) (string, error)
	Update(ctx context.Context, req *models.CreateUpdateEntity) error
	UpdateEntityStatus(ctx context.Context, entityID, statusID string) error
	UpdateEntityDrafts(ctx context.Context, entityID, entityDraftID string) error
	Delete(ctx context.Context, id string) error
}
