package review

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ErrNotFound is returned when no review matches the given lookup.
var ErrNotFound = errors.New("review: not found")

type Repository interface {
	Create(ctx context.Context, review entity.Review) (entity.Review, error)
	// FindByPartyIDReviewerIDAndTargetUserID looks up the review reviewerID
	// has already left for targetUserID within partyID, if any.
	FindByPartyIDReviewerIDAndTargetUserID(ctx context.Context, partyID, reviewerID, targetUserID string) (entity.Review, error)
	// FindByPartyIDAndReviewerID returns every review reviewerID has left
	// within partyID.
	FindByPartyIDAndReviewerID(ctx context.Context, partyID, reviewerID string) ([]entity.Review, error)
}

type repositoryImpl struct {
	collection *mongo.Collection
}

// @WireSet("Repository")
func New(ctx context.Context, db *mongo.Database) (Repository, error) {
	collection := db.Collection("reviews")

	// A reviewer may only leave one review per target per party.
	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "party_id", Value: 1},
			{Key: "reviewer_id", Value: 1},
			{Key: "target_user_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{collection: collection}, nil
}
