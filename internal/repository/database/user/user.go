package user

import (
	"context"
	"errors"
	"regexp"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// searchResultLimit caps how many users a single Search call can return,
// since there is no pagination convention yet for this endpoint.
const searchResultLimit = 20

func (r *repositoryImpl) Create(ctx context.Context, user entity.User) (entity.User, error) {
	doc, err := fromEntity(user)
	if err != nil {
		return entity.User{}, err
	}

	if _, err := r.collection.InsertOne(ctx, doc); err != nil {
		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByID(ctx context.Context, id string) (entity.User, error) {
	var doc model

	objID, err := mongoid.ToObjectID(id)
	if err != nil {
		return entity.User{}, ErrNotFound
	}

	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&doc)
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

func (r *repositoryImpl) UpdateProfile(ctx context.Context, id, displayName, username, emergencyContactName, emergencyContactPhone string) (entity.User, error) {
	var doc model

	objID, err := mongoid.ToObjectID(id)
	if err != nil {
		return entity.User{}, ErrNotFound
	}

	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{
			"display_name":            displayName,
			"username":                username,
			"emergency_contact_name":  emergencyContactName,
			"emergency_contact_phone": emergencyContactPhone,
		}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, ErrNotFound
		}

		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) UpdateProfileImage(ctx context.Context, id, profileImageURL string) (entity.User, error) {
	var doc model

	objID, err := mongoid.ToObjectID(id)
	if err != nil {
		return entity.User{}, ErrNotFound
	}

	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"profile_image_url": profileImageURL}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, ErrNotFound
		}

		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) IncrementOnTimeCount(ctx context.Context, id string) (entity.User, error) {
	return r.incrementCount(ctx, id, "on_time_count")
}

func (r *repositoryImpl) IncrementLateCount(ctx context.Context, id string) (entity.User, error) {
	return r.incrementCount(ctx, id, "late_count")
}

func (r *repositoryImpl) incrementCount(ctx context.Context, id, field string) (entity.User, error) {
	var doc model

	objID, err := mongoid.ToObjectID(id)
	if err != nil {
		return entity.User{}, ErrNotFound
	}

	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$inc": bson.M{field: 1}},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, ErrNotFound
		}

		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) RecordRating(ctx context.Context, id string, score int) (entity.User, error) {
	var doc model

	objID, err := mongoid.ToObjectID(id)
	if err != nil {
		return entity.User{}, ErrNotFound
	}

	// An aggregation-pipeline update computes the new running average from
	// the document's own current fields in a single atomic operation,
	// avoiding a separate read-then-write race between concurrent reviews.
	update := mongo.Pipeline{
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "rating", Value: bson.D{{Key: "$divide", Value: bson.A{
				bson.D{{Key: "$add", Value: bson.A{
					bson.D{{Key: "$multiply", Value: bson.A{"$rating", "$rating_count"}}},
					score,
				}}},
				bson.D{{Key: "$add", Value: bson.A{"$rating_count", 1}}},
			}}}},
			{Key: "rating_count", Value: bson.D{{Key: "$add", Value: bson.A{"$rating_count", 1}}}},
		}}},
	}

	err = r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": objID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, ErrNotFound
		}

		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) FindByUsername(ctx context.Context, username string) (entity.User, error) {
	var doc model

	filter := bson.M{"username": bson.M{"$regex": "^" + regexp.QuoteMeta(username) + "$", "$options": "i"}}
	err := r.collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, ErrNotFound
		}

		return entity.User{}, err
	}

	return doc.toEntity(), nil
}

func (r *repositoryImpl) Search(ctx context.Context, query string) ([]entity.User, error) {
	clauses := []bson.M{
		{"display_name": bson.M{"$regex": regexp.QuoteMeta(query), "$options": "i"}},
		{"google_user_id": query},
	}

	// query only matches _id when it is a well-formed ObjectID hex string;
	// anything else could never match a document's _id anyway.
	if objID, err := mongoid.ToObjectID(query); err == nil {
		clauses = append(clauses, bson.M{"_id": objID})
	}

	filter := bson.M{"$or": clauses}

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
