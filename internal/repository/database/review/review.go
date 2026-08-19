package review

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *repositoryImpl) Create(ctx context.Context, review entity.Review) (entity.Review, error) {
	doc, err := fromEntity(review)
	if err != nil {
		return entity.Review{}, err
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.Review{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByPartyIDReviewerIDAndTargetUserID(ctx context.Context, partyID, reviewerID, targetUserID string) (entity.Review, error) {
	var doc model

	filter := bson.M{"party_id": partyID, "reviewer_id": reviewerID, "target_user_id": targetUserID}

	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.Review{}, ErrNotFound
		}

		return entity.Review{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByPartyIDAndReviewerID(ctx context.Context, partyID, reviewerID string) ([]entity.Review, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"party_id": partyID, "reviewer_id": reviewerID})
	if err != nil {
		return nil, err
	}

	var docs []model
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	reviews := make([]entity.Review, 0, len(docs))
	for _, doc := range docs {
		reviews = append(reviews, doc.toEntity())
	}

	return reviews, nil
}
