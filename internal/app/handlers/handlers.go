package handlers

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
)

type HTTPHandler struct {
	options config.Options
	storage storage.Repository
}

func New(options config.Options, storage storage.Repository) HTTPHandler {
	return HTTPHandler{options, storage}
}
