package story

import (
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
)

// model is the MongoDB document shape for the stories collection. It stays
// private to this package - callers only ever see entity.Story.
type model struct {
	ID        string    `bson:"_id"`
	PartyID   string    `bson:"party_id"`
	UserID    string    `bson:"user_id"`
	Image     string    `bson:"image"`
	CreatedAt time.Time `bson:"created_at"`
}

func fromEntity(s entity.Story) model {
	return model{
		ID:        s.ID,
		PartyID:   s.PartyID,
		UserID:    s.UserID,
		Image:     s.Image,
		CreatedAt: s.CreatedAt,
	}
}

func (m model) toEntity() entity.Story {
	return entity.Story{
		ID:        m.ID,
		PartyID:   m.PartyID,
		UserID:    m.UserID,
		Image:     m.Image,
		CreatedAt: m.CreatedAt,
	}
}
