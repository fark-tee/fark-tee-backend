package upload

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

func (h *handlerImpl) Create(ctx context.Context, req *dto.UploadImageRequest) (*dto.UploadImageResponse, error) {
	image := req.RawBody.Data().Image

	url, err := h.service.UploadImage(ctx, image, image.Size, image.ContentType, image.Filename)
	if err != nil {
		return nil, err
	}

	return &dto.UploadImageResponse{URL: url}, nil
}
