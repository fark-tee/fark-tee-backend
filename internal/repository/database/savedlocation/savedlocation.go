package savedlocation

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *repositoryImpl) Create(ctx context.Context, location entity.SavedLocation) (entity.SavedLocation, error) {
	doc, err := fromEntity(location)
	if err != nil {
		return entity.SavedLocation{}, err
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.SavedLocation{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByID(ctx context.Context, id string) (entity.SavedLocation, error) {
	var doc model

	objID, err := mongoid.ToObjectID(id)
	if err != nil {
		return entity.SavedLocation{}, ErrNotFound
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.SavedLocation{}, ErrNotFound
		}

		return entity.SavedLocation{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByUserID(ctx context.Context, userID string) ([]entity.SavedLocation, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	var docs []model
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	locations := make([]entity.SavedLocation, 0, len(docs))
	for _, doc := range docs {
		locations = append(locations, doc.toEntity())
	}

	return locations, nil
}

func (r *repositoryImpl) Update(ctx context.Context, location entity.SavedLocation) (entity.SavedLocation, error) {
	doc, err := fromEntity(location)
	if err != nil {
		return entity.SavedLocation{}, err
	}

	err = r.collection.FindOneAndReplace(
		ctx,
		bson.M{"_id": doc.ID},
		doc,
		options.FindOneAndReplace().SetReturnDocument(options.After),
	).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.SavedLocation{}, ErrNotFound
		}

		return entity.SavedLocation{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) Delete(ctx context.Context, id string) error {
	objID, err := mongoid.ToObjectID(id)
	if err != nil {
		return ErrNotFound
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
