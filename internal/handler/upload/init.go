package upload

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/service/upload"
	"github.com/fark-tee/fark-tee-backend/pkg/dto"
)

type Handler interface {
	Create(ctx context.Context, req *dto.UploadImageRequest) (*dto.UploadImageResponse, error)
}

type handlerImpl struct {
	service upload.Service
}

// @WireSet("Handler")
func New(service upload.Service) Handler {
	return &handlerImpl{service: service}
}
