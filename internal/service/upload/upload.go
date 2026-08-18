package upload

import (
	"context"
	"io"
	"path/filepath"

	"github.com/fark-tee/go-kit/idx"
)

func (s *serviceImpl) UploadImage(ctx context.Context, image io.Reader, size int64, contentType, filename string) (string, error) {
	key := "uploads/" + idx.NewUUID() + filepath.Ext(filename)

	return s.uploader.UploadPublic(ctx, key, image, size, contentType)
}
