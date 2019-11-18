package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

// Server is the http server
type Server struct {
	logger *zerolog.Logger
	server *http.Server
	host   string
}

// New creates a new http server
func New(logger *zerolog.Logger, host string) *Server {
	l := logger.With().Str("pkg", "http-server").Logger()
	return &Server{
		logger: &l,
		host:   host,
	}
}

// Start starts accepting http requests
func (s *Server) Start(routes func(router chi.Router), middlewares ...func(http.Handler) http.Handler) {
	// create http router
	router := chi.NewRouter()

	// setup middlewares
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Recoverer)
	for _, m := range middlewares {
		router.Use(m)
	}

	// setup http routes
	routes(router)

	// create server
	s.server = &http.Server{
		Addr:    fmt.Sprintf(s.host),
		Handler: router,
	}

	// start http server
	s.logger.Info().Str("host", s.host).Msg("starting http server...")
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		s.logger.Fatal().Err(err).Str("host", s.host).Msg("failed to start http server")
		return
	}
	go s.server.Serve(listener)
}

// Done terminates the http server
func (s *Server) Done() {
	s.logger.Info().Msg("stopping http server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	s.server.Shutdown(ctx)
	s.logger.Info().Msg("stopped http server")
}
