package googleoauth

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/auth"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	Start(ctx context.Context, req *dto.GoogleStartRequest) (*dto.RedirectResponse, error)
	Callback(ctx context.Context, req *dto.GoogleCallbackRequest) (*dto.RedirectResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.TokenPairResponse, error)
}

type handlerImpl struct {
	service auth.Service
}

// @WireSet("Handler")
func New(service auth.Service) Handler {
	return &handlerImpl{service: service}
}
