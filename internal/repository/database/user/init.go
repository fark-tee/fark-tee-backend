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
	FindByGoogleUserID(ctx context.Context, googleUserID string) (entity.User, error)
	// FindByUsername looks up a user by username, case-insensitively.
	FindByUsername(ctx context.Context, username string) (entity.User, error)
	// Search returns users whose display name, ID, or Google user ID match
	// query, capped at searchResultLimit results.
	Search(ctx context.Context, query string) ([]entity.User, error)
	// UpdateProfile sets id's display name and username and returns the
	// updated user.
	UpdateProfile(ctx context.Context, id, displayName, username string) (entity.User, error)
	// UpdateProfileImage sets id's profile image URL and returns the updated
	// user.
	UpdateProfileImage(ctx context.Context, id, profileImageURL string) (entity.User, error)
	// IncrementOnTimeCount increments id's on-time arrival count and returns
	// the updated user.
	IncrementOnTimeCount(ctx context.Context, id string) (entity.User, error)
	// IncrementLateCount increments id's late arrival count and returns the
	// updated user.
	IncrementLateCount(ctx context.Context, id string) (entity.User, error)
	// RecordRating folds score into id's running average rating, incrementing
	// its rating count, and returns the updated user.
	RecordRating(ctx context.Context, id string, score int) (entity.User, error)
}

type repositoryImpl struct {
	collection *mongo.Collection
}

// @WireSet("Repository")
func New(ctx context.Context, db *mongo.Database) (Repository, error) {
	collection := db.Collection("users")

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "google_user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	// Partial: only enforced once a user actually has a username, so
	// existing users created before this field existed (empty string) don't
	// collide with each other. $ne isn't supported in partial filter
	// expressions (it compiles to $not), so $gt "" is used instead — any
	// non-empty string is greater than "".
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "username", Value: 1}},
		Options: options.Index().
			SetUnique(true).
			SetPartialFilterExpression(bson.M{"username": bson.M{"$exists": true, "$gt": ""}}),
	})
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{collection: collection}, nil
}
