package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerSavedLocationRoutes(api huma.API, handlers *handler.Handlers) {
	security := []map[string][]string{
		{humax.BearerAuthSecurityScheme: {}},
	}

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/saved-locations",
		Summary:  "Create a saved location for the current user",
		Security: security,
	}, humax.Wrap201(handlers.SavedLocation.Create))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/saved-locations",
		Summary:  "List the current user's saved locations",
		Security: security,
	}, humax.Wrap200(handlers.SavedLocation.List))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/saved-locations/{id}",
		Summary:  "Get a saved location by ID",
		Security: security,
	}, humax.Wrap200(handlers.SavedLocation.Get))

	huma.Register(api, huma.Operation{
		Method:   http.MethodPut,
		Path:     "/saved-locations/{id}",
		Summary:  "Update a saved location",
		Security: security,
	}, humax.Wrap200(handlers.SavedLocation.Update))

	huma.Register(api, huma.Operation{
		Method:        http.MethodDelete,
		Path:          "/saved-locations/{id}",
		Summary:       "Delete a saved location",
		Security:      security,
		DefaultStatus: http.StatusNoContent,
	}, handlers.SavedLocation.Delete)
}
