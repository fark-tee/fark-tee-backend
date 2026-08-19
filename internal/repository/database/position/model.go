package position

import (
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// model is the MongoDB document shape for the positions collection. It stays
// private to this package - callers only ever see entity.Position.
type model struct {
	ID                       bson.ObjectID `bson:"_id"`
	TripID                   string        `bson:"trip_id"`
	PartyID                  string        `bson:"party_id"`
	UserID                   string        `bson:"user_id"`
	Lat                      float64       `bson:"lat"`
	Lng                      float64       `bson:"lng"`
	RecordedAt               time.Time     `bson:"recorded_at"`
	EstimatedDurationSeconds int           `bson:"estimated_duration_seconds"`
	EstimatedArrivalAt       time.Time     `bson:"estimated_arrival_at"`
}

func fromEntity(p entity.Position) (model, error) {
	id, err := mongoid.ToObjectID(p.ID)
	if err != nil {
		return model{}, err
	}

	return model{
		ID:                       id,
		TripID:                   p.TripID,
		PartyID:                  p.PartyID,
		UserID:                   p.UserID,
		Lat:                      p.Lat,
		Lng:                      p.Lng,
		RecordedAt:               p.RecordedAt,
		EstimatedDurationSeconds: p.EstimatedDurationSeconds,
		EstimatedArrivalAt:       p.EstimatedArrivalAt,
	}, nil
}

func (m model) toEntity() entity.Position {
	return entity.Position{
		ID:                       mongoid.FromObjectID(m.ID),
		TripID:                   m.TripID,
		PartyID:                  m.PartyID,
		UserID:                   m.UserID,
		Lat:                      m.Lat,
		Lng:                      m.Lng,
		RecordedAt:               m.RecordedAt,
		EstimatedDurationSeconds: m.EstimatedDurationSeconds,
		EstimatedArrivalAt:       m.EstimatedArrivalAt,
	}
}
