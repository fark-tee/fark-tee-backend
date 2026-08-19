package partymember

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ErrNotFound is returned when no party member matches the given lookup.
var ErrNotFound = errors.New("partymember: not found")

type Repository interface {
	Create(ctx context.Context, member entity.PartyMember) (entity.PartyMember, error)
	FindByPartyIDAndUserID(ctx context.Context, partyID, userID string) (entity.PartyMember, error)
	FindPendingByUserID(ctx context.Context, userID string) ([]entity.PartyMember, error)
	FindAcceptedByUserID(ctx context.Context, userID string) ([]entity.PartyMember, error)
	FindByPartyID(ctx context.Context, partyID string) ([]entity.PartyMember, error)
	UpdateStatus(ctx context.Context, id string, status entity.PartyMemberStatus) (entity.PartyMember, error)
	UpdateTripStatus(ctx context.Context, id string, tripStatus entity.TripStatus) (entity.PartyMember, error)
	UpdateCheckIn(ctx context.Context, id string, status entity.CheckInStatus, requestedByUserID string) (entity.PartyMember, error)
	Delete(ctx context.Context, id string) error
}

type repositoryImpl struct {
	collection *mongo.Collection
}

// @WireSet("Repository")
func New(ctx context.Context, db *mongo.Database) (Repository, error) {
	collection := db.Collection("party_members")

	_, err := collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "party_id", Value: 1}, {Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "status", Value: 1}},
	})
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{collection: collection}, nil
}
