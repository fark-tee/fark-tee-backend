package router

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
	"github.com/fark-tee/fark-tee-backend/internal/middleware/authmw"
)

func RegisterRoutes(api huma.API, handlers *handler.Handlers, authMW *authmw.Middleware) {
	humax.RegisterHealthCheckHandler(api)

	v1 := huma.NewGroup(api, "/v1")

	registerInstagramOAuthRoutes(v1, handlers)

	protected := huma.NewGroup(v1, "")
	protected.UseMiddleware(authMW.RequireAuth(protected))

	registerSavedLocationRoutes(protected, handlers)
	registerUserRoutes(protected, handlers)
	registerPartyRoutes(protected, handlers)
}
