package savedlocation

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) Create(ctx context.Context, req *dto.CreateSavedLocationRequest) (*dto.SavedLocationResponse, error) {
	location, err := h.service.Create(ctx, req.UserID, req.Body.Name, req.Body.Lat, req.Body.Lng)
	if err != nil {
		return nil, err
	}

	return toSavedLocationResponse(location), nil
}

func (h *handlerImpl) Get(ctx context.Context, req *dto.GetSavedLocationRequest) (*dto.SavedLocationResponse, error) {
	location, err := h.service.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return toSavedLocationResponse(location), nil
}

func (h *handlerImpl) List(ctx context.Context, req *dto.ListSavedLocationsRequest) (*dto.SavedLocationsResponse, error) {
	locations, err := h.service.ListByUserID(ctx, req.UserID)
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
	location, err := h.service.Update(ctx, req.ID, req.Body.Name, req.Body.Lat, req.Body.Lng)
	if err != nil {
		return nil, err
	}

	return toSavedLocationResponse(location), nil
}

func (h *handlerImpl) Delete(ctx context.Context, req *dto.DeleteSavedLocationRequest) (*dto.DeleteSavedLocationResponse, error) {
	if err := h.service.Delete(ctx, req.ID); err != nil {
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
