package trip

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/osrm"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/party"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/partymember"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/position"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/trip"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

type Service interface {
	// StartTrip starts a new depart or return trip for actorID within
	// partyID, heading to destination, and records its first position along
	// with an OSRM-computed estimated arrival time.
	StartTrip(ctx context.Context, actorID, partyID string, direction entity.TripDirection, lat, lng float64, destination entity.Destination) (entity.Trip, entity.Position, error)
	// UpdatePosition records a new position for actorID's current trip
	// within partyID.
	UpdatePosition(ctx context.Context, actorID, partyID string, lat, lng float64) (entity.Position, error)
	// GetMemberPosition returns the latest recorded position of
	// targetUserID within partyID.
	GetMemberPosition(ctx context.Context, actorID, partyID, targetUserID string) (entity.Position, error)
	// GetPartyPositions returns the latest recorded position of every
	// party member that has recorded one.
	GetPartyPositions(ctx context.Context, actorID, partyID string) ([]entity.Position, error)
	// UpdateTripStatus advances actorID's trip status within partyID.
	// Trip status only ever moves forward.
	UpdateTripStatus(ctx context.Context, actorID, partyID string, status entity.TripStatus) (entity.PartyMember, error)
	// GetMemberTrip returns targetUserID's latest trip (depart or return)
	// within partyID, including its OSRM-computed route polyline.
	GetMemberTrip(ctx context.Context, actorID, partyID, targetUserID string) (entity.Trip, error)
}

type serviceImpl struct {
	partyRepo    party.Repository
	memberRepo   partymember.Repository
	tripRepo     trip.Repository
	positionRepo position.Repository
	userRepo     user.Repository
	osrmClient   *osrm.Client
}

// @WireSet("Service")
func New(partyRepo party.Repository, memberRepo partymember.Repository, tripRepo trip.Repository, positionRepo position.Repository, userRepo user.Repository, osrmClient *osrm.Client) Service {
	return &serviceImpl{
		partyRepo:    partyRepo,
		memberRepo:   memberRepo,
		tripRepo:     tripRepo,
		positionRepo: positionRepo,
		userRepo:     userRepo,
		osrmClient:   osrmClient,
	}
}
