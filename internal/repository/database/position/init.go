package position

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ErrNotFound is returned when no position matches the given lookup.
var ErrNotFound = errors.New("position: not found")

type Repository interface {
	Create(ctx context.Context, position entity.Position) (entity.Position, error)
	FindLatestByPartyIDAndUserID(ctx context.Context, partyID, userID string) (entity.Position, error)
	FindLatestByPartyID(ctx context.Context, partyID string) ([]entity.Position, error)
}

type repositoryImpl struct {
	collection *mongo.Collection
}

// @WireSet("Repository")
func New(ctx context.Context, db *mongo.Database) (Repository, error) {
	collection := db.Collection("positions")

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "party_id", Value: 1},
			{Key: "user_id", Value: 1},
			{Key: "recorded_at", Value: -1},
		},
	})
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{collection: collection}, nil
}
