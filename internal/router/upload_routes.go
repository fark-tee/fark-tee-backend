package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerUploadRoutes(api huma.API, handlers *handler.Handlers) {
	huma.Register(api, huma.Operation{
		Method:  http.MethodPost,
		Path:    "/upload",
		Summary: "Upload an image and get back its public URL",
	}, humax.Wrap200(handlers.Upload.Create))
}
