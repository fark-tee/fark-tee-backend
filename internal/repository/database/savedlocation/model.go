package savedlocation

import (
	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// model is the MongoDB document shape for the saved_locations collection. It
// stays private to this package - callers only ever see entity.SavedLocation.
type model struct {
	ID     bson.ObjectID `bson:"_id"`
	UserID string        `bson:"user_id"`
	Name   string        `bson:"name"`
	Lat    float64       `bson:"lat"`
	Lng    float64       `bson:"lng"`
}

func fromEntity(l entity.SavedLocation) (model, error) {
	id, err := mongoid.ToObjectID(l.ID)
	if err != nil {
		return model{}, err
	}

	return model{
		ID:     id,
		UserID: l.UserID,
		Name:   l.Name,
		Lat:    l.Lat,
		Lng:    l.Lng,
	}, nil
}

func (m model) toEntity() entity.SavedLocation {
	return entity.SavedLocation{
		ID:     mongoid.FromObjectID(m.ID),
		UserID: m.UserID,
		Name:   m.Name,
		Lat:    m.Lat,
		Lng:    m.Lng,
	}
}
