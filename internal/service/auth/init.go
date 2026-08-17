package auth

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/config"
	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/googleoauth"
	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/token"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

// GoogleLoginResult is the outcome of a completed Google OAuth login: the
// (possibly newly created) user, a fresh access/refresh token pair, and the
// mobile app deeplink the caller asked to be sent back to when it started
// the flow.
type GoogleLoginResult struct {
	User         entity.User
	AccessToken  string
	RefreshToken string
	RedirectURI  string
}

type Service interface {
	// GoogleAuthCodeURL builds the URL to send the user to in order to grant
	// this app access to their Google account, encoding redirectURI (the
	// mobile app deeplink to return to once login completes) into the OAuth
	// state parameter. It returns an error if redirectURI does not use the
	// app's registered deeplink scheme.
	GoogleAuthCodeURL(redirectURI string) (string, error)
	// LoginWithGoogle logs the user in via Google, creating them if needed,
	// and returns a signed access/refresh token pair alongside their profile
	// and the deeplink to redirect back to.
	LoginWithGoogle(ctx context.Context, code, state string) (GoogleLoginResult, error)
	// RefreshAccessToken redeems a refresh token for a new access/refresh
	// token pair.
	RefreshAccessToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error)
}

type serviceImpl struct {
	cfg          *config.Config
	userRepo     user.Repository
	verifier     *googleoauth.Verifier
	tokenManager *token.Manager
}

// @WireSet("Service")
func New(cfg *config.Config, userRepo user.Repository, verifier *googleoauth.Verifier, tokenManager *token.Manager) Service {
	return &serviceImpl{
		cfg:          cfg,
		userRepo:     userRepo,
		verifier:     verifier,
		tokenManager: tokenManager,
	}
}
