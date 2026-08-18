package trip

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *repositoryImpl) Create(ctx context.Context, trip entity.Trip) (entity.Trip, error) {
	doc, err := fromEntity(trip)
	if err != nil {
		return entity.Trip{}, err
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.Trip{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindLatestByPartyIDAndUserID(ctx context.Context, partyID, userID string) (entity.Trip, error) {
	var doc model

	opts := options.FindOne().SetSort(bson.D{{Key: "started_at", Value: -1}})

	err := r.collection.FindOne(ctx, bson.M{"party_id": partyID, "user_id": userID}, opts).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Trip{}, ErrNotFound
		}

		return entity.Trip{}, err
	}

	return doc.toEntity(), nil
}
