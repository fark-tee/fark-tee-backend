package config

type InstagramOAuthConfig struct {
	ClientID     string `env:"INSTAGRAM_CLIENT_ID" validate:"required"`
	ClientSecret string `env:"INSTAGRAM_CLIENT_SECRET" validate:"required"`
	RedirectURL  string `env:"INSTAGRAM_REDIRECT_URL" validate:"required"`
}
