package storage

import (
	"context"
	"io"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/fark-tee/fark-tee-backend/internal/config"
)

// Uploader stores objects in an S3-compatible bucket (AWS S3 or MinIO) and
// makes them publicly readable.
type Uploader struct {
	client        *s3.Client
	bucket        string
	publicURLBase string
}

// @WireSet("Infrastructure")
func NewUploader(ctx context.Context, cfg *config.Config) (*Uploader, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(cfg.Storage.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.Storage.AccessKey, cfg.Storage.SecretKey, "")),
	)
	if err != nil {
		return nil, err
	}

	endpoint := cfg.Storage.Endpoint
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = &endpoint
		o.UsePathStyle = true
	})

	return &Uploader{
		client:        client,
		bucket:        cfg.Storage.Bucket,
		publicURLBase: strings.TrimRight(endpoint, "/") + "/" + cfg.Storage.Bucket,
	}, nil
}

// UploadPublic uploads data under key, marks the object publicly readable,
// and returns its public URL.
func (u *Uploader) UploadPublic(ctx context.Context, key string, data io.Reader, size int64, contentType string) (string, error) {
	_, err := u.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        &u.bucket,
		Key:           &key,
		Body:          data,
		ContentLength: &size,
		ContentType:   &contentType,
		ACL:           types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", err
	}

	return u.publicURLBase + "/" + key, nil
}

// DeleteObject removes the object stored under key.
func (u *Uploader) DeleteObject(ctx context.Context, key string) error {
	_, err := u.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &u.bucket,
		Key:    &key,
	})

	return err
}

// KeyFromPublicURL extracts the object key from a URL previously returned by
// UploadPublic, so a caller holding only the public URL can still delete the
// underlying object. ok is false if publicURL wasn't produced by this
// uploader's bucket.
func (u *Uploader) KeyFromPublicURL(publicURL string) (key string, ok bool) {
	return strings.CutPrefix(publicURL, u.publicURLBase+"/")
}
