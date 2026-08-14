package instagramoauth

import (
	"context"
	"net/http"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) Start(_ context.Context, _ *dto.InstagramStartRequest) (*dto.RedirectResponse, error) {
	return &dto.RedirectResponse{
		Status:   http.StatusFound,
		Location: h.service.InstagramAuthCodeURL(""),
	}, nil
}

func (h *handlerImpl) Callback(ctx context.Context, req *dto.InstagramCallbackRequest) (*dto.UserResponse, error) {
	user, err := h.service.LoginWithInstagram(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

func toUserResponse(u entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:              u.ID,
		ProfileImageURL: u.ProfileImageURL,
		DisplayName:     u.DisplayName,
		InstagramUserID: u.InstagramUserID,
		Rating:          u.Rating,
		RatingCount:     u.RatingCount,
		OnTimeCount:     u.OnTimeCount,
		LateCount:       u.LateCount,
	}
}
