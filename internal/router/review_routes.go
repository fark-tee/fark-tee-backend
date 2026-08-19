package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerReviewRoutes(api huma.API, handlers *handler.Handlers) {
	security := []map[string][]string{
		{humax.BearerAuthSecurityScheme: {}},
	}

	huma.Register(api, huma.Operation{
		Method:   http.MethodPost,
		Path:     "/parties/{partyId}/members/{userId}/review",
		Summary:  "Review a party member who has arrived at the destination",
		Security: security,
	}, humax.Wrap201(handlers.Review.CreateReview))

	huma.Register(api, huma.Operation{
		Method:   http.MethodGet,
		Path:     "/parties/{partyId}/reviews",
		Summary:  "List the reviews the current user has left for this party",
		Security: security,
	}, humax.Wrap200(handlers.Review.ListMyReviews))
}
