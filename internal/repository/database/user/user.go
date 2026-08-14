package user

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *repositoryImpl) Create(ctx context.Context, user entity.User) (entity.User, error) {
	doc := fromEntity(user)

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByInstagramUserID(ctx context.Context, instagramUserID string) (entity.User, error) {
	var doc model

	err := r.collection.FindOne(ctx, bson.M{"instagram_user_id": instagramUserID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, ErrNotFound
		}

		return entity.User{}, err
	}

	return doc.toEntity(), nil
}
