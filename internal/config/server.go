package config

type ServerConfig struct {
	Port      string `env:"PORT" envDefault:"8080" validate:"required"`
	LogLevel  string `env:"LOG_LEVEL" envDefault:"INFO" validate:"required,oneofci=DEBUG INFO WARN ERROR"`
	LogFormat string `env:"LOG_FORMAT" envDefault:"TEXT" validate:"required,oneofci=TEXT JSON"`
}
