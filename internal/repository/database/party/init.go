package party

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ErrNotFound is returned when no party matches the given lookup.
var ErrNotFound = errors.New("party: not found")

type Repository interface {
	Create(ctx context.Context, party entity.Party) (entity.Party, error)
	FindByID(ctx context.Context, id string) (entity.Party, error)
	FindByIDs(ctx context.Context, ids []string) ([]entity.Party, error)
}

type repositoryImpl struct {
	collection *mongo.Collection
}

// @WireSet("Repository")
func New(_ context.Context, db *mongo.Database) (Repository, error) {
	return &repositoryImpl{collection: db.Collection("parties")}, nil
}
