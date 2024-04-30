package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"

	"github.com/Uikola/carCatalog/internal/db"
	"github.com/Uikola/carCatalog/internal/db/repository/postgres"
	"github.com/Uikola/carCatalog/internal/server"
	"github.com/Uikola/carCatalog/internal/server/car"
	"github.com/Uikola/carCatalog/internal/usecase/car_usecase"
	"github.com/Uikola/carCatalog/pkg/zlog"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//	@title			Car Catalog API
//	@version		1.0
//	@description	This is an API for managing the car catalog.

//	@contact.name	API Support
//	@contact.url	https://t.me/uikola
//	@contact.email	ugulaev806@yandex.ru

// @host		localhost:8000
// @BasePath	/api
func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

	logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
	log.Logger = zlog.Default(true, "dev", zerolog.Level(logLevel))

	if err = run(ctx); err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	database, err := db.InitDB(os.Getenv("POSTGRES_CONN"))
	if err != nil {
		return err
	}

	m := createMigrate()
	if err = m.Up(); err != nil {
		log.Error().Err(err).Msg("failed to create users table")
	}

	carRepository := postgres.NewCarRepository(database)
	carUseCase := car_usecase.NewUseCaseImpl(carRepository)
	carHandler := car.NewHandler(carUseCase)

	srv := server.NewServer(carHandler)
	httpServer := &http.Server{
		Addr:    os.Getenv("HTTP_PORT"),
		Handler: srv,
	}

	go func() {
		log.Info().Msgf("listening http on %s", httpServer.Addr)
		if err = httpServer.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err = httpServer.Shutdown(ctx); err != nil {
			log.Error().Msg(err.Error())
			os.Exit(1)
		}
	}()
	wg.Wait()

	return nil
}

func createMigrate() *migrate.Migrate {
	m, err := migrate.New(
		"file://migrations/",
		os.Getenv("PGX_URL"),
	)

	if err != nil {
		log.Error().Msg(err.Error())
	}

	return m
}
