package trip

import (
	"context"
	"errors"
	"time"

	"github.com/fark-tee/go-kit/apperror"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/party"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/partymember"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/position"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/trip"
)

func (s *serviceImpl) StartTrip(ctx context.Context, actorID, partyID string, direction entity.TripDirection, lat, lng float64) (entity.Trip, entity.Position, error) {
	if _, err := s.partyRepo.FindByID(ctx, partyID); err != nil {
		return entity.Trip{}, entity.Position{}, toAppError(err)
	}

	if err := s.requireMember(ctx, partyID, actorID); err != nil {
		return entity.Trip{}, entity.Position{}, err
	}

	createdTrip, err := s.tripRepo.Create(ctx, entity.Trip{
		ID:        mongoid.New(),
		PartyID:   partyID,
		UserID:    actorID,
		Direction: direction,
		StartedAt: time.Now(),
	})
	if err != nil {
		return entity.Trip{}, entity.Position{}, err
	}

	createdPosition, err := s.positionRepo.Create(ctx, entity.Position{
		ID:         mongoid.New(),
		TripID:     createdTrip.ID,
		PartyID:    partyID,
		UserID:     actorID,
		Lat:        lat,
		Lng:        lng,
		RecordedAt: createdTrip.StartedAt,
	})
	if err != nil {
		return entity.Trip{}, entity.Position{}, err
	}

	return createdTrip, createdPosition, nil
}

func (s *serviceImpl) UpdatePosition(ctx context.Context, actorID, partyID string, lat, lng float64) (entity.Position, error) {
	if err := s.requireMember(ctx, partyID, actorID); err != nil {
		return entity.Position{}, err
	}

	currentTrip, err := s.tripRepo.FindLatestByPartyIDAndUserID(ctx, partyID, actorID)
	if err != nil {
		if errors.Is(err, trip.ErrNotFound) {
			return entity.Position{}, apperror.NewBadRequestError("NO_ACTIVE_TRIP", "you must start a trip before updating your position")
		}

		return entity.Position{}, toAppError(err)
	}

	return s.positionRepo.Create(ctx, entity.Position{
		ID:         mongoid.New(),
		TripID:     currentTrip.ID,
		PartyID:    partyID,
		UserID:     actorID,
		Lat:        lat,
		Lng:        lng,
		RecordedAt: time.Now(),
	})
}

func (s *serviceImpl) GetMemberPosition(ctx context.Context, actorID, partyID, targetUserID string) (entity.Position, error) {
	if err := s.requireMember(ctx, partyID, actorID); err != nil {
		return entity.Position{}, err
	}

	pos, err := s.positionRepo.FindLatestByPartyIDAndUserID(ctx, partyID, targetUserID)
	if err != nil {
		return entity.Position{}, toAppError(err)
	}

	return pos, nil
}

func (s *serviceImpl) GetPartyPositions(ctx context.Context, actorID, partyID string) ([]entity.Position, error) {
	if err := s.requireMember(ctx, partyID, actorID); err != nil {
		return nil, err
	}

	return s.positionRepo.FindLatestByPartyID(ctx, partyID)
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
	case errors.Is(err, party.ErrNotFound):
		return apperror.NewNotFoundError("PARTY_NOT_FOUND", "party not found", err)
	case errors.Is(err, partymember.ErrNotFound):
		return apperror.NewNotFoundError("PARTY_MEMBER_NOT_FOUND", "party member not found", err)
	case errors.Is(err, trip.ErrNotFound):
		return apperror.NewNotFoundError("TRIP_NOT_FOUND", "trip not found", err)
	case errors.Is(err, position.ErrNotFound):
		return apperror.NewNotFoundError("POSITION_NOT_FOUND", "position not found", err)
	default:
		return err
	}
}
