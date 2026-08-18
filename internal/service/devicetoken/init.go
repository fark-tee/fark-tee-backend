package devicetoken

import (
	"context"

	"github.com/fark-tee/fark-tee-backend/internal/repository/database/devicetoken"
)

type Service interface {
	// Register upserts a device token for userID, creating it (or moving it
	// from a previous owner) if it doesn't already point at userID.
	Register(ctx context.Context, userID, token, platform string) error
	// Unregister removes token, e.g. on logout.
	Unregister(ctx context.Context, token string) error
}

type serviceImpl struct {
	repo devicetoken.Repository
}

// @WireSet("Service")
func New(repo devicetoken.Repository) Service {
	return &serviceImpl{repo: repo}
}
