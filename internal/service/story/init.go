package story

import (
	"context"
	"io"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/storage"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/partymember"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/story"
)

type Service interface {
	// Create uploads image as a public object and records a story for
	// actorID in partyID. actorID must be an accepted member of the party.
	Create(ctx context.Context, actorID, partyID string, image io.Reader, size int64, contentType, filename string) (entity.Story, error)
	// Delete removes the story if it belongs to partyID and was created by
	// actorID.
	Delete(ctx context.Context, actorID, partyID, storyID string) error
	// ListByUser returns targetUserID's stories in partyID. actorID must be
	// an accepted member of the party.
	ListByUser(ctx context.Context, actorID, partyID, targetUserID string) ([]entity.Story, error)
}

type serviceImpl struct {
	repo       story.Repository
	memberRepo partymember.Repository
	uploader   *storage.Uploader
}

// @WireSet("Service")
func New(repo story.Repository, memberRepo partymember.Repository, uploader *storage.Uploader) Service {
	return &serviceImpl{
		repo:       repo,
		memberRepo: memberRepo,
		uploader:   uploader,
	}
}
