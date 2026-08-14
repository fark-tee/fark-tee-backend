package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerSavedLocationRoutes(api huma.API, handlers *handler.Handlers) {
	huma.Register(api, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/users/{userId}/saved-locations",
		Summary: "Create a saved location for a user",
	}, humax.Wrap201(handlers.SavedLocation.Create))

	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/users/{userId}/saved-locations",
		Summary: "List a user's saved locations",
	}, humax.Wrap200(handlers.SavedLocation.List))

	huma.Register(api, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/saved-locations/{id}",
		Summary: "Get a saved location by ID",
	}, humax.Wrap200(handlers.SavedLocation.Get))

	huma.Register(api, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/saved-locations/{id}",
		Summary: "Update a saved location",
	}, humax.Wrap200(handlers.SavedLocation.Update))

	huma.Register(api, huma.Operation{
		Method:        http.MethodDelete,
		Path:          "/saved-locations/{id}",
		Summary:       "Delete a saved location",
		DefaultStatus: http.StatusNoContent,
	}, handlers.SavedLocation.Delete)
}
