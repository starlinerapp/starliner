//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"
	"github.com/google/wire"
	http "starliner.app/pkg/api"
	"starliner.app/pkg/api/handler"
	"starliner.app/pkg/api/middleware"
	"starliner.app/pkg/config"
	"starliner.app/pkg/db"
	"starliner.app/pkg/db/sqlc"
	"starliner.app/pkg/repository"
	"starliner.app/pkg/service"
)

func ProvideQueries(db *sql.DB) *sqlc.Queries {
	return sqlc.New(db)
}

func InitializeAPI(cfg *config.Config) (*http.Server, error) {
	wire.Build(
		db.Connect,
		ProvideQueries,
		repository.NewUserRepository,
		service.NewUserService,
		middleware.NewBasicAuthMiddleware,
		middleware.NewUserMiddleware,
		http.NewServer,
		handler.NewRootHandler,
		handler.NewUserHandler,
	)

	return &http.Server{}, nil
}
