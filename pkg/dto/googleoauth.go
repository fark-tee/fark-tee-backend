package dto

type GoogleStartRequest struct {
	// RedirectURI is the mobile app deeplink to send the user back to once
	// Google login completes, e.g. "farktee://auth/callback".
	RedirectURI string `query:"redirect_uri" required:"true"`
}

type GoogleCallbackRequest struct {
	Code  string `query:"code" required:"true"`
	State string `query:"state" required:"true"`
}
