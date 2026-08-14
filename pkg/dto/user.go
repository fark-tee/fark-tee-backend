package dto

type UserResponse struct {
	ID              string  `json:"id"`
	ProfileImageURL string  `json:"profileImageUrl"`
	DisplayName     string  `json:"displayName"`
	InstagramUserID string  `json:"instagramUserId"`
	Rating          float64 `json:"rating"`
	RatingCount     int     `json:"ratingCount"`
	OnTimeCount     int     `json:"onTimeCount"`
	LateCount       int     `json:"lateCount"`
}
