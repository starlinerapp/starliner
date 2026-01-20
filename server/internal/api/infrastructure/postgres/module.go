package postgres

import (
	"database/sql"
	"go.uber.org/fx"
	"starliner.app/internal/api/infrastructure/postgres/sqlc"
)

func NewQueries(db *sql.DB) *sqlc.Queries {
	return sqlc.New(db)
}

var Module = fx.Module(
	"postgres",
	fx.Provide(
		Connect,
		NewQueries,
	),
)
