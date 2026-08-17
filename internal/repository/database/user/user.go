package user

import (
	"context"
	"errors"
	"regexp"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// searchResultLimit caps how many users a single Search call can return,
// since there is no pagination convention yet for this endpoint.
const searchResultLimit = 20

func (r *repositoryImpl) Create(ctx context.Context, user entity.User) (entity.User, error) {
	doc := fromEntity(user)

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByID(ctx context.Context, id string) (entity.User, error) {
	var doc model

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, ErrNotFound
		}

		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByGoogleUserID(ctx context.Context, googleUserID string) (entity.User, error) {
	var doc model

	err := r.collection.FindOne(ctx, bson.M{"google_user_id": googleUserID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, ErrNotFound
		}

		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) Search(ctx context.Context, query string) ([]entity.User, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"display_name": bson.M{"$regex": regexp.QuoteMeta(query), "$options": "i"}},
			{"google_user_id": query},
			{"_id": query},
		},
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().SetLimit(searchResultLimit))
	if err != nil {
		return nil, err
	}

	var docs []model
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	users := make([]entity.User, 0, len(docs))
	for _, doc := range docs {
		users = append(users, doc.toEntity())
	}

	return users, nil
}
