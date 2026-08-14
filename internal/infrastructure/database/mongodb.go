package database

import (
	"context"
	"log/slog"

	"github.com/fark-tee/fark-tee-backend/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// @WireSet("Infrastructure")
func NewMongoClient(ctx context.Context, cfg *config.Config) (*mongo.Client, func(), error) {
	client, err := mongo.Connect(options.Client().ApplyURI(cfg.Database.URI))
	if err != nil {
		return nil, nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)

		return nil, nil, err
	}

	slog.Info("Connected to MongoDB", slog.String("database", cfg.Database.Database))

	return client, func() { _ = client.Disconnect(ctx) }, nil
}

// @WireSet("Infrastructure")
func NewMongoDatabase(client *mongo.Client, cfg *config.Config) *mongo.Database {
	return client.Database(cfg.Database.Database)
}
