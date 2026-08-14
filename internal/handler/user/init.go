package user

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/user"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	Search(ctx context.Context, req *dto.SearchUsersRequest) (*dto.UsersResponse, error)
}

type handlerImpl struct {
	service user.Service
}

// @WireSet("Handler")
func New(service user.Service) Handler {
	return &handlerImpl{service: service}
}
