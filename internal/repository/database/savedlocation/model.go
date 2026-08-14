package savedlocation

import "github.com/fark-tee/fark-tee-backend/internal/entity"

// model is the MongoDB document shape for the saved_locations collection. It
// stays private to this package - callers only ever see entity.SavedLocation.
type model struct {
	ID     string  `bson:"_id"`
	UserID string  `bson:"user_id"`
	Name   string  `bson:"name"`
	Lat    float64 `bson:"lat"`
	Lng    float64 `bson:"lng"`
}

func fromEntity(l entity.SavedLocation) model {
	return model{
		ID:     l.ID,
		UserID: l.UserID,
		Name:   l.Name,
		Lat:    l.Lat,
		Lng:    l.Lng,
	}
}

func (m model) toEntity() entity.SavedLocation {
	return entity.SavedLocation{
		ID:     m.ID,
		UserID: m.UserID,
		Name:   m.Name,
		Lat:    m.Lat,
		Lng:    m.Lng,
	}
}
