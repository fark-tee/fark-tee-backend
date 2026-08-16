package story

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ErrNotFound is returned when no story matches the given lookup.
var ErrNotFound = errors.New("story: not found")

type Repository interface {
	Create(ctx context.Context, story entity.Story) (entity.Story, error)
	FindByID(ctx context.Context, id string) (entity.Story, error)
	FindByPartyIDAndUserID(ctx context.Context, partyID, userID string) ([]entity.Story, error)
	Delete(ctx context.Context, id string) error
}

type repositoryImpl struct {
	collection *mongo.Collection
}

// @WireSet("Repository")
func New(ctx context.Context, db *mongo.Database) (Repository, error) {
	collection := db.Collection("stories")

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "party_id", Value: 1}, {Key: "user_id", Value: 1}},
	})
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{collection: collection}, nil
}
