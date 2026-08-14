package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerUserRoutes(api huma.API, handlers *handler.Handlers) {
	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/users/search",
		Summary: "Search users by display name, ID, or Instagram user ID",
		Security: []map[string][]string{
			{humax.BearerAuthSecurityScheme: {}},
		},
	}, humax.Wrap200(handlers.User.Search))
}
