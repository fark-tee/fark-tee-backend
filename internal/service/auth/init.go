package auth

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/instagramoauth"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

type Service interface {
	InstagramAuthCodeURL(state string) string
	LoginWithInstagram(ctx context.Context, code string) (entity.User, error)
}

type serviceImpl struct {
	userRepo user.Repository
	verifier *instagramoauth.Verifier
}

// @WireSet("Service")
func New(userRepo user.Repository, verifier *instagramoauth.Verifier) Service {
	return &serviceImpl{
		userRepo: userRepo,
		verifier: verifier,
	}
}
