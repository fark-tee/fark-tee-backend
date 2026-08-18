package devicetoken

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
)

func (s *serviceImpl) Register(ctx context.Context, userID, token, platform string) error {
	return s.repo.Upsert(ctx, entity.DeviceToken{
		UserID:   userID,
		Token:    token,
		Platform: entity.DevicePlatform(platform),
	})
}

func (s *serviceImpl) Unregister(ctx context.Context, token string) error {
	return s.repo.DeleteByToken(ctx, token)
}
