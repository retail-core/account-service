package dtos

import "github.com/gofrs/uuid"

type BusinessCreatedEvent struct {
	OwnerId   uuid.UUID `json:"ownerId"`   // suppose to be owner_id but auth-serevice is sending ownerId
	OwnerName string    `json:"ownerName"` // suppose to be owner_name but auth-serevice is sending ownerName
}

type StaffCreatedEvent struct {
	OwnerID uuid.UUID `json:"owner_id"`
	Role    string    `json:"role"`
	UserID  uuid.UUID `json:"user_id"`
	StoreID uuid.UUID `json:"store_id"`
	Email   string    `json:"email"`
	UserName string    `json:"user_name"`
}

type VerifyStaffEvent struct {
	UserID uuid.UUID `json:"user_id"`
}

type EditStaffRequest struct {
	Role *string `json:"role" validate:"omitempty"`
	UserName *string `json:"username" validate:"omitempty"`
}