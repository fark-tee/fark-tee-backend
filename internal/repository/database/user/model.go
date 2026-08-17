package user

import "github.com/fark-tee/fark-tee-backend/internal/entity"

// model is the MongoDB document shape for the users collection. It stays
// private to this package - callers only ever see entity.User.
type model struct {
	ID              string  `bson:"_id"`
	ProfileImageURL string  `bson:"profile_image_url"`
	DisplayName     string  `bson:"display_name"`
	GoogleUserID    string  `bson:"google_user_id"`
	Rating          float64 `bson:"rating"`
	RatingCount     int     `bson:"rating_count"`
	OnTimeCount     int     `bson:"on_time_count"`
	LateCount       int     `bson:"late_count"`
}

func fromEntity(u entity.User) model {
	return model{
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

func (m model) toEntity() entity.User {
	return entity.User{
		ID:              m.ID,
		ProfileImageURL: m.ProfileImageURL,
		DisplayName:     m.DisplayName,
		GoogleUserID:    m.GoogleUserID,
		Rating:          m.Rating,
		RatingCount:     m.RatingCount,
		OnTimeCount:     m.OnTimeCount,
		LateCount:       m.LateCount,
	}
}
