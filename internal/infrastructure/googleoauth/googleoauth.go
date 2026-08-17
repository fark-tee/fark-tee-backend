package googleoauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fark-tee/fark-tee-backend/internal/config"
)

const (
	authorizeURL = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL     = "https://oauth2.googleapis.com/token"
	userInfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"
)

// Profile is the subset of a Google account's profile returned by the
// OpenID Connect userinfo endpoint.
type Profile struct {
	GoogleUserID    string
	Name            string
	Email           string
	ProfileImageURL string
}

type Verifier struct {
	cfg *config.Config
}

// @WireSet("Infrastructure")
func NewVerifier(cfg *config.Config) *Verifier {
	return &Verifier{cfg: cfg}
}

// AuthCodeURL builds the URL the client should be redirected to in order to
// grant this app access to their Google account. state is echoed back on the
// callback so the caller can correlate the request; passing "" omits it.
func (v *Verifier) AuthCodeURL(state string) string {
	params := url.Values{
		"client_id":     {v.cfg.GoogleOAuth.ClientID},
		"redirect_uri":  {v.cfg.GoogleOAuth.RedirectURL},
		"response_type": {"code"},
		"scope":         {"openid email profile"},
	}

	if state != "" {
		params.Set("state", state)
	}

	return authorizeURL + "?" + params.Encode()
}

// Exchange trades an authorization code for the account's profile, following
// the code -> access token -> profile fields flow.
func (v *Verifier) Exchange(ctx context.Context, code string) (Profile, error) {
	accessToken, err := v.exchangeCodeForToken(ctx, code)
	if err != nil {
		return Profile{}, err
	}

	return v.fetchProfile(ctx, accessToken)
}

func (v *Verifier) exchangeCodeForToken(ctx context.Context, code string) (accessToken string, err error) {
	form := url.Values{
		"client_id":     {v.cfg.GoogleOAuth.ClientID},
		"client_secret": {v.cfg.GoogleOAuth.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {v.cfg.GoogleOAuth.RedirectURL},
		"code":          {code},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var body struct {
		AccessToken string `json:"access_token"`
	}

	if err := doJSON(req, &body); err != nil {
		return "", fmt.Errorf("exchange google code: %w", err)
	}

	return body.AccessToken, nil
}

func (v *Verifier) fetchProfile(ctx context.Context, accessToken string) (Profile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userInfoURL, nil)
	if err != nil {
		return Profile{}, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	var body struct {
		Sub     string `json:"sub"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}

	if err := doJSON(req, &body); err != nil {
		return Profile{}, fmt.Errorf("fetch google profile: %w", err)
	}

	return Profile{
		GoogleUserID:    body.Sub,
		Name:            body.Name,
		Email:           body.Email,
		ProfileImageURL: body.Picture,
	}, nil
}

func doJSON(req *http.Request, out any) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return json.Unmarshal(respBody, out)
}
