package trip

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/middleware/authmw"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) StartTrip(ctx context.Context, req *dto.StartTripRequest) (*dto.StartTripResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	createdTrip, createdPosition, err := h.service.StartTrip(
		ctx,
		actorID,
		req.PartyID,
		entity.TripDirection(req.Body.Direction),
		req.Body.Lat,
		req.Body.Lng,
		entity.Destination{
			Name: req.Body.DestinationName,
			Lat:  req.Body.DestinationLat,
			Lng:  req.Body.DestinationLng,
		},
	)
	if err != nil {
		return nil, err
	}

	return &dto.StartTripResponse{
		Trip:     *toTripResponse(createdTrip),
		Position: *toPositionResponse(createdPosition),
	}, nil
}

func (h *handlerImpl) UpdatePosition(ctx context.Context, req *dto.UpdatePositionRequest) (*dto.PositionResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	updated, err := h.service.UpdatePosition(ctx, actorID, req.PartyID, req.Body.Lat, req.Body.Lng)
	if err != nil {
		return nil, err
	}

	return toPositionResponse(updated), nil
}

func (h *handlerImpl) GetMemberPosition(ctx context.Context, req *dto.GetMemberPositionRequest) (*dto.PositionResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	pos, err := h.service.GetMemberPosition(ctx, actorID, req.PartyID, req.UserID)
	if err != nil {
		return nil, err
	}

	return toPositionResponse(pos), nil
}

func (h *handlerImpl) GetPartyPositions(ctx context.Context, req *dto.GetPartyPositionsRequest) (*dto.PositionsResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	positions, err := h.service.GetPartyPositions(ctx, actorID, req.PartyID)
	if err != nil {
		return nil, err
	}

	resp := &dto.PositionsResponse{
		Positions: make([]dto.PositionResponse, 0, len(positions)),
	}
	for _, pos := range positions {
		resp.Positions = append(resp.Positions, *toPositionResponse(pos))
	}

	return resp, nil
}

func (h *handlerImpl) UpdateTripStatus(ctx context.Context, req *dto.UpdateTripStatusRequest) (*dto.PartyMemberResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	member, err := h.service.UpdateTripStatus(ctx, actorID, req.PartyID, entity.TripStatus(req.Body.Status))
	if err != nil {
		return nil, err
	}

	return toPartyMemberResponse(member), nil
}

func toPartyMemberResponse(m entity.PartyMember) *dto.PartyMemberResponse {
	return &dto.PartyMemberResponse{
		ID:               m.ID,
		PartyID:          m.PartyID,
		UserID:           m.UserID,
		UserDisplayName:  m.UserDisplayName,
		UserProfileImage: m.UserProfileImage,
		Status:           string(m.Status),
		TripStatus:       string(m.TripStatus),
	}
}

func toTripResponse(t entity.Trip) *dto.TripResponse {
	return &dto.TripResponse{
		ID:              t.ID,
		PartyID:         t.PartyID,
		UserID:          t.UserID,
		Direction:       string(t.Direction),
		DestinationName: t.Destination.Name,
		DestinationLat:  t.Destination.Lat,
		DestinationLng:  t.Destination.Lng,
		StartedAt:       t.StartedAt,
	}
}

func toPositionResponse(p entity.Position) *dto.PositionResponse {
	return &dto.PositionResponse{
		ID:                       p.ID,
		TripID:                   p.TripID,
		PartyID:                  p.PartyID,
		UserID:                   p.UserID,
		Lat:                      p.Lat,
		Lng:                      p.Lng,
		RecordedAt:               p.RecordedAt,
		EstimatedDurationSeconds: p.EstimatedDurationSeconds,
		EstimatedArrivalAt:       p.EstimatedArrivalAt,
	}
}
