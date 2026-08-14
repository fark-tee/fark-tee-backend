package router

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func RegisterRoutes(api huma.API, handlers *handler.Handlers) {
	humax.RegisterHealthCheckHandler(api)

	v1 := huma.NewGroup(api, "/v1")

	registerInstagramOAuthRoutes(v1, handlers)
	registerSavedLocationRoutes(v1, handlers)
}
