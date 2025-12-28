package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/retail-core/account-service/internal/models"
)

type AccountRepository interface {
	CreateBusiness(ctx context.Context, business *models.Business) (*models.Business, error)

	CreateStore(ctx context.Context, store *models.Store) (*models.Store, error)

	UpdateStore(ctx context.Context, store *models.Store) (*models.Store, error)

	GetStoresByUserID(ctx context.Context, id uuid.UUID) ([]models.Store, error)

	SaveBusinessOnboardingData(ctx context.Context, answers *models.BusinessOnboarding) error

	GetBusinessByUserID(ctx context.Context, userID uuid.UUID) (*models.Business, error)

	GetStoreByID(ctx context.Context, storeID uuid.UUID) (*models.Store, error)
	
	GetPermissionByName(ctx context.Context, name string) (*models.Permission, error)
	CreateStaff(ctx context.Context, staff *models.Staff) (*models.Staff, error)

	GetUsersByStoreID(ctx context.Context, storeID uuid.UUID) ([]uuid.UUID, error)
	GetStaffsByStoreID(ctx context.Context, storeID uuid.UUID) ([]models.Staff, error)
	VerifyStaff(ctx context.Context, userID uuid.UUID) error
	CreateBusinessOnboarding(ctx context.Context, onboarding *models.BusinessOnboarding) (*models.BusinessOnboarding, error)
	UpdateBusinessName(ctx context.Context, businessID uuid.UUID, name string) error
	UpdateFirstStoreName(ctx context.Context, businessID uuid.UUID, name string) error
	DeleteStaff(ctx context.Context, storeID uuid.UUID, staffID uuid.UUID) (uuid.UUID, error)
}
