package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/fark-tee/fark-tee-backend/internal/config"
	"github.com/fark-tee/go-kit/loggy"
)

// @WireSet("Infrastructure")
func NewLogger(cfg *config.Config) *slog.Logger {
	var level slog.Level
	if err := level.UnmarshalText([]byte(cfg.Server.LogLevel)); err != nil {
		level = slog.LevelInfo
	}

	var handler slog.Handler
	if strings.EqualFold(cfg.Server.LogFormat, "JSON") {
		handler = loggy.NewJSONHandler(os.Stdout, &loggy.JSONHandlerOptions{
			Level: level.Level(),
		})
	} else {
		handler = loggy.NewTextHandler(os.Stdout, &loggy.TextHandlerOptions{
			Level: level.Level(),
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
