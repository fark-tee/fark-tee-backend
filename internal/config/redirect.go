package config

// RedirectConfig controls where the backend is allowed to send a client
// (the mobile app, or a browser during local debugging) after a
// browser-based flow such as OAuth completes.
type RedirectConfig struct {
	// AllowedPrefixes lists the redirect_uri prefixes the backend will
	// redirect to, e.g. the mobile app's deeplink scheme ("farktee://") and,
	// for local debugging, a web URL ("http://localhost:3000"). Requests
	// asking to redirect anywhere else are rejected so the backend can't be
	// used as an open redirect.
	AllowedPrefixes []string `env:"ALLOWED_REDIRECT_PREFIXES" envSeparator:"," validate:"required"`
}
