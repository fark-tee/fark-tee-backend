package router

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
	"github.com/fark-tee/fark-tee-backend/internal/middleware/authmw"
)

func RegisterRoutes(api huma.API, handlers *handler.Handlers, authMW *authmw.Middleware) {
	// Registered on the root API (not a sub-group) so every route logs,
	// including the health check and the public Google OAuth routes below -
	// huma.Group.Middlewares() prepends its parent's, so anything added
	// here is inherited by /v1 and the protected group too.
	//
	// Full request/response bodies are logged, which means /v1/auth/refresh
	// (and the Google OAuth callback) will write access/refresh tokens to
	// the log output - fine for local/dev, but reconsider before running
	// this in an environment where logs aren't tightly access-controlled.
	// Use humax.LogWithOptions(humax.LogOptions{HideRequestBody: true,
	// HideResponseBody: true}) instead if that's a concern.
	api.UseMiddleware(humax.LogFull)

	humax.RegisterHealthCheckHandler(api)

	v1 := huma.NewGroup(api, "/v1")

	registerGoogleOAuthRoutes(v1, handlers)
	registerUploadRoutes(v1, handlers)

	protected := huma.NewGroup(v1, "")
	protected.UseMiddleware(authMW.RequireAuth(protected))

	registerSavedLocationRoutes(protected, handlers)
	registerUserRoutes(protected, handlers)
	registerPartyRoutes(protected, handlers)
	registerStoryRoutes(protected, handlers)
	registerTripRoutes(protected, handlers)
	registerDeviceTokenRoutes(protected, handlers)
}
