package auth

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/instagramoauth"
	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/token"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

type Service interface {
	InstagramAuthCodeURL(state string) string
	// LoginWithInstagram logs the user in via Instagram, creating them if
	// needed, and returns a signed access token alongside their profile.
	LoginWithInstagram(ctx context.Context, code string) (entity.User, string, error)
}

type serviceImpl struct {
	userRepo     user.Repository
	verifier     *instagramoauth.Verifier
	tokenManager *token.Manager
}

// @WireSet("Service")
func New(userRepo user.Repository, verifier *instagramoauth.Verifier, tokenManager *token.Manager) Service {
	return &serviceImpl{
		userRepo:     userRepo,
		verifier:     verifier,
		tokenManager: tokenManager,
	}
}
