package user

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
)

func (s *serviceImpl) Search(ctx context.Context, query string) ([]entity.User, error) {
	return s.repo.Search(ctx, query)
}
