package story

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/story"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	Create(ctx context.Context, req *dto.CreateStoryRequest) (*dto.StoryResponse, error)
	ListByUser(ctx context.Context, req *dto.ListStoriesByUserRequest) (*dto.StoriesResponse, error)
	Delete(ctx context.Context, req *dto.DeleteStoryRequest) (*dto.DeleteStoryResponse, error)
}

type handlerImpl struct {
	service story.Service
}

// @WireSet("Handler")
func New(service story.Service) Handler {
	return &handlerImpl{service: service}
}
