package googleoauth

import (
	"context"
	"net/http"
	"net/url"

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
		Location: appendTokens(result.RedirectURI, result.AccessToken, result.RefreshToken),
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

// appendTokens adds the issued access/refresh tokens to redirectURI as query
// parameters, preserving any query string it already has.
func appendTokens(redirectURI, accessToken, refreshToken string) string {
	u, err := url.Parse(redirectURI)
	if err != nil {
		return redirectURI
	}

	q := u.Query()
	q.Set("accessToken", accessToken)
	q.Set("refreshToken", refreshToken)
	u.RawQuery = q.Encode()

	return u.String()
}
