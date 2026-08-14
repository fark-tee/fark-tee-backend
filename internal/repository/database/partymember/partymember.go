package partymember

import (
	"context"
	"errors"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *repositoryImpl) Create(ctx context.Context, member entity.PartyMember) (entity.PartyMember, error) {
	doc := fromEntity(member)

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.PartyMember{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByPartyIDAndUserID(ctx context.Context, partyID, userID string) (entity.PartyMember, error) {
	var doc model

	err := r.collection.FindOne(ctx, bson.M{"party_id": partyID, "user_id": userID}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.PartyMember{}, ErrNotFound
		}

		return entity.PartyMember{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindPendingByUserID(ctx context.Context, userID string) ([]entity.PartyMember, error) {
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id": userID,
		"status":  string(entity.PartyMemberStatusPending),
	})
	if err != nil {
		return nil, err
	}

	var docs []model
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	members := make([]entity.PartyMember, 0, len(docs))
	for _, doc := range docs {
		members = append(members, doc.toEntity())
	}

	return members, nil
}

func (r *repositoryImpl) UpdateStatus(ctx context.Context, id string, status entity.PartyMemberStatus) (entity.PartyMember, error) {
	var doc model

	err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": string(status)}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.PartyMember{}, ErrNotFound
		}

		return entity.PartyMember{}, err
	}

	return doc.toEntity(), nil
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
