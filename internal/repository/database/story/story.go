package story

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *repositoryImpl) Create(ctx context.Context, story entity.Story) (entity.Story, error) {
	doc := fromEntity(story)

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.Story{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByID(ctx context.Context, id string) (entity.Story, error) {
	var doc model

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Story{}, ErrNotFound
		}

		return entity.Story{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByPartyIDAndUserID(ctx context.Context, partyID, userID string) ([]entity.Story, error) {
	cursor, err := r.collection.Find(
		ctx,
		bson.M{"party_id": partyID, "user_id": userID},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}

	var docs []model
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	stories := make([]entity.Story, 0, len(docs))
	for _, doc := range docs {
		stories = append(stories, doc.toEntity())
	}

	return stories, nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
