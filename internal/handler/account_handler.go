package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"github.com/retail-core/account-service/internal/dtos"
	"github.com/retail-core/account-service/internal/errors"
	"github.com/retail-core/account-service/internal/httpx"
	"github.com/retail-core/account-service/internal/mapping"
	"github.com/retail-core/account-service/internal/service"
	"github.com/retail-core/account-service/internal/validation"
)

type AccountHandler struct {
	Service service.AccountService
}

func NewAccountHandler(s service.AccountService) *AccountHandler {
	return &AccountHandler{Service: s}
}

func (h *AccountHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateStoreRequest
	userIDParam := chi.URLParam(r, "user_id")

	userID, err := uuid.FromString(userIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("user_id must be valid UUID string"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, errors.BadRequest("Invalid Request Body"))
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		httpx.WriteError(w, err)
		return
	}

	store, err := h.Service.CreateStore(r.Context(), userID, req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	response := map[string]any{
		"resource_id": store.ID,
		"message":     "Store created successfully",
	}

	httpx.WriteJSON(w, http.StatusCreated, response)
}

func (h *AccountHandler) CreateStaff(w http.ResponseWriter, r *http.Request) {
	storeIDParam := chi.URLParam(r, "store_id")

	storeID, err := uuid.FromString(storeIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("store_id must be valid UUID string"))
		return
	}
	var req dtos.CreateStaffRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, errors.BadRequest("Invalid Request Body"))
		return
	}

	staff, err := h.Service.CreateStaff(r.Context(), storeID, req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	response := map[string]any{
		"resource_id": staff.ID,
		"message":     "Staff created successfully",
	}

	httpx.WriteJSON(w, http.StatusOK, response)
}

func (h *AccountHandler) GetStoresByUserID(w http.ResponseWriter, r *http.Request) {
	userIDParam := chi.URLParam(r, "user_id")

	userID, err := uuid.FromString(userIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("user_id must be valid UUID string"))
		return
	}

	stores, err := h.Service.GetStoresByUserID(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	var res = make([]dtos.StoreResponse, 0, len(stores))
	for _, store := range stores {
		res = append(res, mapping.ToStoreResponse(store))
	}

	httpx.WriteJSON(w, http.StatusOK, res)
}

func (h *AccountHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateStoreRequest

	userIDParam := chi.URLParam(r, "user_id")
	storeIDParam := chi.URLParam(r, "store_id")

	userID, err := uuid.FromString(userIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("user_id must be valid UUID string"))
		return
	}

	storeID, err := uuid.FromString(storeIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("store_id must be valid UUID string"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, errors.BadRequest("Invalid Request Body"))
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		httpx.WriteError(w, err)
		return
	}

	store, err := h.Service.UpdateStore(r.Context(), userID, storeID, req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	response := map[string]any{
		"resource_id": store.ID,
		"message":     "Store updated successfully",
	}
	httpx.WriteJSON(w, http.StatusOK, response)
}

func (h *AccountHandler) GetUsersByStoreID(w http.ResponseWriter, r *http.Request) {
	storeIDParam := chi.URLParam(r, "store_id")

	storeID, err := uuid.FromString(storeIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("store_id must be valid UUID string"))
		return
	}

	userIDs, err := h.Service.GetUsersByStoreID(r.Context(), storeID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	var res = make([]dtos.UserResponse, 0, len(userIDs))

	for _, userID := range userIDs {
		res = append(res, dtos.UserResponse{UserID: userID})
	}

	httpx.WriteJSON(w, http.StatusOK, res)
}

func (h *AccountHandler) GetStaffsByStoreID(w http.ResponseWriter, r *http.Request) {
	storeIDParam := chi.URLParam(r, "store_id")

	storeID, err := uuid.FromString(storeIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("store_id must be valid UUID string"))
		return
	}

	userIDs, err := h.Service.GetStaffsByStoreID(r.Context(), storeID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	var res = make([]dtos.StaffResponse, 0, len(userIDs))

	for _, userID := range userIDs {
		res = append(res, dtos.StaffResponse{
			Id:         userID.ID,
			StoreID:    userID.StoreID,
			Username:   userID.UserName,
			Email:      userID.Email,
			Role:       *userID.Role,
			IsVerified: userID.IsVerified,
		})
	}

	httpx.WriteJSON(w, http.StatusOK, res)
}

func (h *AccountHandler) OnboardUser(w http.ResponseWriter, r *http.Request) {
	var req dtos.BusinessOnboardingRequest
	userIDParam := chi.URLParam(r, "user_id")

	userID, err := uuid.FromString(userIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("user_id must be valid UUID string"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.WriteError(w, errors.BadRequest("Invalid Request Body"))
		return
	}

	if err := validation.ValidateStruct(req); err != nil {
		httpx.WriteError(w, err)
		return
	}

	onboarding, err := h.Service.OnboardUser(r.Context(), userID, req)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	response := map[string]any{
		"resource_id": onboarding.ID,
		"message":     "User onboarded successfully",
	}

	httpx.WriteJSON(w, http.StatusCreated, response)
}

func (h *AccountHandler) DeleteStaff(w http.ResponseWriter, r *http.Request) {
	storeIDParam := chi.URLParam(r, "store_id")
	staffOIDParam := chi.URLParam(r, "staff_id")

	storeID, err := uuid.FromString(storeIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("store_id must be valid UUID string"))
		return
	}

	staffID, err := uuid.FromString(staffOIDParam)
	if err != nil {
		httpx.WriteError(w, errors.BadRequest("staff_id must be valid UUID string"))
		return
	}

	err = h.Service.DeleteStaff(r.Context(), storeID, staffID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	response := map[string]any{
		"message": "Staff deleted successfully",
	}
	httpx.WriteJSON(w, http.StatusOK, response)
}