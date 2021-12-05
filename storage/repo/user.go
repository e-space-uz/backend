package repo

import (
	"context"

	"github.com/e-space-uz/backend/models"
)

type UserI interface {
	Create(ctx context.Context, req *models.User) (string, error)
	GetAll(ctx context.Context, page, limit uint32) ([]*models.User, uint32, error)
	Get(ctx context.Context, id string) (*models.User, error)
	Update(ctx context.Context, req *models.User) error
	LoginExists(ctx context.Context, login string) (bool, error)
	Login(ctx context.Context, login string) (*models.LoginInfo, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	GetCount(ctx context.Context, soato string, organizationId string) (int32, error)
	UpdatePassword(ctx context.Context, oldPassword, newPassword string, userID string) error
	SetRoleID(ctx context.Context, roleID, staffID string) error
	SetStaffSoato(ctx context.Context, staffID, soato string) error
}
