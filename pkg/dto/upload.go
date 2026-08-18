package dto

import "github.com/danielgtaylor/huma/v2"

type UploadImageFormData struct {
	Image huma.FormFile `form:"image" contentType:"image/*" required:"true"`
}

type UploadImageRequest struct {
	RawBody huma.MultipartFormFiles[UploadImageFormData]
}

type UploadImageResponse struct {
	URL string `json:"url"`
}
