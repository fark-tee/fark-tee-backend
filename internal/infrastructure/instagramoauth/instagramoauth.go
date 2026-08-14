package instagramoauth

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
	authorizeURL    = "https://www.instagram.com/oauth/authorize"
	accessTokenURL  = "https://api.instagram.com/oauth/access_token"
	graphProfileURL = "https://graph.instagram.com/me"
)

// Profile is the subset of an Instagram account's public profile the
// "Instagram API with Instagram Login" exposes to a Basic/Business
// authorized app. Instagram does not return a display name or profile
// picture through this API - Username is the closest available field.
type Profile struct {
	InstagramUserID string
	Username        string
}

type Verifier struct {
	cfg *config.Config
}

// @WireSet("Infrastructure")
func NewVerifier(cfg *config.Config) *Verifier {
	return &Verifier{cfg: cfg}
}

// AuthCodeURL builds the URL the client should be redirected to in order to
// grant this app access to their Instagram account. state is echoed back on
// the callback so the caller can correlate the request; passing "" omits it.
func (v *Verifier) AuthCodeURL(state string) string {
	params := url.Values{
		"client_id":     {v.cfg.InstagramOAuth.ClientID},
		"redirect_uri":  {v.cfg.InstagramOAuth.RedirectURL},
		"response_type": {"code"},
		"scope":         {"instagram_business_basic"},
	}

	if state != "" {
		params.Set("state", state)
	}

	return authorizeURL + "?" + params.Encode()
}

// Exchange trades an authorization code for the account's profile, following
// the code -> short-lived access token -> profile fields flow.
func (v *Verifier) Exchange(ctx context.Context, code string) (Profile, error) {
	token, userID, err := v.exchangeCodeForToken(ctx, code)
	if err != nil {
		return Profile{}, err
	}

	return v.fetchProfile(ctx, token, userID)
}

func (v *Verifier) exchangeCodeForToken(ctx context.Context, code string) (accessToken, userID string, err error) {
	form := url.Values{
		"client_id":     {v.cfg.InstagramOAuth.ClientID},
		"client_secret": {v.cfg.InstagramOAuth.ClientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {v.cfg.InstagramOAuth.RedirectURL},
		"code":          {code},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, accessTokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var body struct {
		AccessToken string `json:"access_token"`
		UserID      string `json:"user_id"`
	}

	if err := doJSON(req, &body); err != nil {
		return "", "", fmt.Errorf("exchange instagram code: %w", err)
	}

	return body.AccessToken, body.UserID, nil
}

func (v *Verifier) fetchProfile(ctx context.Context, accessToken, userID string) (Profile, error) {
	params := url.Values{
		"fields":       {"id,username"},
		"access_token": {accessToken},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, graphProfileURL+"?"+params.Encode(), nil)
	if err != nil {
		return Profile{}, err
	}

	var body struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	if err := doJSON(req, &body); err != nil {
		return Profile{}, fmt.Errorf("fetch instagram profile: %w", err)
	}

	if body.ID == "" {
		body.ID = userID
	}

	return Profile{InstagramUserID: body.ID, Username: body.Username}, nil
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
