package server

import (
	"context"
	"errors"
	"net/http"
	"syscall"
	"time"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/handlers"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/SergeyRonzhin/url-shortener/internal/app/middlewares"
	"github.com/SergeyRonzhin/url-shortener/internal/app/migrator"
	"github.com/SergeyRonzhin/url-shortener/internal/app/service"
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/skovtunenko/graterm"
)

type Server struct {
	httpHandlers handlers.HTTPHandler
	options      *config.Options
	logger       *logger.Logger
	middlewares  *middlewares.Middlewares
	term         *graterm.Terminator
}

func New(options *config.Options, logger *logger.Logger) (*Server, context.Context, error) {
	s, err := getStorage(options, logger)

	if err != nil {
		return nil, nil, err
	}

	if _, ok := s.(*storage.DBStorage); ok {
		m := migrator.New(options, logger)

		if err := m.ApplyMigrations(); err != nil {
			panic(err)
		}
	}

	term, ctx := graterm.NewWithSignals(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	term.WithOrder(1).
		WithName("Database").
		Register(500*time.Millisecond, func(ctx context.Context) {
			logger.Info("terminating DB...")
			defer logger.Info("...DB terminated")

			if err := s.Close(); err != nil {
				logger.Error(err)
			}
		})

	m := middlewares.New(logger)

	server := Server{
		httpHandlers: handlers.New(options, logger, service.New(logger, options, s)),
		options:      options,
		logger:       logger,
		middlewares:  &m,
		term:         term,
	}

	return &server, ctx, err
}

func (s Server) Run(ctx context.Context) {
	r := chi.NewRouter()

	r.Use(s.middlewares.Logging, s.middlewares.Compression)

	r.Get("/ping", s.httpHandlers.Ping)
	r.Route("/api/shorten", func(r chi.Router) {
		r.Post("/", s.httpHandlers.Shorten)
		r.Post("/batch", s.httpHandlers.ShortenBatch)
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", s.httpHandlers.GET)
		r.Post("/", s.httpHandlers.POST)
	})

	httpServer := &http.Server{
		Addr:    s.options.ServerAddress,
		Handler: r,
	}

	s.term.WithOrder(1).
		WithName("HTTPServer").
		Register(500*time.Millisecond, func(ctx context.Context) {
			s.logger.Info("terminating HTTPServer...")
			defer s.logger.Info("...HTTPServer terminated")

			if err := httpServer.Shutdown(ctx); err != nil {
				s.logger.Error(err)
			}
		})

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("terminated HTTP Server: %+v\n", err)
		}
	}()

	s.logger.Info("server started on address: \"" + s.options.ServerAddress + "\"")

	if err := s.term.Wait(ctx, 5*time.Second); err != nil {
		s.logger.Error("graceful termination period was timed out: %+v", err)
	}
}

func getStorage(o *config.Options, logger *logger.Logger) (service.Repository, error) {
	if o.DatabaseDsn != "" {
		return storage.NewDBStorage(o, logger)
	}

	if o.FileStoragePath != "" {
		return storage.NewFileStorage(o)
	}

	return storage.NewMemStorage(), nil
}
