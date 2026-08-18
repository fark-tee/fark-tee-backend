package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerPartyRoutes(api huma.API, handlers *handler.Handlers) {
	security := []map[string][]string{
		{humax.BearerAuthSecurityScheme: {}},
	}

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/parties",
		Summary:  "Create a party",
		Security: security,
	}, humax.Wrap201(handlers.Party.Create))

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/parties/{partyId}/members",
		Summary:  "Invite a user to a party (owner only)",
		Security: security,
	}, humax.Wrap201(handlers.Party.Invite))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/me/invites",
		Summary:  "List the current user's pending party invites",
		Security: security,
	}, humax.Wrap200(handlers.Party.MyInvites))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/me/parties",
		Summary:  "List the current user's parties",
		Security: security,
	}, humax.Wrap200(handlers.Party.MyParties))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/parties/{partyId}",
		Summary:  "Get a party by ID (members only)",
		Security: security,
	}, humax.Wrap200(handlers.Party.Get))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/parties/{partyId}/members",
		Summary:  "List a party's members",
		Security: security,
	}, humax.Wrap200(handlers.Party.ListMembers))

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/parties/{partyId}/members/accept",
		Summary:  "Accept a pending party invite",
		Security: security,
	}, humax.Wrap200(handlers.Party.AcceptInvite))

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/parties/{partyId}/members/decline",
		Summary:       "Decline a pending party invite, removing yourself from the party",
		Security:      security,
		DefaultStatus: http.StatusNoContent,
	}, handlers.Party.DeclineInvite)

	huma.Register(api, huma.Operation{
		Method:        http.MethodDelete,
		Path:          "/parties/{partyId}/members/{userId}",
		Summary:       "Remove a member from a party (owner only)",
		Security:      security,
		DefaultStatus: http.StatusNoContent,
	}, handlers.Party.RemoveMember)
}
