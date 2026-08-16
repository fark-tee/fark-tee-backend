package config

type StorageConfig struct {
	Endpoint  string `env:"S3_ENDPOINT" validate:"required"`
	Region    string `env:"S3_REGION" validate:"required"`
	Bucket    string `env:"S3_BUCKET" validate:"required"`
	AccessKey string `env:"S3_ACCESS_KEY" validate:"required"`
	SecretKey string `env:"S3_SECRET_KEY" validate:"required"`
}
