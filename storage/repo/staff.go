package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type StaffI interface {
	Create(ctx context.Context, req *models.Applicant) (string, error)
	GetAll(ctx context.Context, page, limit uint32) ([]*models.Applicant, uint32, error)
	Get(ctx context.Context, id string) (*models.Applicant, error)
	Update(ctx context.Context, req *models.Applicant) error
	LoginExists(ctx context.Context, login string) (bool, error)
	Login(ctx context.Context, login string) (*models.LoginInfo, error)
	Delete(ctx context.Context, id string) error
}
