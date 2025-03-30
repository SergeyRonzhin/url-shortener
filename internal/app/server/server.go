package server

import (
	"net/http"

	"github.com/SergeyRonzhin/url-shortener/internal/app/handlers"
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.httpHandlers.Index)

	return http.ListenAndServe(":8080", mux)
}
