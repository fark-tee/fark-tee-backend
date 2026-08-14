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
	AcceptInvite(ctx context.Context, req *dto.AcceptInviteRequest) (*dto.PartyMemberResponse, error)
	DeclineInvite(ctx context.Context, req *dto.DeclineInviteRequest) (*dto.DeclineInviteResponse, error)
	RemoveMember(ctx context.Context, req *dto.RemovePartyMemberRequest) (*dto.RemovePartyMemberResponse, error)
}

type handlerImpl struct {
	service party.Service
}

// @WireSet("Handler")
func New(service party.Service) Handler {
	return &handlerImpl{service: service}
}
