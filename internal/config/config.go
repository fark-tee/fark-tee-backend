package config

import (
	"os"

	env "github.com/caarlos0/env/v11"
	"github.com/fark-tee/go-kit/validatex"
	"github.com/joho/godotenv"
)

type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	CORS           CORSConfig
	InstagramOAuth InstagramOAuthConfig
	JWT            JWTConfig
	Storage        StorageConfig
}

// @WireSet("Config")
func New() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, err
	}

	if err := validatex.Struct(cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
