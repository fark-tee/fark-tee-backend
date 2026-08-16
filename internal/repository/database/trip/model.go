package trip

import (
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
)

// model is the MongoDB document shape for the trips collection. It stays
// private to this package - callers only ever see entity.Trip.
type model struct {
	ID        string    `bson:"_id"`
	PartyID   string    `bson:"party_id"`
	UserID    string    `bson:"user_id"`
	Direction string    `bson:"direction"`
	StartedAt time.Time `bson:"started_at"`
}

func fromEntity(t entity.Trip) model {
	return model{
		ID:        t.ID,
		PartyID:   t.PartyID,
		UserID:    t.UserID,
		Direction: string(t.Direction),
		StartedAt: t.StartedAt,
	}
}

func (m model) toEntity() entity.Trip {
	return entity.Trip{
		ID:        m.ID,
		PartyID:   m.PartyID,
		UserID:    m.UserID,
		Direction: entity.TripDirection(m.Direction),
		StartedAt: m.StartedAt,
	}
}
