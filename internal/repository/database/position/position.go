package position

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *repositoryImpl) Create(ctx context.Context, position entity.Position) (entity.Position, error) {
	doc := fromEntity(position)

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.Position{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindLatestByPartyIDAndUserID(ctx context.Context, partyID, userID string) (entity.Position, error) {
	var doc model

	opts := options.FindOne().SetSort(bson.D{{Key: "recorded_at", Value: -1}})

	err := r.collection.FindOne(ctx, bson.M{"party_id": partyID, "user_id": userID}, opts).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Position{}, ErrNotFound
		}

		return entity.Position{}, err
	}

	return doc.toEntity(), nil
}

// FindLatestByPartyID returns the most recently recorded position for every
// user that has one within the given party.
func (r *repositoryImpl) FindLatestByPartyID(ctx context.Context, partyID string) ([]entity.Position, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{Key: "party_id", Value: partyID}}}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "recorded_at", Value: -1}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "doc", Value: bson.D{{Key: "$first", Value: "$$ROOT"}}},
		}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$doc"}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var docs []model
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	positions := make([]entity.Position, 0, len(docs))
	for _, doc := range docs {
		positions = append(positions, doc.toEntity())
	}

	return positions, nil
}
