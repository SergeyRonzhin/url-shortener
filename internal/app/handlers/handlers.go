package handlers

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
)

type HTTPHandler struct {
	storage storage.Repository
}

func New(storage storage.Repository) HTTPHandler {
	return HTTPHandler{storage}
}
