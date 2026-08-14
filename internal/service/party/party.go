package party

import (
	"context"
	"errors"
	"time"

	"github.com/fark-tee/go-kit/apperror"
	"github.com/fark-tee/go-kit/idx"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/party"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/partymember"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

func (s *serviceImpl) Create(ctx context.Context, actorID, name, destinationName string, destinationLat, destinationLng float64, targetTime time.Time) (entity.Party, error) {
	actor, err := s.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return entity.Party{}, toAppError(err)
	}

	created, err := s.partyRepo.Create(ctx, entity.Party{
		ID:              idx.NewUUID(),
		Name:            name,
		DestinationName: destinationName,
		DestinationLat:  destinationLat,
		DestinationLng:  destinationLng,
		TargetTime:      targetTime,
		CreatedByID:     actor.ID,
		CreatedByName:   actor.DisplayName,
	})
	if err != nil {
		return entity.Party{}, err
	}

	if _, err := s.memberRepo.Create(ctx, entity.PartyMember{
		ID:              idx.NewUUID(),
		PartyID:         created.ID,
		UserID:          actor.ID,
		UserDisplayName: actor.DisplayName,
		Status:          entity.PartyMemberStatusAccepted,
	}); err != nil {
		return entity.Party{}, err
	}

	return created, nil
}

func (s *serviceImpl) Invite(ctx context.Context, actorID, partyID, targetUserID string) (entity.PartyMember, error) {
	p, err := s.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return entity.PartyMember{}, toAppError(err)
	}

	if p.CreatedByID != actorID {
		return entity.PartyMember{}, apperror.NewForbiddenError("NOT_PARTY_OWNER", "only the party owner can invite members")
	}

	if targetUserID == p.CreatedByID {
		return entity.PartyMember{}, apperror.NewConflictError("ALREADY_PARTY_MEMBER", "user is already the party owner")
	}

	target, err := s.userRepo.FindByID(ctx, targetUserID)
	if err != nil {
		return entity.PartyMember{}, toAppError(err)
	}

	if _, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, targetUserID); err == nil {
		return entity.PartyMember{}, apperror.NewConflictError("ALREADY_PARTY_MEMBER", "user is already invited or a member of this party")
	} else if !errors.Is(err, partymember.ErrNotFound) {
		return entity.PartyMember{}, err
	}

	return s.memberRepo.Create(ctx, entity.PartyMember{
		ID:              idx.NewUUID(),
		PartyID:         partyID,
		UserID:          target.ID,
		UserDisplayName: target.DisplayName,
		Status:          entity.PartyMemberStatusPending,
	})
}

func (s *serviceImpl) MyInvites(ctx context.Context, actorID string) ([]Invite, error) {
	members, err := s.memberRepo.FindPendingByUserID(ctx, actorID)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return []Invite{}, nil
	}

	partyIDs := make([]string, 0, len(members))
	for _, member := range members {
		partyIDs = append(partyIDs, member.PartyID)
	}

	parties, err := s.partyRepo.FindByIDs(ctx, partyIDs)
	if err != nil {
		return nil, err
	}

	partiesByID := make(map[string]entity.Party, len(parties))
	for _, p := range parties {
		partiesByID[p.ID] = p
	}

	invites := make([]Invite, 0, len(members))
	for _, member := range members {
		p, ok := partiesByID[member.PartyID]
		if !ok {
			continue
		}

		invites = append(invites, Invite{Party: p, Member: member})
	}

	return invites, nil
}

func (s *serviceImpl) AcceptInvite(ctx context.Context, actorID, partyID string) (entity.PartyMember, error) {
	member, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, actorID)
	if err != nil {
		return entity.PartyMember{}, toAppError(err)
	}

	if member.Status != entity.PartyMemberStatusPending {
		return entity.PartyMember{}, apperror.NewConflictError("INVITE_NOT_PENDING", "invite is no longer pending")
	}

	return s.memberRepo.UpdateStatus(ctx, member.ID, entity.PartyMemberStatusAccepted)
}

func (s *serviceImpl) DeclineInvite(ctx context.Context, actorID, partyID string) error {
	member, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, actorID)
	if err != nil {
		return toAppError(err)
	}

	if member.Status != entity.PartyMemberStatusPending {
		return apperror.NewConflictError("INVITE_NOT_PENDING", "invite is no longer pending")
	}

	return s.memberRepo.Delete(ctx, member.ID)
}

func (s *serviceImpl) RemoveMember(ctx context.Context, actorID, partyID, targetUserID string) error {
	p, err := s.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return toAppError(err)
	}

	if p.CreatedByID != actorID {
		return apperror.NewForbiddenError("NOT_PARTY_OWNER", "only the party owner can remove members")
	}

	if targetUserID == p.CreatedByID {
		return apperror.NewBadRequestError("CANNOT_REMOVE_OWNER", "the party owner cannot be removed")
	}

	member, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, targetUserID)
	if err != nil {
		return toAppError(err)
	}

	return s.memberRepo.Delete(ctx, member.ID)
}

func toAppError(err error) error {
	switch {
	case errors.Is(err, party.ErrNotFound):
		return apperror.NewNotFoundError("PARTY_NOT_FOUND", "party not found", err)
	case errors.Is(err, partymember.ErrNotFound):
		return apperror.NewNotFoundError("PARTY_MEMBER_NOT_FOUND", "party member not found", err)
	case errors.Is(err, user.ErrNotFound):
		return apperror.NewNotFoundError("USER_NOT_FOUND", "user not found", err)
	default:
		return err
	}
}
