package models

type TPermission string;

const (
	INVENTORY_UPDATE  TPermission = "UPDATE_INVENTORY"
)

type Permission struct {
	Base
	Name       string `gorm:"type:varchar(100);not null;uniqueIndex"` // pls make it unique
}

func (Permission) TableName() string {
	return "permissions"
}