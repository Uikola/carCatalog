package main

import (
	"context"
	"github.com/Uikola/carCatalog/internal/db"
	"github.com/Uikola/carCatalog/pkg/zlog"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
)

func run(ctx context.Context) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	database, err := db.InitDB(os.Getenv("POSTGRES_CONN"))
	if err != nil {
		return err
	}
	_ = database
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log.Logger = zlog.Default(true, "dev", zerolog.InfoLevel)

	if err := run(ctx); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}
