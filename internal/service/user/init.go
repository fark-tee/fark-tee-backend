package user

import (
	"context"
	"io"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/storage"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

type Service interface {
	Search(ctx context.Context, query string) ([]entity.User, error)
	// GetByID returns the user with id, or a NOT_FOUND apperror if none exists.
	GetByID(ctx context.Context, id string) (entity.User, error)
	// UpdateProfile sets id's display name and username and returns the
	// updated user. Returns a NOT_FOUND apperror if no such user exists, or a
	// CONFLICT apperror if username is already taken by a different user.
	UpdateProfile(ctx context.Context, id, displayName, username string) (entity.User, error)
	// UploadProfileImage uploads image as id's new profile picture and
	// returns the updated user.
	UploadProfileImage(ctx context.Context, id string, image io.Reader, size int64, contentType, filename string) (entity.User, error)
}

type serviceImpl struct {
	repo     user.Repository
	uploader *storage.Uploader
}

// @WireSet("Service")
func New(repo user.Repository, uploader *storage.Uploader) Service {
	return &serviceImpl{repo: repo, uploader: uploader}
}
