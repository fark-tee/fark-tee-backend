package review

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/middleware/authmw"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) CreateReview(ctx context.Context, req *dto.CreateReviewRequest) (*dto.ReviewResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	created, err := h.service.CreateReview(ctx, actorID, req.PartyID, req.UserID, req.Body.Score)
	if err != nil {
		return nil, err
	}

	return toReviewResponse(created), nil
}

func (h *handlerImpl) ListMyReviews(ctx context.Context, req *dto.ListMyReviewsRequest) (*dto.ReviewsResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	reviews, err := h.service.ListMyReviews(ctx, actorID, req.PartyID)
	if err != nil {
		return nil, err
	}

	resp := &dto.ReviewsResponse{
		Reviews: make([]dto.ReviewResponse, 0, len(reviews)),
	}
	for _, r := range reviews {
		resp.Reviews = append(resp.Reviews, *toReviewResponse(r))
	}

	return resp, nil
}

func toReviewResponse(r entity.Review) *dto.ReviewResponse {
	return &dto.ReviewResponse{
		ID:           r.ID,
		PartyID:      r.PartyID,
		ReviewerID:   r.ReviewerID,
		TargetUserID: r.TargetUserID,
		Score:        r.Score,
		CreatedAt:    r.CreatedAt,
	}
}
