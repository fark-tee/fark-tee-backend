package story

import (
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// model is the MongoDB document shape for the stories collection. It stays
// private to this package - callers only ever see entity.Story.
type model struct {
	ID        bson.ObjectID `bson:"_id"`
	PartyID   string        `bson:"party_id"`
	UserID    string        `bson:"user_id"`
	Image     string        `bson:"image"`
	CreatedAt time.Time     `bson:"created_at"`
}

func fromEntity(s entity.Story) (model, error) {
	id, err := mongoid.ToObjectID(s.ID)
	if err != nil {
		return model{}, err
	}

	return model{
		ID:        id,
		PartyID:   s.PartyID,
		UserID:    s.UserID,
		Image:     s.Image,
		CreatedAt: s.CreatedAt,
	}, nil
}

func (m model) toEntity() entity.Story {
	return entity.Story{
		ID:        mongoid.FromObjectID(m.ID),
		PartyID:   m.PartyID,
		UserID:    m.UserID,
		Image:     m.Image,
		CreatedAt: m.CreatedAt,
	}
}
