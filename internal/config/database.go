package config

type DatabaseConfig struct {
	URI      string `env:"MONGODB_URI" validate:"required"`
	Database string `env:"MONGODB_DATABASE" validate:"required"`
}
