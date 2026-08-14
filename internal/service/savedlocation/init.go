package savedlocation

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/savedlocation"
)

type Service interface {
	Create(ctx context.Context, userID, name string, lat, lng float64) (entity.SavedLocation, error)
	Get(ctx context.Context, id string) (entity.SavedLocation, error)
	ListByUserID(ctx context.Context, userID string) ([]entity.SavedLocation, error)
	Update(ctx context.Context, id, name string, lat, lng float64) (entity.SavedLocation, error)
	Delete(ctx context.Context, id string) error
}

type serviceImpl struct {
	repo savedlocation.Repository
}

// @WireSet("Service")
func New(repo savedlocation.Repository) Service {
	return &serviceImpl{repo: repo}
}
