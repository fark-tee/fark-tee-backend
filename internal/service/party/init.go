package party

import (
	"context"
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/party"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/partymember"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

// Invite pairs a pending party member row with the party it belongs to, for
// callers that need both to render an invite (e.g. "my invite list").
type Invite struct {
	Party  entity.Party
	Member entity.PartyMember
}

type Service interface {
	Create(ctx context.Context, actorID, name, destinationName string, destinationLat, destinationLng float64, targetTime time.Time) (entity.Party, error)
	// Invite adds targetUserID to partyID as a pending member. Only the
	// party's owner may invite.
	Invite(ctx context.Context, actorID, partyID, targetUserID string) (entity.PartyMember, error)
	MyInvites(ctx context.Context, actorID string) ([]Invite, error)
	AcceptInvite(ctx context.Context, actorID, partyID string) (entity.PartyMember, error)
	DeclineInvite(ctx context.Context, actorID, partyID string) error
	// RemoveMember removes targetUserID from partyID. Only the party's
	// owner may remove members, and the owner cannot be removed.
	RemoveMember(ctx context.Context, actorID, partyID, targetUserID string) error
}

type serviceImpl struct {
	partyRepo  party.Repository
	memberRepo partymember.Repository
	userRepo   user.Repository
}

// @WireSet("Service")
func New(partyRepo party.Repository, memberRepo partymember.Repository, userRepo user.Repository) Service {
	return &serviceImpl{
		partyRepo:  partyRepo,
		memberRepo: memberRepo,
		userRepo:   userRepo,
	}
}
