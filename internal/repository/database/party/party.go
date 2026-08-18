package party

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *repositoryImpl) Create(ctx context.Context, party entity.Party) (entity.Party, error) {
	doc, err := fromEntity(party)
	if err != nil {
		return entity.Party{}, err
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.Party{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByID(ctx context.Context, id string) (entity.Party, error) {
	var doc model

	objID, err := mongoid.ToObjectID(id)
	if err != nil {
		return entity.Party{}, ErrNotFound
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Party{}, ErrNotFound
		}

		return entity.Party{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByIDs(ctx context.Context, ids []string) ([]entity.Party, error) {
	objIDs := make([]bson.ObjectID, 0, len(ids))
	for _, id := range ids {
		// ids that aren't well-formed ObjectID hex strings can't match any
		// document, so they're skipped rather than failing the whole query.
		if objID, err := mongoid.ToObjectID(id); err == nil {
			objIDs = append(objIDs, objID)
		}
	}

	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": objIDs}})
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
