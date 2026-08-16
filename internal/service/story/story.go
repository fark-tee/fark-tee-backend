package story

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/fark-tee/go-kit/apperror"
	"github.com/fark-tee/go-kit/idx"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/partymember"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/story"
)

func (s *serviceImpl) Create(ctx context.Context, actorID, partyID string, image io.Reader, size int64, contentType, filename string) (entity.Story, error) {
	if err := s.requireAcceptedMember(ctx, partyID, actorID); err != nil {
		return entity.Story{}, err
	}

	id := idx.NewUUID()

	imageURL, err := s.uploader.UploadPublic(ctx, "stories/"+id+filepath.Ext(filename), image, size, contentType)
	if err != nil {
		return entity.Story{}, err
	}

	return s.repo.Create(ctx, entity.Story{
		ID:        id,
		PartyID:   partyID,
		UserID:    actorID,
		Image:     imageURL,
		CreatedAt: time.Now(),
	})
}

func (s *serviceImpl) Delete(ctx context.Context, actorID, partyID, storyID string) error {
	existing, err := s.repo.FindByID(ctx, storyID)
	if err != nil {
		return toAppError(err)
	}

	if existing.PartyID != partyID {
		return apperror.NewNotFoundError("STORY_NOT_FOUND", "story not found")
	}

	if existing.UserID != actorID {
		return apperror.NewForbiddenError("NOT_OWNER", "story does not belong to this user")
	}

	if err := s.repo.Delete(ctx, storyID); err != nil {
		return toAppError(err)
	}

	if key, ok := s.uploader.KeyFromPublicURL(existing.Image); ok {
		if err := s.uploader.DeleteObject(ctx, key); err != nil {
			slog.Warn("failed to delete story image", slog.String("storyId", storyID), slog.Any("error", err))
		}
	}

	return nil
}

func (s *serviceImpl) ListByUser(ctx context.Context, actorID, partyID, targetUserID string) ([]entity.Story, error) {
	if err := s.requireAcceptedMember(ctx, partyID, actorID); err != nil {
		return nil, err
	}

	return s.repo.FindByPartyIDAndUserID(ctx, partyID, targetUserID)
}

func (s *serviceImpl) requireAcceptedMember(ctx context.Context, partyID, userID string) error {
	member, err := s.memberRepo.FindByPartyIDAndUserID(ctx, partyID, userID)
	if err != nil {
		if errors.Is(err, partymember.ErrNotFound) {
			return apperror.NewForbiddenError("NOT_PARTY_MEMBER", "user is not a member of this party")
		}

		return err
	}

	if member.Status != entity.PartyMemberStatusAccepted {
		return apperror.NewForbiddenError("NOT_PARTY_MEMBER", "user is not an accepted member of this party")
	}

	return nil
}

func toAppError(err error) error {
	if errors.Is(err, story.ErrNotFound) {
		return apperror.NewNotFoundError("STORY_NOT_FOUND", "story not found", err)
	}

	return err
}
