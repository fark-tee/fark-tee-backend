package party

import (
	"context"
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/fcm"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/devicetoken"
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
	MyParties(ctx context.Context, actorID string) ([]entity.Party, error)
	// Get returns the party identified by partyID, provided actorID is a
	// member of it. If actorID is not a member, a not-found error is
	// returned so as not to leak the existence of parties the actor isn't
	// in.
	Get(ctx context.Context, actorID, partyID string) (entity.Party, error)
	// ListMembers returns every member of partyID, provided actorID is a
	// member of it.
	ListMembers(ctx context.Context, actorID, partyID string) ([]entity.PartyMember, error)
	AcceptInvite(ctx context.Context, actorID, partyID string) (entity.PartyMember, error)
	DeclineInvite(ctx context.Context, actorID, partyID string) error
	// RemoveMember removes targetUserID from partyID. Only the party's
	// owner may remove members, and the owner cannot be removed.
	RemoveMember(ctx context.Context, actorID, partyID, targetUserID string) error
	// Nudge sends a best-effort push notification to targetUserID's devices
	// asking them to hurry up. Purely social - never mutates the target
	// member's real status, and a failure to deliver (missing devices, a
	// dead token, FCM being unconfigured) is not surfaced as an error to the
	// caller.
	Nudge(ctx context.Context, actorID, partyID, targetUserID string) error
}

type serviceImpl struct {
	partyRepo       party.Repository
	memberRepo      partymember.Repository
	userRepo        user.Repository
	deviceTokenRepo devicetoken.Repository
	fcmClient       fcm.Client
}

// @WireSet("Service")
func New(
	partyRepo party.Repository,
	memberRepo partymember.Repository,
	userRepo user.Repository,
	deviceTokenRepo devicetoken.Repository,
	fcmClient fcm.Client,
) Service {
	return &serviceImpl{
		partyRepo:       partyRepo,
		memberRepo:      memberRepo,
		userRepo:        userRepo,
		deviceTokenRepo: deviceTokenRepo,
		fcmClient:       fcmClient,
	}
}
