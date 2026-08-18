package user

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"path/filepath"

	"github.com/fark-tee/go-kit/apperror"
	"github.com/fark-tee/go-kit/idx"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

func (s *serviceImpl) Search(ctx context.Context, query string) ([]entity.User, error) {
	return s.repo.Search(ctx, query)
}

func (s *serviceImpl) GetByID(ctx context.Context, id string) (entity.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, toAppError(err)
	}

	return u, nil
}

func (s *serviceImpl) UpdateProfile(ctx context.Context, id, displayName, username string) (entity.User, error) {
	existing, err := s.repo.FindByUsername(ctx, username)
	if err == nil && existing.ID != id {
		return entity.User{}, apperror.NewConflictError("USERNAME_TAKEN", "username is already taken")
	} else if err != nil && !errors.Is(err, user.ErrNotFound) {
		return entity.User{}, err
	}

	u, err := s.repo.UpdateProfile(ctx, id, displayName, username)
	if err != nil {
		return entity.User{}, toAppError(err)
	}

	return u, nil
}

func (s *serviceImpl) UploadProfileImage(ctx context.Context, id string, image io.Reader, size int64, contentType, filename string) (entity.User, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return entity.User{}, toAppError(err)
	}

	key := "profile-images/" + idx.NewUUID() + filepath.Ext(filename)

	imageURL, err := s.uploader.UploadPublic(ctx, key, image, size, contentType)
	if err != nil {
		return entity.User{}, err
	}

	u, err := s.repo.UpdateProfileImage(ctx, id, imageURL)
	if err != nil {
		return entity.User{}, toAppError(err)
	}

	if oldKey, ok := s.uploader.KeyFromPublicURL(existing.ProfileImageURL); ok {
		if err := s.uploader.DeleteObject(ctx, oldKey); err != nil {
			slog.Warn("failed to delete old profile image", slog.String("userId", id), slog.Any("error", err))
		}
	}

	return u, nil
}

func toAppError(err error) error {
	switch {
	case errors.Is(err, user.ErrNotFound):
		return apperror.NewNotFoundError("USER_NOT_FOUND", "user not found", err)
	default:
		return err
	}
}
