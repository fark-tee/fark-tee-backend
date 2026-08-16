package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerStoryRoutes(api huma.API, handlers *handler.Handlers) {
	security := []map[string][]string{
		{humax.BearerAuthSecurityScheme: {}},
	}

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/parties/{partyId}/stories",
		Summary:  "Create a story in a party (party members only)",
		Security: security,
	}, humax.Wrap201(handlers.Story.Create))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/parties/{partyId}/members/{userId}/stories",
		Summary:  "List a party member's stories",
		Security: security,
	}, humax.Wrap200(handlers.Story.ListByUser))

	huma.Register(api, huma.Operation{
		Method:        http.MethodDelete,
		Path:          "/parties/{partyId}/stories/{storyId}",
		Summary:       "Delete a story",
		Security:      security,
		DefaultStatus: http.StatusNoContent,
	}, handlers.Story.Delete)
}
