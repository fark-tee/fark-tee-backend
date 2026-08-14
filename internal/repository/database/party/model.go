package party

import (
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
)

// model is the MongoDB document shape for the parties collection. It stays
// private to this package - callers only ever see entity.Party.
type model struct {
	ID              string    `bson:"_id"`
	Name            string    `bson:"name"`
	DestinationName string    `bson:"destination_name"`
	DestinationLat  float64   `bson:"destination_lat"`
	DestinationLng  float64   `bson:"destination_lng"`
	TargetTime      time.Time `bson:"target_time"`
	CreatedByID     string    `bson:"created_by_id"`
	CreatedByName   string    `bson:"created_by_name"`
}

func fromEntity(p entity.Party) model {
	return model{
		ID:              p.ID,
		Name:            p.Name,
		DestinationName: p.DestinationName,
		DestinationLat:  p.DestinationLat,
		DestinationLng:  p.DestinationLng,
		TargetTime:      p.TargetTime,
		CreatedByID:     p.CreatedByID,
		CreatedByName:   p.CreatedByName,
	}
}

func (m model) toEntity() entity.Party {
	return entity.Party{
		ID:              m.ID,
		Name:            m.Name,
		DestinationName: m.DestinationName,
		DestinationLat:  m.DestinationLat,
		DestinationLng:  m.DestinationLng,
		TargetTime:      m.TargetTime,
		CreatedByID:     m.CreatedByID,
		CreatedByName:   m.CreatedByName,
	}
}
