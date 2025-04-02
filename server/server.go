// Package http has the [server] and HTTP handlers.
package server

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"recipeze/repo"
)

// server holds dependencies for the HTTP server as well as the HTTP server itself.
type server struct {
	queries *repo.Queries
	db      *pgxpool.Pool
	log     *slog.Logger
	mux     chi.Router
	server  *http.Server
}

type NewServerOptions struct {
	DB  *pgxpool.Pool
	Log *slog.Logger
}

func NewServer(opts NewServerOptions) *server {
	if opts.Log == nil {
		opts.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	mux := chi.NewMux()

	return &server{
		queries: repo.New(opts.DB),
		db:      opts.DB,
		log:     opts.Log,
		mux:     mux,
		server: &http.Server{
			Addr:              ":8080",
			Handler:           mux,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       5 * time.Second,
		},
	}
}

// Start the server and set up routes.
func (s *server) Start() error {
	s.log.Info("Starting http server", "address", "0.0.0.0:8080")

	// Important - maps paths to handlers
	s.SetupRoutes()

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop the server gracefully.
func (s *server) Stop() error {
	s.log.Info("Stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.log.Info("Stopped http server")
	return nil
}
