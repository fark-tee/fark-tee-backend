package user

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) Search(ctx context.Context, req *dto.SearchUsersRequest) (*dto.UsersResponse, error) {
	users, err := h.service.Search(ctx, req.Query)
	if err != nil {
		return nil, err
	}

	resp := &dto.UsersResponse{
		Users: make([]dto.UserResponse, 0, len(users)),
	}
	for _, u := range users {
		resp.Users = append(resp.Users, *toUserResponse(u))
	}

	return resp, nil
}

func toUserResponse(u entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:              u.ID,
		ProfileImageURL: u.ProfileImageURL,
		DisplayName:     u.DisplayName,
		GoogleUserID:    u.GoogleUserID,
		Rating:          u.Rating,
		RatingCount:     u.RatingCount,
		OnTimeCount:     u.OnTimeCount,
		LateCount:       u.LateCount,
	}
}
