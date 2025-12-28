package models

import (
	"github.com/gofrs/uuid"
)

type Staff struct {
	Base
	BusinessID uuid.UUID `gorm:"index;not null"`
	UserName   string    `gorm:"type:varchar(100);not null"`
	Email      string    `json:"email" validate:"required,email"`
	StoreID    uuid.UUID `gorm:"index;not null"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"` // usere_id from auth service
	Role       *string   `gorm:"type:varchar(100);default:null"`
	IsVerified bool      `gorm:"default:false"` // whether the staff has verified their email

	// Associations

	Business Business
	Store    Store

	Permissions []Permission `gorm:"many2many:staff_permissions;"`
}
