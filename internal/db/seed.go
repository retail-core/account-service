package db

import (
	"gorm.io/gorm"
	"github.com/retail-core/account-service/internal/models"
)

func SeedPermissions(db *gorm.DB) error {
	permissions := []models.Permission{
		{Name: "INVENTORY_UPDATE"},
	}

	for _, p := range permissions {
		if err := db.
			FirstOrCreate(&models.Permission{}, models.Permission{Name: p.Name}).
			Error; err != nil {
			return err
		}
	}
	return nil
}
