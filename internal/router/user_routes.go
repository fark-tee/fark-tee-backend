package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerUserRoutes(api huma.API, handlers *handler.Handlers) {
	security := []map[string][]string{
		{humax.BearerAuthSecurityScheme: {}},
	}

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/users/search",
		Summary:  "Search users by display name, ID, or Google user ID",
		Security: security,
	}, humax.Wrap200(handlers.User.Search))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/me",
		Summary:  "Get the current user's profile",
		Security: security,
	}, humax.Wrap200(handlers.User.GetMe))

	huma.Register(api, huma.Operation{
		Method:   http.MethodPatch,
		Path:     "/me",
		Summary:  "Update the current user's display name",
		Security: security,
	}, humax.Wrap200(handlers.User.UpdateMe))

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/me/profile-image",
		Summary:  "Upload the current user's profile image",
		Security: security,
	}, humax.Wrap200(handlers.User.UploadProfileImage))
}
