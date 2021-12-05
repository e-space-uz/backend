package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type ApplicantI interface {
	Create(ctx context.Context, req *models.Applicant) (string, error)
	GetAll(ctx context.Context, page, limit uint32) ([]*models.Applicant, uint32, error)
	Get(ctx context.Context, id string) (*models.Applicant, error)
	Update(ctx context.Context, req *models.Applicant) error
}
