package user

import (
	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// model is the MongoDB document shape for the users collection. It stays
// private to this package - callers only ever see entity.User.
type model struct {
	ID              bson.ObjectID `bson:"_id"`
	ProfileImageURL string        `bson:"profile_image_url"`
	DisplayName     string        `bson:"display_name"`
	Username        string        `bson:"username"`
	GoogleUserID    string        `bson:"google_user_id"`
	Rating          float64       `bson:"rating"`
	RatingCount     int           `bson:"rating_count"`
	OnTimeCount     int           `bson:"on_time_count"`
	LateCount       int           `bson:"late_count"`

	EmergencyContactName  string `bson:"emergency_contact_name"`
	EmergencyContactPhone string `bson:"emergency_contact_phone"`
}

func fromEntity(u entity.User) (model, error) {
	id, err := mongoid.ToObjectID(u.ID)
	if err != nil {
		return model{}, err
	}

	return model{
		ID:                    id,
		ProfileImageURL:       u.ProfileImageURL,
		DisplayName:           u.DisplayName,
		Username:              u.Username,
		GoogleUserID:          u.GoogleUserID,
		Rating:                u.Rating,
		RatingCount:           u.RatingCount,
		OnTimeCount:           u.OnTimeCount,
		LateCount:             u.LateCount,
		EmergencyContactName:  u.EmergencyContactName,
		EmergencyContactPhone: u.EmergencyContactPhone,
	}, nil
}

func (m model) toEntity() entity.User {
	return entity.User{
		ID:                    mongoid.FromObjectID(m.ID),
		ProfileImageURL:       m.ProfileImageURL,
		DisplayName:           m.DisplayName,
		Username:              m.Username,
		GoogleUserID:          m.GoogleUserID,
		Rating:                m.Rating,
		RatingCount:           m.RatingCount,
		OnTimeCount:           m.OnTimeCount,
		LateCount:             m.LateCount,
		EmergencyContactName:  m.EmergencyContactName,
		EmergencyContactPhone: m.EmergencyContactPhone,
	}
}
