package auth

import (
	"context"
	"errors"

	"github.com/fark-tee/go-kit/apperror"
	"github.com/fark-tee/go-kit/idx"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

func (s *serviceImpl) InstagramAuthCodeURL(state string) string {
	return s.verifier.AuthCodeURL(state)
}

func (s *serviceImpl) LoginWithInstagram(ctx context.Context, code string) (entity.User, error) {
	profile, err := s.verifier.Exchange(ctx, code)
	if err != nil {
		return entity.User{}, apperror.NewUnauthorizedError("INSTAGRAM_AUTH_FAILED", "failed to authenticate with Instagram", err)
	}

	existing, err := s.userRepo.FindByInstagramUserID(ctx, profile.InstagramUserID)
	if err == nil {
		return existing, nil
	}

	if !errors.Is(err, user.ErrNotFound) {
		return entity.User{}, err
	}

	return s.userRepo.Create(ctx, entity.User{
		ID:              idx.NewUUID(),
		DisplayName:     profile.Username,
		InstagramUserID: profile.InstagramUserID,
	})
}
