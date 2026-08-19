package review

import (
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// model is the MongoDB document shape for the reviews collection. It stays
// private to this package - callers only ever see entity.Review.
type model struct {
	ID           bson.ObjectID `bson:"_id"`
	PartyID      string        `bson:"party_id"`
	ReviewerID   string        `bson:"reviewer_id"`
	TargetUserID string        `bson:"target_user_id"`
	Score        int           `bson:"score"`
	Comment      string        `bson:"comment"`
	CreatedAt    time.Time     `bson:"created_at"`
}

func fromEntity(r entity.Review) (model, error) {
	id, err := mongoid.ToObjectID(r.ID)
	if err != nil {
		return model{}, err
	}

	return model{
		ID:           id,
		PartyID:      r.PartyID,
		ReviewerID:   r.ReviewerID,
		TargetUserID: r.TargetUserID,
		Score:        r.Score,
		Comment:      r.Comment,
		CreatedAt:    r.CreatedAt,
	}, nil
}

func (m model) toEntity() entity.Review {
	return entity.Review{
		ID:           mongoid.FromObjectID(m.ID),
		PartyID:      m.PartyID,
		ReviewerID:   m.ReviewerID,
		TargetUserID: m.TargetUserID,
		Score:        m.Score,
		Comment:      m.Comment,
		CreatedAt:    m.CreatedAt,
	}
}
