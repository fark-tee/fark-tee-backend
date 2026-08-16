package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerTripRoutes(api huma.API, handlers *handler.Handlers) {
	security := []map[string][]string{
		{humax.BearerAuthSecurityScheme: {}},
	}

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/parties/{partyId}/trips",
		Summary:  "Start a depart or return trip and record its first position",
		Security: security,
	}, humax.Wrap201(handlers.Trip.StartTrip))

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/parties/{partyId}/positions",
		Summary:  "Record a new position for the current user's active trip",
		Security: security,
	}, humax.Wrap201(handlers.Trip.UpdatePosition))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/parties/{partyId}/positions",
		Summary:  "Get the latest recorded position of every party member",
		Security: security,
	}, humax.Wrap200(handlers.Trip.GetPartyPositions))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/parties/{partyId}/members/{userId}/position",
		Summary:  "Get the latest recorded position of a party member",
		Security: security,
	}, humax.Wrap200(handlers.Trip.GetMemberPosition))
}
