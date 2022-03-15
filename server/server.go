package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.autoiterative.com/group-zealous-ishizaka-gates/backendify/logger"
	"go.uber.org/zap"
)

// Server represents the main configuration for the http server.
type Server struct {
	*http.Server
	*chi.Mux
	logger *logger.Logger
}

// logRoutes is used by Zap Logger to register all the routes that the API has.
func (s *Server) logRoutes() {
	if err := chi.Walk(s, s.printRouteInZap()); err != nil {
		s.logger.Error("Failed to walk routes:", zap.Error(err))
	}
}

// printRouteInZap creates the corresponding format to print the Routes that the API has.
func (s *Server) printRouteInZap() chi.WalkFunc {
	return func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")

		s.logger.Debug("Route registered", zap.String("method", method), zap.String("route", route))

		return nil
	}
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (s *Server) Start() {
	s.logger.Info("Starting server...")

	defer func() {
		if err := s.logger.Sync(); err != nil {
			s.logger.Warn(fmt.Sprintf("logger.Sync(): %s", err))
		}
	}()

	// logging current routes
	s.logRoutes()

	go func() {
		if err := s.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal("Could not listen on", zap.String("addr", s.Addr), zap.Error(err))
		}
	}()

	s.logger.Info("Server is ready to handle requests", zap.String("addr", s.Addr))

	s.gracefulShutdown()
}

// gracefulShutdown is used to shut down in a graceful way.
func (s *Server) gracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	sig := <-quit

	s.logger.Info("Server is shutting down", zap.String("reason", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.SetKeepAlivesEnabled(false)

	if err := s.Shutdown(ctx); err != nil {
		s.logger.Fatal("Could not gracefully shutdown the server", zap.Error(err))
	}

	s.logger.Info("Server stopped")
}

// WithOptions defines the possible options which could be passed to the
// server to set extra features.
func (s *Server) WithOptions(opts ...Option) {
	// some options needs to be set first than others.
	sort.Slice(opts, func(i, j int) bool {
		return opts[i].Key() < opts[j].Key()
	})

	for i := range opts {
		opts[i].Apply(s)
	}
}

// NewServer creates and configures a server serving all application routes.
// The server implements a graceful shutdown and utilizes zap.Logger for logging purposes.
// chi.Mux is used for registering some convenient middlewares and easy configuration of
// routes using different http verbs.
func New(opts ...Option) *Server {
	var (
		router = chi.NewRouter()
	)

	s := &Server{
		&http.Server{
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
		router,
		logger.NewLogger(logger.DefaultEnvLoggerConfig()),
	}

	// registered the first middleware as a required to log everything
	router.Use(logger.ChiZapLoggerMiddleware(s.logger))

	// applying the options
	s.WithOptions(opts...)

	return s
}
