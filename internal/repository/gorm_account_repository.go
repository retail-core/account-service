package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	appErrors "github.com/retail-core/account-service/internal/errors"
	"github.com/retail-core/account-service/internal/models"
	"gorm.io/gorm"
)

type GormAccountRepository struct {
	DB *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *GormAccountRepository {
	return &GormAccountRepository{DB: db}
}

func (r *GormAccountRepository) CreateBusiness(ctx context.Context, business *models.Business) (*models.Business, error) {
	if result := r.DB.WithContext(ctx).Create(business); result.Error != nil {
		return nil, fmt.Errorf("failed to create business: %w", result.Error)
	}
	return business, nil
}

func (r *GormAccountRepository) CreateStore(ctx context.Context, store *models.Store) (*models.Store, error) {
	if result := r.DB.WithContext(ctx).Create(store); result.Error != nil {
		return nil, fmt.Errorf("failed to create store: %w", result.Error)
	}
	return store, nil
}

func (r *GormAccountRepository) GetStoresByUserID(ctx context.Context, userID uuid.UUID) ([]models.Store, error) {
	// try to get business
	var business models.Business
	if err := r.DB.WithContext(ctx).Preload("Stores").First(&business, "user_id = ?", userID).Error; err == nil {
		return business.Stores, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get business by user_id: %w", err)
	}

	// fallback to staff
	var staff models.Staff
	if err := r.DB.WithContext(ctx).Preload("Store").First(&staff, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NotFound("User " + userID.String())
		}
		return nil, fmt.Errorf("failed to get staff by user ID: %w", err)
	}

	if staff.Store.ID == uuid.Nil {
		return nil, appErrors.NotFound("User " + userID.String())
	}

	return []models.Store{staff.Store}, nil
}

func (r *GormAccountRepository) SaveBusinessOnboardingData(ctx context.Context, answers *models.BusinessOnboarding) error {
	if result := r.DB.WithContext(ctx).Create(answers); result.Error != nil {
		return fmt.Errorf("failed to save business onboarding data: %w", result.Error)
	}
	return nil
}

func (r *GormAccountRepository) GetBusinessByUserID(ctx context.Context, userID uuid.UUID) (*models.Business, error) {
	var business models.Business

	if result := r.DB.WithContext(ctx).First(&business, "user_id = ?", userID); result.Error != nil {
		return nil, fmt.Errorf("failed to get business by user ID: %w", result.Error)
	}
	return &business, nil
}

func (r *GormAccountRepository) GetStoreByID(ctx context.Context, storeID uuid.UUID) (*models.Store, error) {
	var store models.Store

	if result := r.DB.WithContext(ctx).First(&store, "id = ?", storeID); result.Error != nil {
		return nil, fmt.Errorf("failed to get store by ID: %w", result.Error)
	}
	return &store, nil
}

func (r *GormAccountRepository) GetPermissionByName(ctx context.Context, name string) (*models.Permission, error) {
	var permission models.Permission

	if result := r.DB.WithContext(ctx).First(&permission, "name = ?", name); result.Error != nil {
		return nil, fmt.Errorf("failed to get permission by name: %w", result.Error)
	}
	return &permission, nil
}

func (r *GormAccountRepository) CreateStaff(ctx context.Context, staff *models.Staff) (*models.Staff, error) {
	if result := r.DB.WithContext(ctx).Create(staff); result.Error != nil {
		return nil, fmt.Errorf("failed to create staff: %w", result.Error)
	}
	return staff, nil
}

func (r *GormAccountRepository) UpdateStore(ctx context.Context, store *models.Store) (*models.Store, error) {
	if result := r.DB.WithContext(ctx).Save(store); result.Error != nil {
		return nil, fmt.Errorf("failed to update store: %w", result.Error)
	}
	return store, nil
}

func (r *GormAccountRepository) GetUsersByStoreID(ctx context.Context, storeID uuid.UUID) ([]uuid.UUID, error) {
	var userIDs = make([]uuid.UUID, 0)
	var store models.Store

	err := r.DB.WithContext(ctx).
		Model(&models.Store{}).
		Preload("Staffs").
		Preload("Business").
		Where("id = ?", storeID).
		First(&store).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get store by ID: %w", err)
	}

	userIDs = append(userIDs, store.Business.UserID) // add business owner user ID

	for _, staff := range store.Staffs {
		userIDs = append(userIDs, staff.ID)
	}

	return userIDs, nil
}

func (r *GormAccountRepository) GetStaffsByStoreID(ctx context.Context, storeID uuid.UUID) ([]models.Staff, error) {
	var staffs []models.Staff

	result := r.DB.WithContext(ctx).
		Model(&models.Staff{}).
		Where("store_id = ?", storeID).
		Find(&staffs)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []models.Staff{}, nil
		}

		return nil, fmt.Errorf("failed to get staffs by store ID: %w", result.Error)
	}
	return staffs, nil
}

func (r *GormAccountRepository) VerifyStaff(ctx context.Context, userID uuid.UUID) error {
	result := r.DB.WithContext(ctx).Model(&models.Staff{}).Where("user_id = ?", userID).Update("is_verified", true)
	if result.Error != nil {
		return fmt.Errorf("failed to verify staff: %w", result.Error)
	}
	return nil
}

func (r *GormAccountRepository) CreateBusinessOnboarding(ctx context.Context, onboarding *models.BusinessOnboarding) (*models.BusinessOnboarding, error) {
	if result := r.DB.WithContext(ctx).Create(onboarding); result.Error != nil {
		return nil, fmt.Errorf("failed to create business onboarding: %w", result.Error)
	}
	return onboarding, nil
}

func (r *GormAccountRepository) UpdateBusinessName(ctx context.Context, businessID uuid.UUID, name string) error {
	result := r.DB.WithContext(ctx).Model(&models.Business{}).Where("id = ?", businessID).Update("name", name)
	if result.Error != nil {
		return fmt.Errorf("failed to update business name: %w", result.Error)
	}
	return nil
}

func (r *GormAccountRepository) UpdateFirstStoreName(ctx context.Context, businessID uuid.UUID, name string) error {
	var store models.Store

	err := r.DB.WithContext(ctx).Where("business_id = ?", businessID).Order("created_at asc").First(&store).Error
	if err != nil {
		return fmt.Errorf("failed to get first store for business: %w", err)
	}

	result := r.DB.WithContext(ctx).Model(&models.Store{}).Where("id = ?", store.ID).Update("name", name)
	if result.Error != nil {
		return fmt.Errorf("failed to update first store name: %w", result.Error)
	}
	return nil
}

func (r *GormAccountRepository) DeleteStaff(ctx context.Context, storeID uuid.UUID, staffID uuid.UUID) (uuid.UUID, error) {
	var staff models.Staff

	err := r.DB.WithContext(ctx).
		Where("store_id = ? AND id = ?", storeID, staffID).
		First(&staff).Error

	if err != nil {
		return uuid.Nil, err
	}

	result := r.DB.WithContext(ctx).
		Unscoped().
		Where("store_id = ? AND id = ?", storeID, staffID).
		Delete(&models.Staff{})

	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	if result.RowsAffected == 0 {
		return uuid.Nil, nil
	}

	return staff.UserID, nil

}
