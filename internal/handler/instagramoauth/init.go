package instagramoauth

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/auth"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	Start(ctx context.Context, req *dto.InstagramStartRequest) (*dto.RedirectResponse, error)
	Callback(ctx context.Context, req *dto.InstagramCallbackRequest) (*dto.UserResponse, error)
}

type handlerImpl struct {
	service auth.Service
}

// @WireSet("Handler")
func New(service auth.Service) Handler {
	return &handlerImpl{service: service}
}
