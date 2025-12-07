package db

import (
	"database/sql"
	"go.uber.org/fx"
	"starliner.app/pkg/db/sqlc"
)

func NewQueries(db *sql.DB) *sqlc.Queries {
	return sqlc.New(db)
}

var Module = fx.Module(
	"db",
	fx.Provide(
		Connect,
		NewQueries,
	),
)
