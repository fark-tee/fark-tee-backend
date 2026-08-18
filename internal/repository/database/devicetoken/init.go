package devicetoken

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ErrNotFound is returned when no device token matches the given lookup.
var ErrNotFound = errors.New("devicetoken: not found")

type Repository interface {
	// Upsert stores deviceToken, keyed by its Token. If the token already
	// exists (e.g. the same physical device logged in as a different
	// account), its UserID and Platform are updated in place rather than
	// creating a duplicate row.
	Upsert(ctx context.Context, deviceToken entity.DeviceToken) error
	// FindByUserID returns every device token registered for userID (a user
	// may have multiple devices).
	FindByUserID(ctx context.Context, userID string) ([]entity.DeviceToken, error)
	// DeleteByToken removes a single token, e.g. on logout.
	DeleteByToken(ctx context.Context, token string) error
}

type repositoryImpl struct {
	collection *mongo.Collection
}

// @WireSet("Repository")
func New(ctx context.Context, db *mongo.Database) (Repository, error) {
	collection := db.Collection("device_tokens")

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "token", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	})
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{collection: collection}, nil
}
