//go:build wireinject
// +build wireinject

package di

import (
	"github.com/fark-tee/fark-tee-backend/cmd/server/server"
	"github.com/google/wire"
)

func Initialize() (*server.Server, func(), error) {
	wire.Build(
		ConfigSet,
		InfrastructureSet,
		RepositorySet,
		ServiceSet,
		HandlerSet,
		server.New,
	)

	return nil, nil, nil
}
