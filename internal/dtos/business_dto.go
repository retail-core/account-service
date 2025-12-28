package dtos

import "github.com/gofrs/uuid"

type CreateBusinessRequest struct {
	UserID uuid.UUID `json:"ownerId" binding:"required"`
	Name   string `json:"ownerName" binding:"required"`
}

type BusinessOnboardingRequest struct {
	BusinessCategory  *string `json:"business_category"`
	BusinessName      *string `json:"business_name"`
	HasStaff          *bool   `json:"has_staff"`
	HasMultipleStores *bool   `json:"has_multiple_stores"`
}
