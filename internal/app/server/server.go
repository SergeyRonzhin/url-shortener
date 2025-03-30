package server

import (
	"net/http"

	"github.com/SergeyRonzhin/url-shortener/internal/app/handlers"
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpHandlers handlers.HTTPHandler
}

func New() Server {
	return Server{
		httpHandlers: handlers.New(storage.New()),
	}
}

func (s Server) Run() error {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/{id}", s.httpHandlers.GET)
		r.Post("/", s.httpHandlers.POST)
	})

	return http.ListenAndServe(":8080", r)
}
