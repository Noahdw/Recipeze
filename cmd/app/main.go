package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
	"maragu.dev/env"

	"recipeze/appconfig"
	"recipeze/server"
)

func main() {
	// Set up a logger that is used throughout the app
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))

	if err := start(log); err != nil {
		log.Error("Error starting app", "error", err)
		os.Exit(1)
	}
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func start(log *slog.Logger) error {
	log.Info("Starting app")

	// We load environment variables from .env if it exists
	_ = env.Load()

	// Catch signals to gracefully shut down the app
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	appconfig.Initialize()

	dbConfig := DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		slog.Error("Unable to create connection pool",
			"error", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	slog.Info("Connected to DB")

	// Set up the HTTP server, injecting the database and logger
	s := server.NewServer(server.NewServerOptions{
		DB:  dbpool,
		Log: log,
	})

	// Use an errgroup to wait for separate goroutines which can error
	eg, ctx := errgroup.WithContext(ctx)

	// Start the server within the errgroup.
	// You can do this for other dependencies as well.
	eg.Go(func() error {
		return s.Start()
	})

	// Wait for the context to be done, which happens when a signal is caught
	<-ctx.Done()
	log.Info("Stopping app")

	// Stop the server gracefully
	eg.Go(func() error {
		return s.Stop()
	})

	// Wait for the server to stop
	if err := eg.Wait(); err != nil {
		return err
	}

	log.Info("Stopped app")

	return nil
}
