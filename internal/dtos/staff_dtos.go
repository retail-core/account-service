package dtos

import "github.com/gofrs/uuid"

type CreateStaffRequest struct {
	AdminID     uuid.UUID `json:"admin_id" binding:"required"` 
	UserID   	uuid.UUID `json:"user_id" binding:"required"` // user_id from auth service
	UserName    string    `json:"first_name" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Role        *string   `json:"role" binding:"default=staff"`
	Permissions []string  `json:"permissions" binding:"omitempty"`
}

type UserResponse struct {
	UserID        uuid.UUID `json:"id"`
}

type StaffResponse struct {
	Id     uuid.UUID `json:"id"`
	StoreID     uuid.UUID `json:"store_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Role       string    `json:"role"`
	IsVerified bool      `json:"is_verified"`
}