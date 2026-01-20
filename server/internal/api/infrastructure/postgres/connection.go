package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"starliner.app/internal/api/conf"
)

func Connect(cfg *conf.Config) (*sql.DB, error) {
	var connString = fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Check if the connection is successful
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
