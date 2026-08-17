package config

type GoogleOAuthConfig struct {
	ClientID     string `env:"GOOGLE_CLIENT_ID" validate:"required"`
	ClientSecret string `env:"GOOGLE_CLIENT_SECRET" validate:"required"`
	RedirectURL  string `env:"GOOGLE_REDIRECT_URL" validate:"required"`
}
