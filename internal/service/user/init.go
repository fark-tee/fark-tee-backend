package user

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/user"
)

type Service interface {
	Search(ctx context.Context, query string) ([]entity.User, error)
}

type serviceImpl struct {
	repo user.Repository
}

// @WireSet("Service")
func New(repo user.Repository) Service {
	return &serviceImpl{repo: repo}
}
