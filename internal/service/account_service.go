package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/retail-core/account-service/internal/dtos"
	"github.com/retail-core/account-service/internal/models"
)

type AccountService interface {
	CreateBusiness(ctx context.Context, req dtos.CreateBusinessRequest) (*models.Business, error)
	CreateStaff(ctx context.Context, storeID uuid.UUID, req dtos.CreateStaffRequest) (*models.Staff, error)
	CreateStore(ctx context.Context, businessID uuid.UUID, req dtos.CreateStoreRequest) (*models.Store, error)
	UpdateStore(ctx context.Context, userID uuid.UUID, storeID uuid.UUID, req dtos.CreateStoreRequest) (*models.Store, error)
	GetStoresByUserID(ctx context.Context, userID uuid.UUID) ([]models.Store, error)
	GetUsersByStoreID(ctx context.Context, storeID uuid.UUID) ([]uuid.UUID, error)
	GetStaffsByStoreID(ctx context.Context, storeID uuid.UUID) ([]models.Staff, error)
	VerifyStaff(ctx context.Context, userID uuid.UUID) error
	OnboardUser(ctx context.Context, userID uuid.UUID, req dtos.BusinessOnboardingRequest) (*models.BusinessOnboarding, error)
	DeleteStaff(ctx context.Context, storeID uuid.UUID, staffID uuid.UUID) error
	EditStaff(ctx context.Context, storeID uuid.UUID, staffID uuid.UUID, req dtos.EditStaffRequest) error
	GetStaffsByUserID(ctx context.Context, storeID uuid.UUID) ([]models.Staff, error)
}