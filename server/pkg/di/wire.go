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
		repository.NewOrganizationRepository,
		repository.NewProjectRepository,

		service.NewUserService,
		service.NewOrganizationService,
		service.NewProjectService,

		middleware.NewBasicAuthMiddleware,
		middleware.NewUserMiddleware,

		handler.NewRootHandler,
		handler.NewUserHandler,
		handler.NewOrganizationHandler,
		handler.NewProjectHandler,

		http.NewServer,
	)

	return &http.Server{}, nil
}
