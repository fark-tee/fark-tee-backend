package position

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// idFilter builds a match condition for an ID-reference field (party_id,
// user_id, ...) that matches regardless of whether that field was stored as
// this codebase's hex-string convention or, from a bad write such as a
// manual import, as a raw ObjectID.
func idFilter(field, id string) bson.M {
	objID, err := mongoid.ToObjectID(id)
	if err != nil {
		return bson.M{field: id}
	}

	return bson.M{"$or": bson.A{
		bson.M{field: id},
		bson.M{field: objID},
	}}
}

func (r *repositoryImpl) Create(ctx context.Context, position entity.Position) (entity.Position, error) {
	doc, err := fromEntity(position)
	if err != nil {
		return entity.Position{}, err
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.Position{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindLatestByPartyIDAndUserID(ctx context.Context, partyID, userID string) (entity.Position, error) {
	var doc model

	opts := options.FindOne().SetSort(bson.D{{Key: "recorded_at", Value: -1}})

	filter := bson.M{"$and": bson.A{idFilter("party_id", partyID), idFilter("user_id", userID)}}

	err := r.collection.FindOne(ctx, filter, opts).Decode(&doc)
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
		bson.D{{Key: "$match", Value: idFilter("party_id", partyID)}},
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
