//go:build wireinject

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

		service.NewUserService,
		service.NewOrganizationService,
		service.NewProjectService,
		service.NewEnvironmentService,

		repository.NewUserRepository,
		repository.NewOrganizationRepository,
		repository.NewProjectRepository,
		repository.NewEnvironmentRepository,

		middleware.NewBasicAuthMiddleware,
		middleware.NewUserMiddleware,

		handler.NewRootHandler,
		handler.NewUserHandler,
		handler.NewOrganizationHandler,
		handler.NewProjectHandler,
		handler.NewEnvironmentHandler,

		http.NewServer,
	)

	return &http.Server{}, nil
}
