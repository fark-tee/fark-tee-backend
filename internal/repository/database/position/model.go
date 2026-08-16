package position

import (
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
)

// model is the MongoDB document shape for the positions collection. It stays
// private to this package - callers only ever see entity.Position.
type model struct {
	ID         string    `bson:"_id"`
	TripID     string    `bson:"trip_id"`
	PartyID    string    `bson:"party_id"`
	UserID     string    `bson:"user_id"`
	Lat        float64   `bson:"lat"`
	Lng        float64   `bson:"lng"`
	RecordedAt time.Time `bson:"recorded_at"`
}

func fromEntity(p entity.Position) model {
	return model{
		ID:         p.ID,
		TripID:     p.TripID,
		PartyID:    p.PartyID,
		UserID:     p.UserID,
		Lat:        p.Lat,
		Lng:        p.Lng,
		RecordedAt: p.RecordedAt,
	}
}

func (m model) toEntity() entity.Position {
	return entity.Position{
		ID:         m.ID,
		TripID:     m.TripID,
		PartyID:    m.PartyID,
		UserID:     m.UserID,
		Lat:        m.Lat,
		Lng:        m.Lng,
		RecordedAt: m.RecordedAt,
	}
}
