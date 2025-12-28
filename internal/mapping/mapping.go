package mapping

import (
	"github.com/retail-core/account-service/internal/dtos"
	"github.com/retail-core/account-service/internal/models"
)

func ToStoreResponse(store models.Store) dtos.StoreResponse {
	return dtos.StoreResponse{
		ID:      store.ID.String(),
		Name:    store.Name,
		Address: store.Address,
		Tag:     store.Tag,
	}
}