package models

import "github.com/gofrs/uuid"

type StaffPermission struct {
	StaffID      uuid.UUID `gorm:"primaryKey"`
	PermissionID uuid.UUID `gorm:"primaryKey"`
}

func (StaffPermission) TableName() string {
	return "staff_permissions"
}