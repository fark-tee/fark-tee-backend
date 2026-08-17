package config

import "time"

type JWTConfig struct {
	Secret     string        `env:"JWT_SECRET" validate:"required"`
	TTL        time.Duration `env:"JWT_TTL" envDefault:"24h" validate:"required"`
	RefreshTTL time.Duration `env:"JWT_REFRESH_TTL" envDefault:"720h" validate:"required"`
}
