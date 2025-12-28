package models

import "github.com/gofrs/uuid"


type Business struct {
	Base

	Name string `gorm:"type:varchar(255);not null"`
	UserID	   uuid.UUID `gorm:"type:uuid;not null;index"`

	Stores []Store

}


func (Business) TableName() string {
	return "businesses"
}
