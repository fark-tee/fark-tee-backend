package devicetoken

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/devicetoken"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	Register(ctx context.Context, req *dto.RegisterDeviceTokenRequest) (*dto.RegisterDeviceTokenResponse, error)
	Unregister(ctx context.Context, req *dto.DeleteDeviceTokenRequest) (*dto.DeleteDeviceTokenResponse, error)
}

type handlerImpl struct {
	service devicetoken.Service
}

// @WireSet("Handler")
func New(service devicetoken.Service) Handler {
	return &handlerImpl{service: service}
}
