package devicetoken

import (
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/entity"
	"github.com/fark-tee/fark-tee-backend/internal/repository/database/mongoid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// model is the MongoDB document shape for the device_tokens collection. It
// stays private to this package - callers only ever see entity.DeviceToken.
type model struct {
	ID        bson.ObjectID `bson:"_id"`
	UserID    string        `bson:"user_id"`
	Token     string        `bson:"token"`
	Platform  string        `bson:"platform"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

func (m model) toEntity() entity.DeviceToken {
	return entity.DeviceToken{
		ID:        mongoid.FromObjectID(m.ID),
		UserID:    m.UserID,
		Token:     m.Token,
		Platform:  entity.DevicePlatform(m.Platform),
		UpdatedAt: m.UpdatedAt,
	}
}
