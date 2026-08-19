package dto

import "github.com/danielgtaylor/huma/v2"

type UserResponse struct {
	ID              string  `json:"id"`
	ProfileImageURL string  `json:"profileImageUrl"`
	DisplayName     string  `json:"displayName"`
	Username        string  `json:"username"`
	GoogleUserID    string  `json:"googleUserId"`
	Rating          float64 `json:"rating"`
	RatingCount     int     `json:"ratingCount"`
	OnTimeCount     int     `json:"onTimeCount"`
	LateCount       int     `json:"lateCount"`
	AccessToken     string  `json:"accessToken,omitempty"`

	EmergencyContactName  string `json:"emergencyContactName,omitempty"`
	EmergencyContactPhone string `json:"emergencyContactPhone,omitempty"`
}

type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

type SearchUsersRequest struct {
	Query string `query:"q" required:"true"`
}

type GetMeRequest struct{}

type UpdateMeRequest struct {
	Body struct {
		DisplayName string `json:"displayName" required:"true"`
		Username    string `json:"username" required:"true" minLength:"3" maxLength:"20" pattern:"^[a-zA-Z0-9_]+$"`

		// EmergencyContactName/Phone are optional - see
		// entity.User.EmergencyContactName's doc comment for why.
		EmergencyContactName  string `json:"emergencyContactName,omitempty" maxLength:"100"`
		EmergencyContactPhone string `json:"emergencyContactPhone,omitempty" maxLength:"20"`
	}
}

type UploadProfileImageFormData struct {
	Image huma.FormFile `form:"image" contentType:"image/*" required:"true"`
}

type UploadProfileImageRequest struct {
	RawBody huma.MultipartFormFiles[UploadProfileImageFormData]
}
