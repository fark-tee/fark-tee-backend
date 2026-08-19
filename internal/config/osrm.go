package config

type OSRMConfig struct {
	BaseURL string `env:"OSRM_BASE_URL" validate:"required"`
	Profile string `env:"OSRM_PROFILE" envDefault:"driving"`
}
