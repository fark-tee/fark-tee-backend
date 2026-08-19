package review

import (
	"context"
	"errors"
	"time"

	"github.com/fark-tee/go-kit/apperror"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/partymember"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/review"
)

func (s *serviceImpl) CreateReview(ctx context.Context, actorID, partyID, targetUserID string, score int) (entity.Review, error) {
	if err := s.requireMember(ctx, partyID, actorID); err != nil {
		return entity.Review{}, err
	}

	if actorID == targetUserID {
		return entity.Review{}, apperror.NewBadRequestError("CANNOT_REVIEW_SELF", "you cannot review yourself")
	}

	target, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, targetUserID)
	if err != nil {
		return entity.Review{}, toAppError(err)
	}

	arrivedOrder, _ := entity.TripStatusOrder(entity.TripStatusArrived)
	targetOrder, ok := entity.TripStatusOrder(target.TripStatus)
	if !ok || targetOrder < arrivedOrder {
		return entity.Review{}, apperror.NewBadRequestError("MEMBER_NOT_ARRIVED", "this member hasn't arrived at the destination yet")
	}

	if _, err := s.reviewRepo.FindByPartyIDReviewerIDAndTargetUserID(ctx, partyID, actorID, targetUserID); err == nil {
		return entity.Review{}, apperror.NewConflictError("ALREADY_REVIEWED", "you have already reviewed this member for this party")
	} else if !errors.Is(err, review.ErrNotFound) {
		return entity.Review{}, err
	}

	created, err := s.reviewRepo.Create(ctx, entity.Review{
		ID:           mongoid.New(),
		PartyID:      partyID,
		ReviewerID:   actorID,
		TargetUserID: targetUserID,
		Score:        score,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return entity.Review{}, err
	}

	if _, err := s.userRepo.RecordRating(ctx, targetUserID, score); err != nil {
		return entity.Review{}, err
	}

	return created, nil
}

func (s *serviceImpl) ListMyReviews(ctx context.Context, actorID, partyID string) ([]entity.Review, error) {
	if err := s.requireMember(ctx, partyID, actorID); err != nil {
		return nil, err
	}

	return s.reviewRepo.FindByPartyIDAndReviewerID(ctx, partyID, actorID)
}

// requireMember returns a forbidden error unless userID is an accepted
// member of partyID.
func (s *serviceImpl) requireMember(ctx context.Context, partyID, userID string) error {
	member, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, userID)
	if err != nil {
		if errors.Is(err, partymember.ErrNotFound) {
			return apperror.NewForbiddenError("NOT_PARTY_MEMBER", "you are not a member of this party")
		}

		return toAppError(err)
	}

	if member.Status != entity.PartyMemberStatusAccepted {
		return apperror.NewForbiddenError("NOT_PARTY_MEMBER", "you are not an accepted member of this party")
	}

	return nil
}

func toAppError(err error) error {
	switch {
	case errors.Is(err, partymember.ErrNotFound):
		return apperror.NewNotFoundError("PARTY_MEMBER_NOT_FOUND", "party member not found", err)
	case errors.Is(err, review.ErrNotFound):
		return apperror.NewNotFoundError("REVIEW_NOT_FOUND", "review not found", err)
	default:
		return err
	}
}
