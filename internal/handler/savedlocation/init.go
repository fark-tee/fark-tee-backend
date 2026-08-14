package savedlocation

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/savedlocation"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	Create(ctx context.Context, req *dto.CreateSavedLocationRequest) (*dto.SavedLocationResponse, error)
	Get(ctx context.Context, req *dto.GetSavedLocationRequest) (*dto.SavedLocationResponse, error)
	List(ctx context.Context, req *dto.ListSavedLocationsRequest) (*dto.SavedLocationsResponse, error)
	Update(ctx context.Context, req *dto.UpdateSavedLocationRequest) (*dto.SavedLocationResponse, error)
	Delete(ctx context.Context, req *dto.DeleteSavedLocationRequest) (*dto.DeleteSavedLocationResponse, error)
}

type handlerImpl struct {
	service savedlocation.Service
}

// @WireSet("Handler")
func New(service savedlocation.Service) Handler {
	return &handlerImpl{service: service}
}
