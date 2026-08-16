package story

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/middleware/authmw"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) Create(ctx context.Context, req *dto.CreateStoryRequest) (*dto.StoryResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	image := req.RawBody.Data().Image

	created, err := h.service.Create(ctx, actorID, req.PartyID, image, image.Size, image.ContentType, image.Filename)
	if err != nil {
		return nil, err
	}

	return toStoryResponse(created), nil
}

func (h *handlerImpl) ListByUser(ctx context.Context, req *dto.ListStoriesByUserRequest) (*dto.StoriesResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	stories, err := h.service.ListByUser(ctx, actorID, req.PartyID, req.UserID)
	if err != nil {
		return nil, err
	}

	resp := &dto.StoriesResponse{
		Stories: make([]dto.StoryResponse, 0, len(stories)),
	}
	for _, s := range stories {
		resp.Stories = append(resp.Stories, *toStoryResponse(s))
	}

	return resp, nil
}

func (h *handlerImpl) Delete(ctx context.Context, req *dto.DeleteStoryRequest) (*dto.DeleteStoryResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.service.Delete(ctx, actorID, req.PartyID, req.StoryID); err != nil {
		return nil, err
	}

	return &dto.DeleteStoryResponse{}, nil
}

func toStoryResponse(s entity.Story) *dto.StoryResponse {
	return &dto.StoryResponse{
		ID:        s.ID,
		PartyID:   s.PartyID,
		UserID:    s.UserID,
		Image:     s.Image,
		CreatedAt: s.CreatedAt,
	}
}
