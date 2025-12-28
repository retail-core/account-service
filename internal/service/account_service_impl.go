package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/retail-core/account-service/internal/dtos"
	"github.com/retail-core/account-service/internal/logger"
	"github.com/retail-core/account-service/internal/models"
	"github.com/retail-core/account-service/internal/mq"
	"github.com/retail-core/account-service/internal/repository"
	"go.uber.org/zap"
)

const (
	StaffDeletedRoutingKey = "staff.deleted"
)

type AccountServiceImpl struct {
	Repo        repository.AccountRepository
	mqPublisher *mq.Publisher
}

func NewAccountServiceImpl(repo repository.AccountRepository, publisher *mq.Publisher) *AccountServiceImpl {
	return &AccountServiceImpl{
		Repo:        repo,
		mqPublisher: publisher,
	}
}

func (s *AccountServiceImpl) CreateBusiness(ctx context.Context, req dtos.CreateBusinessRequest) (*models.Business, error) {
	business := &models.Business{
		UserID: req.UserID,
		Name:   req.Name,
	}

	createdBusiness, err := s.Repo.CreateBusiness(ctx, business)
	if err != nil {
		logger.L().Error("Failed to create business", zap.Error(err))
		return nil, err
	}

	// create a default store for the business userName + Shop
	store := &models.Store{
		BusinessID: createdBusiness.ID,
		Name:       req.Name + " Shop",
	}

	_, err = s.Repo.CreateStore(ctx, store)
	if err != nil {
		logger.L().Error("Failed to create default store for business", zap.Error(err))
		return nil, err
	}

	return createdBusiness, nil
}

func (s *AccountServiceImpl) CreateStore(ctx context.Context, userID uuid.UUID, req dtos.CreateStoreRequest) (*models.Store, error) {

	business, err := s.Repo.GetBusinessByUserID(ctx, userID)
	if err != nil {
		logger.L().Error("Failed to fetch business for user", zap.Error(err))
		return nil, err
	}

	store := &models.Store{
		BusinessID: business.ID,
		Name:       req.Name,
		Address:    req.Address,
		Tag:        req.Tag,
	}

	createdStore, err := s.Repo.CreateStore(ctx, store)
	if err != nil {
		logger.L().Error("Failed to create store", zap.Error(err))
		return nil, err
	}

	return createdStore, nil
}

func (s *AccountServiceImpl) CreateStaff(ctx context.Context, storeID uuid.UUID, req dtos.CreateStaffRequest) (*models.Staff, error) {

	business, err := s.Repo.GetBusinessByUserID(ctx, req.AdminID)
	if err != nil {
		logger.L().Error("Failed to fetch business for user", zap.Error(err))
		return nil, err
	}

	store, err := s.Repo.GetStoreByID(ctx, storeID)
	if err != nil {
		logger.L().Error("Failed to fetch store for business", zap.Error(err))
		return nil, err
	}

	permissions := make([]models.Permission, 0)

	for _, permName := range req.Permissions {
		perm, err := s.Repo.GetPermissionByName(ctx, permName)
		if err != nil {
			logger.L().Error("Failed to fetch permission by name", zap.Error(err))
		}
		permissions = append(permissions, *perm)
	}

	staff := &models.Staff{
		StoreID:     store.ID,
		BusinessID:  business.ID,
		UserName:    req.UserName,
		Email:       req.Email,
		Role:        req.Role,
		UserID:      req.UserID,
		Permissions: permissions,
	}

	staff, err = s.Repo.CreateStaff(ctx, staff)
	if err != nil {
		logger.L().Error("Failed to create staff", zap.Error(err))
		return nil, err
	}

	return staff, nil
}

func (s *AccountServiceImpl) GetStoresByUserID(ctx context.Context, userID uuid.UUID) ([]models.Store, error) {
	stores, err := s.Repo.GetStoresByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (s *AccountServiceImpl) UpdateStore(ctx context.Context, userID uuid.UUID, storeID uuid.UUID, req dtos.CreateStoreRequest) (*models.Store, error) {
	_, err := s.Repo.GetBusinessByUserID(ctx, userID) // Verify business exists for user, this means only business owner can update store

	if err != nil {
		logger.L().Error("Failed to fetch business for user", zap.Error(err))
		return nil, err
	}

	store, err := s.Repo.GetStoreByID(ctx, storeID)
	if err != nil {
		logger.L().Error("Failed to fetch store for business", zap.Error(err))
		return nil, err
	}

	store.Name = req.Name
	store.Address = req.Address
	store.Tag = req.Tag

	updatedStore, err := s.Repo.UpdateStore(ctx, store)
	if err != nil {
		logger.L().Error("Failed to update store ", zap.Error(err))
		return nil, err
	}

	return updatedStore, nil
}

func (s *AccountServiceImpl) GetUsersByStoreID(ctx context.Context, storeID uuid.UUID) ([]uuid.UUID, error) {
	userIDs, err := s.Repo.GetUsersByStoreID(ctx, storeID)

	if err != nil {
		logger.L().Error("Failed to fetch users for store", zap.Error(err))
		return nil, err
	}
	return userIDs, nil
}

func (s *AccountServiceImpl) GetStaffsByStoreID(ctx context.Context, storeID uuid.UUID) ([]models.Staff, error) {
	staffs, err := s.Repo.GetStaffsByStoreID(ctx, storeID)
	if err != nil {
		logger.L().Error("Failed to fetch staffs for store", zap.Error(err))
		return nil, err
	}

	return staffs, nil
}

func (s *AccountServiceImpl) VerifyStaff(ctx context.Context, userID uuid.UUID) error {
	if err := s.Repo.VerifyStaff(ctx, userID); err != nil {
		logger.L().Error("Failed to fetch staff for verification", zap.Error(err))
		return err
	}
	return nil
}

func (s *AccountServiceImpl) OnboardUser(ctx context.Context, userID uuid.UUID, req dtos.BusinessOnboardingRequest) (*models.BusinessOnboarding, error) {

	// first we get the business associated with the user
	business, err := s.Repo.GetBusinessByUserID(ctx, userID)
	if err != nil {
		logger.L().Error("Failed to fetch business for user", zap.Error(err))
		return nil, err
	}

	if req.BusinessName != nil {
		// update business name
		err = s.Repo.UpdateBusinessName(ctx, business.ID, *req.BusinessName)
		if err != nil {
			logger.L().Error("Failed to update business name", zap.Error(err))
			return nil, err
		}

		// update first store name
		err = s.Repo.UpdateFirstStoreName(ctx, business.ID, *req.BusinessName)
		if err != nil {
			logger.L().Error("Failed to update first store name", zap.Error(err))
			return nil, err
		}
	}

	onboarding := &models.BusinessOnboarding{
		BusinessID:        business.ID,
		BusinessCategory:  req.BusinessCategory,
		BusinessName:      req.BusinessName,
		HasStaff:          req.HasStaff,
		HasMultipleStores: req.HasMultipleStores,
	}

	createdOnboarding, err := s.Repo.CreateBusinessOnboarding(ctx, onboarding)

	if err != nil {
		logger.L().Error("Failed to create business onboarding", zap.Error(err))
		return nil, err
	}

	return createdOnboarding, nil

}

func (s *AccountServiceImpl) DeleteStaff(ctx context.Context, storeID uuid.UUID, staffID uuid.UUID) error {

	userID, err := s.Repo.DeleteStaff(ctx, storeID, staffID)

	if err != nil {
		logger.L().Error("Failed to delete staff", zap.Error(err))
		return err
	}

	if err := s.mqPublisher.PublishDomainEvent(ctx, StaffDeletedRoutingKey, map[string]uuid.UUID{
		"user_id": userID,
	}); err != nil {
		logger.L().Error("Failed to publish staff deleted event", zap.Error(err))
		return err
	}

	return nil
}
