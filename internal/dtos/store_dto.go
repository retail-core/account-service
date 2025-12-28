package dtos

type CreateStoreRequest struct {
	Name    string `json:"name" binding:"required"`
	Address *string `json:"address" binding:"omitempty"`
	Tag     *string `json:"tag" binding:"omitempty"`
}

type StoreResponse struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Address *string `json:"address,omitempty"`
	Tag     *string `json:"tag,omitempty"`
}