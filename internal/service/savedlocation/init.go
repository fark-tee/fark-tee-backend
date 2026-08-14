package savedlocation

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/savedlocation"
)

type Service interface {
	Create(ctx context.Context, userID, name string, lat, lng float64) (entity.SavedLocation, error)
	// Get returns the saved location if it belongs to userID.
	Get(ctx context.Context, userID, id string) (entity.SavedLocation, error)
	ListByUserID(ctx context.Context, userID string) ([]entity.SavedLocation, error)
	// Update updates the saved location if it belongs to userID.
	Update(ctx context.Context, userID, id, name string, lat, lng float64) (entity.SavedLocation, error)
	// Delete deletes the saved location if it belongs to userID.
	Delete(ctx context.Context, userID, id string) error
}

type serviceImpl struct {
	repo savedlocation.Repository
}

// @WireSet("Service")
func New(repo savedlocation.Repository) Service {
	return &serviceImpl{repo: repo}
}
