package models

import "github.com/gofrs/uuid"

type BusinessOnboarding struct {
	Base
	BusinessID uuid.UUID `gorm:"uniqueIndex"`
	BusinessCategory *string `gorm:"type:varchar(100);default:null"`
	BusinessName	 *string `gorm:"type:varchar(255);default:null"`

	HasStaff          *bool
	HasMultipleStores *bool
}

func (BusinessOnboarding) TableName() string {
	return "business_onboardings"
}