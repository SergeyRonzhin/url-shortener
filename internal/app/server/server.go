package server

import (
	"net/http"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/handlers"
	"github.com/SergeyRonzhin/url-shortener/internal/app/middlewares"
	"github.com/SergeyRonzhin/url-shortener/internal/app/service"
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	httpHandlers handlers.HTTPHandler
	options      *config.Options
	logger       *zap.SugaredLogger
	middlewares  *middlewares.Middlewares
}

func New(options *config.Options, logger *zap.SugaredLogger) Server {
	s := storage.New()
	m := middlewares.New(logger)

	return Server{
		httpHandlers: handlers.New(options, logger, service.New(&s)),
		options:      options,
		logger:       logger,
		middlewares:  &m,
	}
}

func (s Server) Run() error {
	r := chi.NewRouter()

	r.Use(s.middlewares.Logging, s.middlewares.Compression)

	r.Post("/api/shorten", s.httpHandlers.Shorten)
	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", s.httpHandlers.GET)
		r.Post("/", s.httpHandlers.POST)
	})

	s.logger.Info("server started")

	return http.ListenAndServe(s.options.ServerAddress, r)
}
