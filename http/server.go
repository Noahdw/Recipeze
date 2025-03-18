// Package http has the [Server] and HTTP handlers.
package http

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

// Server holds dependencies for the HTTP server as well as the HTTP server itself.
type Server struct {
	db     *repo.Queries
	log    *slog.Logger
	mux    chi.Router
	server *http.Server
}

type NewServerOptions struct {
	DB  *pgxpool.Pool
	Log *slog.Logger
}

func NewServer(opts NewServerOptions) *Server {
	if opts.Log == nil {
		opts.Log = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	mux := chi.NewMux()

	return &Server{
		db:  repo.New(opts.DB),
		log: opts.Log,
		mux: mux,
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
func (s *Server) Start() error {
	s.log.Info("Starting http server", "address", "0.0.0.0:8080")

	// Important - maps paths to handlers
	s.setupRoutes()

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Stop the server gracefully.
func (s *Server) Stop() error {
	s.log.Info("Stopping http server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	s.log.Info("Stopped http server")
	return nil
}
