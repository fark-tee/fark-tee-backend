package devicetoken

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/middleware/authmw"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) Register(ctx context.Context, req *dto.RegisterDeviceTokenRequest) (*dto.RegisterDeviceTokenResponse, error) {
	userID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.service.Register(ctx, userID, req.Body.Token, req.Body.Platform); err != nil {
		return nil, err
	}

	return &dto.RegisterDeviceTokenResponse{}, nil
}

func (h *handlerImpl) Unregister(ctx context.Context, req *dto.DeleteDeviceTokenRequest) (*dto.DeleteDeviceTokenResponse, error) {
	if _, err := authmw.RequireUserID(ctx); err != nil {
		return nil, err
	}

	if err := h.service.Unregister(ctx, req.Body.Token); err != nil {
		return nil, err
	}

	return &dto.DeleteDeviceTokenResponse{}, nil
}
