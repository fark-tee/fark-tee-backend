package user

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/middleware/authmw"
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

func (h *handlerImpl) GetMe(ctx context.Context, _ *dto.GetMeRequest) (*dto.UserResponse, error) {
	userID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	u, err := h.service.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

func (h *handlerImpl) UpdateMe(ctx context.Context, req *dto.UpdateMeRequest) (*dto.UserResponse, error) {
	userID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	u, err := h.service.UpdateProfile(ctx, userID, req.Body.DisplayName, req.Body.Username)
	if err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

func (h *handlerImpl) UploadProfileImage(ctx context.Context, req *dto.UploadProfileImageRequest) (*dto.UserResponse, error) {
	userID, err := authmw.RequireUserID(ctx)
	if err != nil {
		return nil, err
	}

	image := req.RawBody.Data().Image

	u, err := h.service.UploadProfileImage(ctx, userID, image, image.Size, image.ContentType, image.Filename)
	if err != nil {
		return nil, err
	}

	return toUserResponse(u), nil
}

func toUserResponse(u entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:              u.ID,
		ProfileImageURL: u.ProfileImageURL,
		DisplayName:     u.DisplayName,
		Username:        u.Username,
		GoogleUserID:    u.GoogleUserID,
		Rating:          u.Rating,
		RatingCount:     u.RatingCount,
		OnTimeCount:     u.OnTimeCount,
		LateCount:       u.LateCount,
	}
}
