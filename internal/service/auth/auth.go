package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/fark-tee/go-kit/apperror"
	"github.com/fark-tee/go-kit/idx"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

func (s *serviceImpl) GoogleAuthCodeURL(redirectURI string) (string, error) {
	if !hasAllowedPrefix(redirectURI, s.cfg.Redirect.AllowedPrefixes) {
		return "", apperror.NewBadRequestError("INVALID_REDIRECT_URI", "redirect_uri is not an allowed destination")
	}

	state := base64.RawURLEncoding.EncodeToString([]byte(redirectURI))

	return s.verifier.AuthCodeURL(state), nil
}

func hasAllowedPrefix(redirectURI string, allowedPrefixes []string) bool {
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(redirectURI, prefix) {
			return true
		}
	}

	return false
}

func (s *serviceImpl) LoginWithGoogle(ctx context.Context, code, state string) (GoogleLoginResult, error) {
	redirectURI, err := decodeRedirectURI(state)
	if err != nil {
		return GoogleLoginResult{}, err
	}

	profile, err := s.verifier.Exchange(ctx, code)
	if err != nil {
		return GoogleLoginResult{}, apperror.NewUnauthorizedError("GOOGLE_AUTH_FAILED", "failed to authenticate with Google", err)
	}

	loggedInUser, err := s.userRepo.FindByGoogleUserID(ctx, profile.GoogleUserID)
	if err != nil {
		if !errors.Is(err, user.ErrNotFound) {
			return GoogleLoginResult{}, err
		}

		loggedInUser, err = s.userRepo.Create(ctx, entity.User{
			ID:              idx.NewUUID(),
			DisplayName:     profile.Name,
			ProfileImageURL: profile.ProfileImageURL,
			GoogleUserID:    profile.GoogleUserID,
		})
		if err != nil {
			return GoogleLoginResult{}, err
		}
	}

	accessToken, refreshToken, err := s.issueTokenPair(loggedInUser.ID)
	if err != nil {
		return GoogleLoginResult{}, err
	}

	return GoogleLoginResult{
		User:         loggedInUser,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		RedirectURI:  redirectURI,
	}, nil
}

func (s *serviceImpl) RefreshAccessToken(_ context.Context, refreshToken string) (string, string, error) {
	claims, err := s.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", "", apperror.NewUnauthorizedError("INVALID_REFRESH_TOKEN", "refresh token is invalid or expired", err)
	}

	return s.issueTokenPair(claims.UserID)
}

func (s *serviceImpl) issueTokenPair(userID string) (accessToken, refreshToken string, err error) {
	accessToken, err = s.tokenManager.Issue(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = s.tokenManager.IssueRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// decodeRedirectURI recovers the mobile app deeplink that was encoded into
// the OAuth state parameter by GoogleAuthCodeURL.
func decodeRedirectURI(state string) (string, error) {
	redirectURI, err := base64.RawURLEncoding.DecodeString(state)
	if err != nil {
		return "", apperror.NewBadRequestError("INVALID_STATE", "state parameter is invalid", err)
	}

	return string(redirectURI), nil
}
