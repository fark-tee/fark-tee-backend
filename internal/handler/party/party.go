package party

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/middleware/authmw"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) Create(ctx context.Context, req *dto.CreatePartyRequest) (*dto.PartyResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	created, err := h.service.Create(
		ctx,
		actorID,
		req.Body.Name,
		req.Body.DestinationName,
		req.Body.DestinationLat,
		req.Body.DestinationLng,
		req.Body.TargetTime,
	)
	if err != nil {
		return nil, err
	}

	return toPartyResponse(created), nil
}

func (h *handlerImpl) Invite(ctx context.Context, req *dto.InviteToPartyRequest) (*dto.PartyMemberResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	member, err := h.service.Invite(ctx, actorID, req.PartyID, req.Body.UserID)
	if err != nil {
		return nil, err
	}

	return toPartyMemberResponse(member), nil
}

func (h *handlerImpl) MyInvites(ctx context.Context, _ *dto.MyInvitesRequest) (*dto.PartyInvitesResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	invites, err := h.service.MyInvites(ctx, actorID)
	if err != nil {
		return nil, err
	}

	resp := &dto.PartyInvitesResponse{
		Invites: make([]dto.PartyInviteResponse, 0, len(invites)),
	}
	for _, invite := range invites {
		resp.Invites = append(resp.Invites, dto.PartyInviteResponse{
			Party:  *toPartyResponse(invite.Party),
			Member: *toPartyMemberResponse(invite.Member),
		})
	}

	return resp, nil
}

func (h *handlerImpl) MyParties(ctx context.Context, _ *dto.MyPartiesRequest) (*dto.PartiesResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	parties, err := h.service.MyParties(ctx, actorID)
	if err != nil {
		return nil, err
	}

	resp := &dto.PartiesResponse{
		Parties: make([]dto.PartyResponse, 0, len(parties)),
	}
	for _, p := range parties {
		resp.Parties = append(resp.Parties, *toPartyResponse(p))
	}

	return resp, nil
}

func (h *handlerImpl) Get(ctx context.Context, req *dto.GetPartyRequest) (*dto.PartyResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	p, err := h.service.Get(ctx, actorID, req.PartyID)
	if err != nil {
		return nil, err
	}

	return toPartyResponse(p), nil
}

func (h *handlerImpl) ListMembers(ctx context.Context, req *dto.ListPartyMembersRequest) (*dto.PartyMembersResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	members, err := h.service.ListMembers(ctx, actorID, req.PartyID)
	if err != nil {
		return nil, err
	}

	resp := &dto.PartyMembersResponse{
		Members: make([]dto.PartyMemberResponse, 0, len(members)),
	}
	for _, m := range members {
		resp.Members = append(resp.Members, *toPartyMemberResponse(m))
	}

	return resp, nil
}

func (h *handlerImpl) AcceptInvite(ctx context.Context, req *dto.AcceptInviteRequest) (*dto.PartyMemberResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	member, err := h.service.AcceptInvite(ctx, actorID, req.PartyID)
	if err != nil {
		return nil, err
	}

	return toPartyMemberResponse(member), nil
}

func (h *handlerImpl) DeclineInvite(ctx context.Context, req *dto.DeclineInviteRequest) (*dto.DeclineInviteResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.service.DeclineInvite(ctx, actorID, req.PartyID); err != nil {
		return nil, err
	}

	return &dto.DeclineInviteResponse{}, nil
}

func (h *handlerImpl) RemoveMember(ctx context.Context, req *dto.RemovePartyMemberRequest) (*dto.RemovePartyMemberResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.service.RemoveMember(ctx, actorID, req.PartyID, req.UserID); err != nil {
		return nil, err
	}

	return &dto.RemovePartyMemberResponse{}, nil
}

func (h *handlerImpl) Nudge(ctx context.Context, req *dto.NudgePartyMemberRequest) (*dto.NudgePartyMemberResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.service.Nudge(ctx, actorID, req.PartyID, req.UserID); err != nil {
		return nil, err
	}

	return &dto.NudgePartyMemberResponse{}, nil
}

func (h *handlerImpl) RequestCheckIn(ctx context.Context, req *dto.RequestCheckInRequest) (*dto.RequestCheckInResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.service.RequestCheckIn(ctx, actorID, req.PartyID, req.UserID); err != nil {
		return nil, err
	}

	return &dto.RequestCheckInResponse{}, nil
}

func (h *handlerImpl) RespondCheckIn(ctx context.Context, req *dto.RespondCheckInRequest) (*dto.PartyMemberResponse, error) {
	actorID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	member, err := h.service.RespondCheckIn(ctx, actorID, req.PartyID, entity.CheckInStatus(req.Body.Status))
	if err != nil {
		return nil, err
	}

	return toPartyMemberResponse(member), nil
}

func toPartyResponse(p entity.Party) *dto.PartyResponse {
	return &dto.PartyResponse{
		ID:              p.ID,
		Name:            p.Name,
		DestinationName: p.DestinationName,
		DestinationLat:  p.DestinationLat,
		DestinationLng:  p.DestinationLng,
		TargetTime:      p.TargetTime,
		CreatedByID:     p.CreatedByID,
		CreatedByName:   p.CreatedByName,
	}
}

func toPartyMemberResponse(m entity.PartyMember) *dto.PartyMemberResponse {
	checkInStatus := m.CheckInStatus
	if checkInStatus == "" {
		checkInStatus = entity.CheckInStatusNone
	}

	return &dto.PartyMemberResponse{
		ID:                       m.ID,
		PartyID:                  m.PartyID,
		UserID:                   m.UserID,
		UserDisplayName:          m.UserDisplayName,
		UserProfileImage:         m.UserProfileImage,
		Status:                   string(m.Status),
		TripStatus:               string(m.TripStatus),
		CheckInStatus:            string(checkInStatus),
		CheckInRequestedByUserID: m.CheckInRequestedByUserID,
	}
}
