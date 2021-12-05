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
	Exists(ctx context.Context, id string) (bool, error)
	GetCount(ctx context.Context, soato string, organizationId string) (int32, error)
	UpdatePassword(ctx context.Context, oldPassword, newPassword string, ApplicantID string) error
	SetRoleID(ctx context.Context, roleID, staffID string) error
	SetStaffSoato(ctx context.Context, staffID, soato string) error
}
