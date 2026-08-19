package party

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/party"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	Create(ctx context.Context, req *dto.CreatePartyRequest) (*dto.PartyResponse, error)
	Invite(ctx context.Context, req *dto.InviteToPartyRequest) (*dto.PartyMemberResponse, error)
	MyInvites(ctx context.Context, req *dto.MyInvitesRequest) (*dto.PartyInvitesResponse, error)
	MyParties(ctx context.Context, req *dto.MyPartiesRequest) (*dto.PartiesResponse, error)
	Get(ctx context.Context, req *dto.GetPartyRequest) (*dto.PartyResponse, error)
	ListMembers(ctx context.Context, req *dto.ListPartyMembersRequest) (*dto.PartyMembersResponse, error)
	AcceptInvite(ctx context.Context, req *dto.AcceptInviteRequest) (*dto.PartyMemberResponse, error)
	DeclineInvite(ctx context.Context, req *dto.DeclineInviteRequest) (*dto.DeclineInviteResponse, error)
	RemoveMember(ctx context.Context, req *dto.RemovePartyMemberRequest) (*dto.RemovePartyMemberResponse, error)
	Nudge(ctx context.Context, req *dto.NudgePartyMemberRequest) (*dto.NudgePartyMemberResponse, error)
	RequestCheckIn(ctx context.Context, req *dto.RequestCheckInRequest) (*dto.RequestCheckInResponse, error)
	RespondCheckIn(ctx context.Context, req *dto.RespondCheckInRequest) (*dto.PartyMemberResponse, error)
}

type handlerImpl struct {
	service party.Service
}

// @WireSet("Handler")
func New(service party.Service) Handler {
	return &handlerImpl{service: service}
}
