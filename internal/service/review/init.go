package review

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/partymember"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/review"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

type Service interface {
	// CreateReview lets actorID leave a score for targetUserID within
	// partyID, once targetUserID has arrived at the destination. It updates
	// targetUserID's aggregate rating.
	CreateReview(ctx context.Context, actorID, partyID, targetUserID string, score int) (entity.Review, error)
	// ListMyReviews returns every review actorID has left within partyID.
	ListMyReviews(ctx context.Context, actorID, partyID string) ([]entity.Review, error)
}

type serviceImpl struct {
	reviewRepo review.Repository
	memberRepo partymember.Repository
	userRepo   user.Repository
}

// @WireSet("Service")
func New(reviewRepo review.Repository, memberRepo partymember.Repository, userRepo user.Repository) Service {
	return &serviceImpl{
		reviewRepo: reviewRepo,
		memberRepo: memberRepo,
		userRepo:   userRepo,
	}
}
