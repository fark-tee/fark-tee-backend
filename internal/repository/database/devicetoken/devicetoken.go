package devicetoken

import (
	"context"
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *repositoryImpl) Upsert(ctx context.Context, deviceToken entity.DeviceToken) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"token": deviceToken.Token},
		bson.M{
			"$set": bson.M{
				"user_id":    deviceToken.UserID,
				"platform":   string(deviceToken.Platform),
				"updated_at": time.Now().UTC(),
			},
			"$setOnInsert": bson.M{
				"_id":   bson.NewObjectID(),
				"token": deviceToken.Token,
			},
		},
		options.UpdateOne().SetUpsert(true),
	)

	return err
}

func (r *repositoryImpl) FindByUserID(ctx context.Context, userID string) ([]entity.DeviceToken, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	var docs []model
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	tokens := make([]entity.DeviceToken, 0, len(docs))
	for _, doc := range docs {
		tokens = append(tokens, doc.toEntity())
	}

	return tokens, nil
}

func (r *repositoryImpl) DeleteByToken(ctx context.Context, token string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"token": token})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	return nil
}
