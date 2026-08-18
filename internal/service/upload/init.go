package upload

import (
	"context"
	"io"

	"github.com/fark-tee/fark-tee-backend/internal/infrastructure/storage"
)

type Service interface {
	// UploadImage uploads image as a public object under uploads/ and
	// returns its public URL.
	UploadImage(ctx context.Context, image io.Reader, size int64, contentType, filename string) (string, error)
}

type serviceImpl struct {
	uploader *storage.Uploader
}

// @WireSet("Service")
func New(uploader *storage.Uploader) Service {
	return &serviceImpl{uploader: uploader}
}
