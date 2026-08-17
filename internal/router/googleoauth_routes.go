package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerGoogleOAuthRoutes(api huma.API, handlers *handler.Handlers) {
	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/auth/google/start",
		Summary:       "Redirect to Google's OAuth consent screen",
		DefaultStatus: http.StatusFound,
	}, handlers.GoogleOAuth.Start)

	huma.Register(api, huma.Operation{
		Method:        http.MethodGet,
		Path:          "/auth/google/callback",
		Summary:       "Complete Google OAuth login, creating the user if needed, then redirect to the mobile app deeplink",
		DefaultStatus: http.StatusFound,
	}, handlers.GoogleOAuth.Callback)

	huma.Register(api, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/auth/refresh",
		Summary: "Exchange a refresh token for a new access/refresh token pair",
	}, humax.Wrap200(handlers.GoogleOAuth.RefreshToken))
}
