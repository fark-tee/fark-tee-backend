package partymember

import (
	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// model is the MongoDB document shape for the party_members collection. It
// stays private to this package - callers only ever see entity.PartyMember.
type model struct {
	ID               bson.ObjectID `bson:"_id"`
	PartyID          string        `bson:"party_id"`
	UserID           string        `bson:"user_id"`
	UserDisplayName  string        `bson:"user_display_name"`
	UserProfileImage string        `bson:"user_profile_image"`
	Status           string        `bson:"status"`
	TripStatus       string        `bson:"trip_status"`

	CheckInStatus            string `bson:"check_in_status"`
	CheckInRequestedByUserID string `bson:"check_in_requested_by_user_id"`
}

func fromEntity(m entity.PartyMember) (model, error) {
	id, err := mongoid.ToObjectID(m.ID)
	if err != nil {
		return model{}, err
	}

	checkInStatus := m.CheckInStatus
	if checkInStatus == "" {
		checkInStatus = entity.CheckInStatusNone
	}

	return model{
		ID:                       id,
		PartyID:                  m.PartyID,
		UserID:                   m.UserID,
		UserDisplayName:          m.UserDisplayName,
		UserProfileImage:         m.UserProfileImage,
		Status:                   string(m.Status),
		TripStatus:               string(m.TripStatus),
		CheckInStatus:            string(checkInStatus),
		CheckInRequestedByUserID: m.CheckInRequestedByUserID,
	}, nil
}

func (m model) toEntity() entity.PartyMember {
	checkInStatus := entity.CheckInStatus(m.CheckInStatus)
	if checkInStatus == "" {
		checkInStatus = entity.CheckInStatusNone
	}

	return entity.PartyMember{
		ID:                       mongoid.FromObjectID(m.ID),
		PartyID:                  m.PartyID,
		UserID:                   m.UserID,
		UserDisplayName:          m.UserDisplayName,
		UserProfileImage:         m.UserProfileImage,
		Status:                   entity.PartyMemberStatus(m.Status),
		TripStatus:               entity.TripStatus(m.TripStatus),
		CheckInStatus:            checkInStatus,
		CheckInRequestedByUserID: m.CheckInRequestedByUserID,
	}
}
