package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"github.com/fark-tee/fark-tee-backend/internal/config"
	"github.com/fark-tee/fark-tee-backend/internal/handler"
	"github.com/fark-tee/fark-tee-backend/internal/router"
	"github.com/fark-tee/go-kit/echox"
	"github.com/fark-tee/go-kit/humax"
)

type Server struct {
	cfg      *config.Config
	handlers *handler.Handlers
}

// Options is intentionally empty: configuration is already loaded via
// internal/config, not through humacli's flag/env parsing.
type Options struct{}

func New(cfg *config.Config, handlers *handler.Handlers) *Server {
	return &Server{
		cfg:      cfg,
		handlers: handlers,
	}
}

func (s *Server) Start() {
	cli := humacli.New(func(hook humacli.Hooks, _ *Options) {
		e := echo.New()

		e.Use(echox.InjectRequestIDMiddleware)
		e.Use(echox.InjectRequestMetaMiddleware)
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins:     s.cfg.CORS.AllowedOrigins,
			AllowMethods:     s.cfg.CORS.AllowedMethods,
			AllowHeaders:     s.cfg.CORS.AllowedHeaders,
			ExposeHeaders:    s.cfg.CORS.ExposedHeaders,
			AllowCredentials: s.cfg.CORS.AllowCredentials,
			MaxAge:           s.cfg.CORS.MaxAge,
		}))

		humaConfig := huma.DefaultConfig("fark-tee-backend", "v1.0.0")
		humaConfig.CreateHooks = nil // skip the default $schema/Link response transformer to keep the response envelope unchanged

		humaAPI := humaecho.New(e, humaConfig)

		humax.UseAppErrors()

		router.RegisterRoutes(humaAPI, s.handlers)

		server := &http.Server{
			Addr:    ":" + s.cfg.Server.Port,
			Handler: e,
		}

		hook.OnStart(func() {
			slog.Info("Starting server on port " + s.cfg.Server.Port)

			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("Failed to start server", slog.Any("error", err))

				os.Exit(1)
			}
		})

		hook.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			server.Shutdown(ctx)
		})
	})

	cli.Run()
}
