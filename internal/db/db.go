package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func InitDB(pgURL string) (*sqlx.DB, error) {
	log.Info().Msg("Database initializing")
	db, err := sqlx.Connect("pgx", pgURL)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("database initialized successfully")
	return db, nil
}
