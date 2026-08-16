package trip

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/trip"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	StartTrip(ctx context.Context, req *dto.StartTripRequest) (*dto.StartTripResponse, error)
	UpdatePosition(ctx context.Context, req *dto.UpdatePositionRequest) (*dto.PositionResponse, error)
	GetMemberPosition(ctx context.Context, req *dto.GetMemberPositionRequest) (*dto.PositionResponse, error)
	GetPartyPositions(ctx context.Context, req *dto.GetPartyPositionsRequest) (*dto.PositionsResponse, error)
}

type handlerImpl struct {
	service trip.Service
}

// @WireSet("Handler")
func New(service trip.Service) Handler {
	return &handlerImpl{service: service}
}
