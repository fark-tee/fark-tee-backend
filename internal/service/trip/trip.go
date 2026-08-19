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

func (s *serviceImpl) StartTrip(ctx context.Context, actorID, partyID string, direction entity.TripDirection, lat, lng float64, destination entity.Destination) (entity.Trip, entity.Position, error) {
	if _, err := s.partyRepo.FindByID(ctx, partyID); err != nil {
		return entity.Trip{}, entity.Position{}, toAppError(err)
	}

	if err := s.requireMember(ctx, partyID, actorID); err != nil {
		return entity.Trip{}, entity.Position{}, err
	}

	createdTrip, err := s.tripRepo.Create(ctx, entity.Trip{
		ID:          mongoid.New(),
		PartyID:     partyID,
		UserID:      actorID,
		Direction:   direction,
		Destination: destination,
		StartedAt:   time.Now(),
	})
	if err != nil {
		return entity.Trip{}, entity.Position{}, err
	}

	duration, arrivalAt, err := s.estimateArrival(ctx, lat, lng, destination, createdTrip.StartedAt)
	if err != nil {
		return entity.Trip{}, entity.Position{}, err
	}

	createdPosition, err := s.positionRepo.Create(ctx, entity.Position{
		ID:                       mongoid.New(),
		TripID:                   createdTrip.ID,
		PartyID:                  partyID,
		UserID:                   actorID,
		Lat:                      lat,
		Lng:                      lng,
		RecordedAt:               createdTrip.StartedAt,
		EstimatedDurationSeconds: duration,
		EstimatedArrivalAt:       arrivalAt,
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

	recordedAt := time.Now()

	duration, arrivalAt, err := s.estimateArrival(ctx, lat, lng, currentTrip.Destination, recordedAt)
	if err != nil {
		return entity.Position{}, err
	}

	return s.positionRepo.Create(ctx, entity.Position{
		ID:                       mongoid.New(),
		TripID:                   currentTrip.ID,
		PartyID:                  partyID,
		UserID:                   actorID,
		Lat:                      lat,
		Lng:                      lng,
		RecordedAt:               recordedAt,
		EstimatedDurationSeconds: duration,
		EstimatedArrivalAt:       arrivalAt,
	})
}

// estimateArrival calls OSRM to estimate the travel time from (lat, lng) to
// destination, returning the duration and the resulting arrival time
// relative to from.
func (s *serviceImpl) estimateArrival(ctx context.Context, lat, lng float64, destination entity.Destination, from time.Time) (int, time.Time, error) {
	route, err := s.osrmClient.Route(ctx, lat, lng, destination.Lat, destination.Lng)
	if err != nil {
		return 0, time.Time{}, apperror.NewServiceUnavailableError("ROUTE_ESTIMATION_FAILED", "could not estimate travel time", err)
	}

	return route.DurationSeconds, from.Add(time.Duration(route.DurationSeconds) * time.Second), nil
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

// requireAcceptedMember returns the caller's own party member row, or a
// forbidden error unless userID is an accepted member of partyID.
func (s *serviceImpl) requireAcceptedMember(ctx context.Context, partyID, userID string) (entity.PartyMember, error) {
	member, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, userID)
	if err != nil {
		if errors.Is(err, partymember.ErrNotFound) {
			return entity.PartyMember{}, apperror.NewForbiddenError("NOT_PARTY_MEMBER", "you are not a member of this party")
		}

		return entity.PartyMember{}, toAppError(err)
	}

	if member.Status != entity.PartyMemberStatusAccepted {
		return entity.PartyMember{}, apperror.NewForbiddenError("NOT_PARTY_MEMBER", "you are not an accepted member of this party")
	}

	return member, nil
}

// requireMember returns a forbidden error unless userID is an accepted
// member of partyID.
func (s *serviceImpl) requireMember(ctx context.Context, partyID, userID string) error {
	_, err := s.requireAcceptedMember(ctx, partyID, userID)
	return err
}

// UpdateTripStatus advances actorID's trip status within partyID. Trip
// status only ever moves forward (e.g. PENDING_DEPARTURE -> ARRIVED is
// allowed, skipping DEPARTED, but ARRIVED -> DEPARTED is not).
func (s *serviceImpl) UpdateTripStatus(ctx context.Context, actorID, partyID string, status entity.TripStatus) (entity.PartyMember, error) {
	newOrder, ok := entity.TripStatusOrder(status)
	if !ok {
		return entity.PartyMember{}, apperror.NewBadRequestError("INVALID_TRIP_STATUS", "invalid trip status")
	}

	member, err := s.requireAcceptedMember(ctx, partyID, actorID)
	if err != nil {
		return entity.PartyMember{}, err
	}

	currentOrder, _ := entity.TripStatusOrder(member.TripStatus)
	if newOrder <= currentOrder {
		return entity.PartyMember{}, apperror.NewConflictError("INVALID_TRIP_STATUS_TRANSITION", "trip status can only move forward")
	}

	return s.memberRepo.UpdateTripStatus(ctx, member.ID, status)
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
