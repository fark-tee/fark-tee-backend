package user

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ErrNotFound is returned when no user matches the given lookup.
var ErrNotFound = errors.New("user: not found")

type Repository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	FindByID(ctx context.Context, id string) (entity.User, error)
	FindByInstagramUserID(ctx context.Context, instagramUserID string) (entity.User, error)
	// Search returns users whose display name, ID, or Instagram user ID
	// match query, capped at searchResultLimit results.
	Search(ctx context.Context, query string) ([]entity.User, error)
}

type repositoryImpl struct {
	collection *mongo.Collection
}

// @WireSet("Repository")
func New(ctx context.Context, db *mongo.Database) (Repository, error) {
	collection := db.Collection("users")

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "instagram_user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{collection: collection}, nil
}
