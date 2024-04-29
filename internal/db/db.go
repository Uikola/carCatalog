package db

import (
	"github.com/jmoiron/sqlx"
)

func InitDB(pgURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", pgURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
