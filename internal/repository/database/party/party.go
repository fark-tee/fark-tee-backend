package party

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *repositoryImpl) Create(ctx context.Context, party entity.Party) (entity.Party, error) {
	doc := fromEntity(party)

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.Party{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByID(ctx context.Context, id string) (entity.Party, error) {
	var doc model

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Party{}, ErrNotFound
		}

		return entity.Party{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]entity.Party, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}

	var docs []model
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	parties := make([]entity.Party, 0, len(docs))
	for _, doc := range docs {
		parties = append(parties, doc.toEntity())
	}

	return parties, nil
}
