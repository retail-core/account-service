package models

import (
	"github.com/gofrs/uuid"
)

type Store struct {
	Base

	BusinessID uuid.UUID `gorm:"index;not null"`

	Name string `gorm:"type:varchar(255);not null"`

	Address *string `gorm:"type:varchar(500);default:null"`

	Tag *string `gorm:"type:varchar(255);default:null"`

	Business Business

	Staffs []Staff `gorm:"constraint:OnDelete:CASCADE;"`
}

func (Store) TableName() string {
	return "stores"
}
