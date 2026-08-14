package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerInstagramOAuthRoutes(api huma.API, handlers *handler.Handlers) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/auth/instagram/start",
		Summary:       "Redirect to Instagram's OAuth consent screen",
		DefaultStatus: http.StatusFound,
	}, handlers.InstagramOAuth.Start)

	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/auth/instagram/callback",
		Summary: "Complete Instagram OAuth login, creating the user if needed",
	}, humax.Wrap200(handlers.InstagramOAuth.Callback))
}
