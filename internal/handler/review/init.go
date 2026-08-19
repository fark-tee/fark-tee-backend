package review

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/review"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	CreateReview(ctx context.Context, req *dto.CreateReviewRequest) (*dto.ReviewResponse, error)
	ListMyReviews(ctx context.Context, req *dto.ListMyReviewsRequest) (*dto.ReviewsResponse, error)
}

type handlerImpl struct {
	service review.Service
}

// @WireSet("Handler")
func New(service review.Service) Handler {
	return &handlerImpl{service: service}
}
