package savedlocation

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/middleware/authmw"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) Create(ctx context.Context, req *dto.CreateSavedLocationRequest) (*dto.SavedLocationResponse, error) {
	userID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	location, err := h.service.Create(ctx, userID, req.Body.Name, req.Body.Lat, req.Body.Lng)
	if err != nil {
		return nil, err
	}

	return toSavedLocationResponse(location), nil
}

func (h *handlerImpl) Get(ctx context.Context, req *dto.GetSavedLocationRequest) (*dto.SavedLocationResponse, error) {
	userID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	location, err := h.service.Get(ctx, userID, req.ID)
	if err != nil {
		return nil, err
	}

	return toSavedLocationResponse(location), nil
}

func (h *handlerImpl) List(ctx context.Context, _ *dto.ListSavedLocationsRequest) (*dto.SavedLocationsResponse, error) {
	userID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	locations, err := h.service.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &dto.SavedLocationsResponse{
		SavedLocations: make([]dto.SavedLocationResponse, 0, len(locations)),
	}
	for _, location := range locations {
		resp.SavedLocations = append(resp.SavedLocations, *toSavedLocationResponse(location))
	}

	return resp, nil
}

func (h *handlerImpl) Update(ctx context.Context, req *dto.UpdateSavedLocationRequest) (*dto.SavedLocationResponse, error) {
	userID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	location, err := h.service.Update(ctx, userID, req.ID, req.Body.Name, req.Body.Lat, req.Body.Lng)
	if err != nil {
		return nil, err
	}

	return toSavedLocationResponse(location), nil
}

func (h *handlerImpl) Delete(ctx context.Context, req *dto.DeleteSavedLocationRequest) (*dto.DeleteSavedLocationResponse, error) {
	userID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.service.Delete(ctx, userID, req.ID); err != nil {
		return nil, err
	}

	return &dto.DeleteSavedLocationResponse{}, nil
}

func toSavedLocationResponse(l entity.SavedLocation) *dto.SavedLocationResponse {
	return &dto.SavedLocationResponse{
		ID:     l.ID,
		UserID: l.UserID,
		Name:   l.Name,
		Lat:    l.Lat,
		Lng:    l.Lng,
	}
}
