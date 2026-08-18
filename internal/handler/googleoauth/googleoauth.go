package googleoauth

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) Start(_ context.Context, req *dto.GoogleStartRequest) (*dto.RedirectResponse, error) {
	authCodeURL, err := h.service.GoogleAuthCodeURL(req.RedirectURI)
	if err != nil {
		return nil, err
	}

	return &dto.RedirectResponse{
		Status:   http.StatusFound,
		Location: authCodeURL,
	}, nil
}

func (h *handlerImpl) Callback(ctx context.Context, req *dto.GoogleCallbackRequest) (*dto.RedirectResponse, error) {
	result, err := h.service.LoginWithGoogle(ctx, req.Code, req.State)
	if err != nil {
		return nil, err
	}

	return &dto.RedirectResponse{
		Status:   http.StatusFound,
		Location: appendTokens(result.RedirectURI, result.AccessToken, result.RefreshToken, result.IsNewUser, result.GoogleDisplayName, result.GoogleProfileImageURL),
	}, nil
}

func (h *handlerImpl) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.TokenPairResponse, error) {
	accessToken, refreshToken, err := h.service.RefreshAccessToken(ctx, req.Body.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// appendTokens adds the issued access/refresh tokens and whether the login
// created a new user to redirectURI as query parameters, preserving any
// query string it already has. For a new user, it also adds Google's name
// and profile picture as query parameters so the app can prefill its
// profile-creation form; the account itself is created without them.
func appendTokens(redirectURI, accessToken, refreshToken string, isNewUser bool, googleDisplayName, googleProfileImageURL string) string {
	u, err := url.Parse(redirectURI)
	if err != nil {
		return redirectURI
	}

	q := u.Query()
	q.Set("accessToken", accessToken)
	q.Set("refreshToken", refreshToken)
	q.Set("isNewUser", strconv.FormatBool(isNewUser))
	if isNewUser {
		if googleDisplayName != "" {
			q.Set("name", googleDisplayName)
		}
		if googleProfileImageURL != "" {
			q.Set("profileImageUrl", googleProfileImageURL)
		}
	}
	u.RawQuery = q.Encode()

	return u.String()
}
