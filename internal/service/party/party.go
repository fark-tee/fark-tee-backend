package party

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/fark-tee/go-kit/apperror"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/party"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/partymember"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

func (s *serviceImpl) Create(ctx context.Context, actorID, name, destinationName string, destinationLat, destinationLng float64, targetTime time.Time, note string) (entity.Party, error) {
	actor, err := s.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return entity.Party{}, toAppError(err)
	}

	created, err := s.partyRepo.Create(ctx, entity.Party{
		ID:              mongoid.New(),
		Name:            name,
		DestinationName: destinationName,
		DestinationLat:  destinationLat,
		DestinationLng:  destinationLng,
		TargetTime:      targetTime,
		CreatedByID:     actor.ID,
		CreatedByName:   actor.DisplayName,
		Note:            note,
	})
	if err != nil {
		return entity.Party{}, err
	}

	if _, err := s.memberRepo.Create(ctx, entity.PartyMember{
		ID:               mongoid.New(),
		PartyID:          created.ID,
		UserID:           actor.ID,
		UserDisplayName:  actor.DisplayName,
		UserProfileImage: actor.ProfileImageURL,
		Status:           entity.PartyMemberStatusAccepted,
		TripStatus:       entity.TripStatusPendingDeparture,
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

	created, err := s.memberRepo.Create(ctx, entity.PartyMember{
		ID:               mongoid.New(),
		PartyID:          partyID,
		UserID:           target.ID,
		UserDisplayName:  target.DisplayName,
		UserProfileImage: target.ProfileImageURL,
		Status:           entity.PartyMemberStatusPending,
		TripStatus:       entity.TripStatusPendingDeparture,
	})
	if err != nil {
		return entity.PartyMember{}, err
	}

	s.notifyPartyInvite(ctx, p, target.ID)

	return created, nil
}

// notifyPartyInvite best-effort pushes a "you've been invited" notification
// to every device targetUserID has registered. Matches Nudge/RequestCheckIn:
// missing devices, a dead token, or FCM being unconfigured are logged and
// swallowed rather than surfaced as an error, since the invite itself is
// already persisted regardless of whether the push lands.
func (s *serviceImpl) notifyPartyInvite(ctx context.Context, p entity.Party, targetUserID string) {
	tokens, err := s.deviceTokenRepo.FindByUserID(ctx, targetUserID)
	if err != nil {
		slog.Warn("failed to load device tokens for party invite notification",
			slog.String("targetUserId", targetUserID), slog.String("partyId", p.ID), slog.Any("error", err))

		return
	}

	for _, t := range tokens {
		if err := s.fcmClient.SendPartyInvite(ctx, t.Token, p.ID, p.Name, p.CreatedByID, p.CreatedByName); err != nil {
			slog.Warn("failed to send party invite notification",
				slog.String("targetUserId", targetUserID),
				slog.String("partyId", p.ID),
				slog.Any("error", err))
		}
	}
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

func (s *serviceImpl) MyParties(ctx context.Context, actorID string) ([]entity.Party, error) {
	members, err := s.memberRepo.FindAcceptedByUserID(ctx, actorID)
	if err != nil {
		return nil, err
	}

	if len(members) == 0 {
		return []entity.Party{}, nil
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

	result := make([]entity.Party, 0, len(members))
	for _, member := range members {
		p, ok := partiesByID[member.PartyID]
		if !ok {
			continue
		}

		result = append(result, p)
	}

	return result, nil
}

// requireMembership verifies that actorID is a member of partyID, returning
// a PARTY_NOT_FOUND error (rather than leaking whether the party exists) if
// not.
func (s *serviceImpl) requireMembership(ctx context.Context, actorID, partyID string) error {
	if _, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, actorID); err != nil {
		if errors.Is(err, partymember.ErrNotFound) {
			return apperror.NewNotFoundError("PARTY_NOT_FOUND", "party not found", err)
		}

		return err
	}

	return nil
}

func (s *serviceImpl) Get(ctx context.Context, actorID, partyID string) (entity.Party, error) {
	p, err := s.partyRepo.FindByID(ctx, partyID)
	if err != nil {
		return entity.Party{}, toAppError(err)
	}

	if err := s.requireMembership(ctx, actorID, partyID); err != nil {
		return entity.Party{}, err
	}

	return p, nil
}

func (s *serviceImpl) ListMembers(ctx context.Context, actorID, partyID string) ([]entity.PartyMember, error) {
	if err := s.requireMembership(ctx, actorID, partyID); err != nil {
		return nil, err
	}

	return s.memberRepo.FindByPartyID(ctx, partyID)
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

func (s *serviceImpl) Nudge(ctx context.Context, actorID, partyID, targetUserID string) error {
	if targetUserID == actorID {
		return apperror.NewBadRequestError("CANNOT_NUDGE_SELF", "cannot nudge yourself")
	}

	if _, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, actorID); err != nil {
		if errors.Is(err, partymember.ErrNotFound) {
			return apperror.NewForbiddenError("NOT_PARTY_MEMBER", "you are not a member of this party")
		}

		return err
	}

	if _, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, targetUserID); err != nil {
		return toAppError(err)
	}

	actor, err := s.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return toAppError(err)
	}

	tokens, err := s.deviceTokenRepo.FindByUserID(ctx, targetUserID)
	if err != nil {
		return err
	}

	if len(tokens) == 0 {
		// No devices to nudge - purely social, so this is a no-op rather
		// than an error.
		return nil
	}

	for _, t := range tokens {
		if err := s.fcmClient.SendNudge(ctx, t.Token, partyID, actor.ID, actor.DisplayName); err != nil {
			slog.Warn("failed to send nudge notification",
				slog.String("targetUserId", targetUserID),
				slog.String("partyId", partyID),
				slog.Any("error", err))
		}
	}

	return nil
}

func (s *serviceImpl) RequestCheckIn(ctx context.Context, actorID, partyID, targetUserID string) error {
	if targetUserID == actorID {
		return apperror.NewBadRequestError("CANNOT_CHECK_IN_SELF", "cannot request a check-in from yourself")
	}

	if _, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, actorID); err != nil {
		if errors.Is(err, partymember.ErrNotFound) {
			return apperror.NewForbiddenError("NOT_PARTY_MEMBER", "you are not a member of this party")
		}

		return err
	}

	target, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, targetUserID)
	if err != nil {
		return toAppError(err)
	}

	if target.TripStatus != entity.TripStatusReturning {
		return apperror.NewConflictError("MEMBER_NOT_HEADING_HOME", "this member is not currently heading home")
	}

	actor, err := s.userRepo.FindByID(ctx, actorID)
	if err != nil {
		return toAppError(err)
	}

	if _, err := s.memberRepo.UpdateCheckIn(ctx, target.ID, entity.CheckInStatusPending, actorID); err != nil {
		return err
	}

	tokens, err := s.deviceTokenRepo.FindByUserID(ctx, targetUserID)
	if err != nil {
		return err
	}

	for _, t := range tokens {
		if err := s.fcmClient.SendCheckInRequest(ctx, t.Token, partyID, actor.ID, actor.DisplayName); err != nil {
			slog.Warn("failed to send check-in request notification",
				slog.String("targetUserId", targetUserID),
				slog.String("partyId", partyID),
				slog.Any("error", err))
		}
	}

	return nil
}

func (s *serviceImpl) RespondCheckIn(ctx context.Context, actorID, partyID string, status entity.CheckInStatus) (entity.PartyMember, error) {
	if status != entity.CheckInStatusOK && status != entity.CheckInStatusNotOK {
		return entity.PartyMember{}, apperror.NewBadRequestError("INVALID_CHECK_IN_STATUS", "status must be OK or NOT_OK")
	}

	member, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, actorID)
	if err != nil {
		return entity.PartyMember{}, toAppError(err)
	}

	if member.CheckInStatus != entity.CheckInStatusPending {
		return entity.PartyMember{}, apperror.NewConflictError("NO_PENDING_CHECK_IN", "there is no pending check-in to respond to")
	}

	updated, err := s.memberRepo.UpdateCheckIn(ctx, member.ID, status, member.CheckInRequestedByUserID)
	if err != nil {
		return entity.PartyMember{}, err
	}

	if status == entity.CheckInStatusNotOK {
		s.broadcastCheckInEmergency(ctx, partyID, actorID)
	}

	return updated, nil
}

// broadcastCheckInEmergency best-effort pushes a full-screen emergency alert
// - carrying actorID's emergency contact - to every other accepted member of
// partyID. Failures (missing devices, FCM unconfigured, a lookup error) are
// logged and swallowed, matching Nudge/RequestCheckIn: this is a safety-net
// notification layered on top of the already-persisted NOT_OK status, not
// something the caller's response should fail over.
func (s *serviceImpl) broadcastCheckInEmergency(ctx context.Context, partyID, actorID string) {
	actor, err := s.userRepo.FindByID(ctx, actorID)
	if err != nil {
		slog.Warn("failed to load actor for check-in emergency broadcast",
			slog.String("actorId", actorID), slog.Any("error", err))

		return
	}

	members, err := s.memberRepo.FindByPartyID(ctx, partyID)
	if err != nil {
		slog.Warn("failed to list party members for check-in emergency broadcast",
			slog.String("partyId", partyID), slog.Any("error", err))

		return
	}

	for _, m := range members {
		if m.UserID == actorID || m.Status != entity.PartyMemberStatusAccepted {
			continue
		}

		tokens, err := s.deviceTokenRepo.FindByUserID(ctx, m.UserID)
		if err != nil {
			slog.Warn("failed to load device tokens for check-in emergency broadcast",
				slog.String("targetUserId", m.UserID), slog.Any("error", err))

			continue
		}

		for _, t := range tokens {
			if err := s.fcmClient.SendCheckInEmergencyAlert(
				ctx, t.Token, partyID, actor.ID, actor.DisplayName,
				actor.EmergencyContactName, actor.EmergencyContactPhone,
			); err != nil {
				slog.Warn("failed to send check-in emergency alert",
					slog.String("targetUserId", m.UserID),
					slog.String("partyId", partyID),
					slog.Any("error", err))
			}
		}
	}
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
