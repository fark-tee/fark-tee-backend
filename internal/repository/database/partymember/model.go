package partymember

import "github.com/fark-tee/fark-tee-backend/internal/entity"

// model is the MongoDB document shape for the party_members collection. It
// stays private to this package - callers only ever see entity.PartyMember.
type model struct {
	ID              string `bson:"_id"`
	PartyID         string `bson:"party_id"`
	UserID          string `bson:"user_id"`
	UserDisplayName string `bson:"user_display_name"`
	Status          string `bson:"status"`
}

func fromEntity(m entity.PartyMember) model {
	return model{
		ID:              m.ID,
		PartyID:         m.PartyID,
		UserID:          m.UserID,
		UserDisplayName: m.UserDisplayName,
		Status:          string(m.Status),
	}
}

func (m model) toEntity() entity.PartyMember {
	return entity.PartyMember{
		ID:              m.ID,
		PartyID:         m.PartyID,
		UserID:          m.UserID,
		UserDisplayName: m.UserDisplayName,
		Status:          entity.PartyMemberStatus(m.Status),
	}
}
