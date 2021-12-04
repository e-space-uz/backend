package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type EntityDraftI interface {
	Create(ctx context.Context, req *models.CreateEntityDraft) (string, error)
	Get(ctx context.Context, id string) (*models.EntityDraft, error)
	GetAll(ctx context.Context, req *models.GetAllEntityDraftsRequest) ([]*models.GetAllEntityDrafts, uint64, error)
}
