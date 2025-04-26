package server

import (
	"net/http"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/handlers"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/SergeyRonzhin/url-shortener/internal/app/middlewares"
	"github.com/SergeyRonzhin/url-shortener/internal/app/service"
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpHandlers handlers.HTTPHandler
	options      *config.Options
	logger       *logger.Logger
	middlewares  *middlewares.Middlewares
}

func New(options *config.Options, logger *logger.Logger) Server {
	s, err := storage.NewFileStorage(options)

	if err != nil {
		panic(err)
	}

	m := middlewares.New(logger)

	return Server{
		httpHandlers: handlers.New(options, logger, service.New(s)),
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
