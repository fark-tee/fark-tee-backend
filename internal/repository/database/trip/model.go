package trip

import (
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// model is the MongoDB document shape for the trips collection. It stays
// private to this package - callers only ever see entity.Trip.
type model struct {
	ID              bson.ObjectID `bson:"_id"`
	PartyID         string        `bson:"party_id"`
	UserID          string        `bson:"user_id"`
	Direction       string        `bson:"direction"`
	DestinationName string        `bson:"destination_name"`
	DestinationLat  float64       `bson:"destination_lat"`
	DestinationLng  float64       `bson:"destination_lng"`
	StartedAt       time.Time     `bson:"started_at"`
	Polyline        string        `bson:"polyline"`
}

func fromEntity(t entity.Trip) (model, error) {
	id, err := mongoid.ToObjectID(t.ID)
	if err != nil {
		return model{}, err
	}

	return model{
		ID:              id,
		PartyID:         t.PartyID,
		UserID:          t.UserID,
		Direction:       string(t.Direction),
		DestinationName: t.Destination.Name,
		DestinationLat:  t.Destination.Lat,
		DestinationLng:  t.Destination.Lng,
		StartedAt:       t.StartedAt,
		Polyline:        t.Polyline,
	}, nil
}

func (m model) toEntity() entity.Trip {
	return entity.Trip{
		ID:        mongoid.FromObjectID(m.ID),
		PartyID:   m.PartyID,
		UserID:    m.UserID,
		Direction: entity.TripDirection(m.Direction),
		Destination: entity.Destination{
			Name: m.DestinationName,
			Lat:  m.DestinationLat,
			Lng:  m.DestinationLng,
		},
		StartedAt: m.StartedAt,
		Polyline:  m.Polyline,
	}
}
