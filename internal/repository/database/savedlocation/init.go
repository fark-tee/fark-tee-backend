package savedlocation

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ErrNotFound is returned when no saved location matches the given lookup.
var ErrNotFound = errors.New("savedlocation: not found")

type Repository interface {
	Create(ctx context.Context, location entity.SavedLocation) (entity.SavedLocation, error)
	FindByID(ctx context.Context, id string) (entity.SavedLocation, error)
	FindByUserID(ctx context.Context, userID string) ([]entity.SavedLocation, error)
	Update(ctx context.Context, location entity.SavedLocation) (entity.SavedLocation, error)
	Delete(ctx context.Context, id string) error
}

type repositoryImpl struct {
	collection *mongo.Collection
}

// @WireSet("Repository")
func New(ctx context.Context, db *mongo.Database) (Repository, error) {
	collection := db.Collection("saved_locations")

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{collection: collection}, nil
}
