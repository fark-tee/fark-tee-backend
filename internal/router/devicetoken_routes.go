package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fark-tee/go-kit/humax"

	"github.com/fark-tee/fark-tee-backend/internal/handler"
)

func registerDeviceTokenRoutes(api huma.API, handlers *handler.Handlers) {
	security := []map[string][]string{
		{humax.BearerAuthSecurityScheme: {}},
	}

	huma.Register(api, huma.Operation{
		Method:        http.MethodPost,
		Path:          "/device-tokens",
		Summary:       "Register (or move ownership of) an FCM device token for the current user",
		Security:      security,
		DefaultStatus: http.StatusNoContent,
	}, handlers.DeviceToken.Register)

	huma.Register(api, huma.Operation{
		Method:        http.MethodDelete,
		Path:          "/device-tokens",
		Summary:       "Unregister an FCM device token (e.g. on logout)",
		Security:      security,
		DefaultStatus: http.StatusNoContent,
	}, handlers.DeviceToken.Unregister)
}
