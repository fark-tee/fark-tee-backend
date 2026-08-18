// Package mongoid converts entity IDs (plain strings) to and from the
// bson.ObjectID values MongoDB stores in every collection's _id field.
package mongoid

import "go.mongodb.org/mongo-driver/v2/bson"

// New generates a fresh entity ID, encoded as the hex string of a newly
// generated ObjectID.
func New() string {
	return bson.NewObjectID().Hex()
}

// ToObjectID parses an entity ID string into the ObjectID used for a
// document's _id field.
func ToObjectID(id string) (bson.ObjectID, error) {
	return bson.ObjectIDFromHex(id)
}

// FromObjectID converts a document's _id back into the entity ID string.
func FromObjectID(id bson.ObjectID) string {
	return id.Hex()
}
