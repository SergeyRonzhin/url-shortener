package handlers

import (
	"net/http"

	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
)

type HTTPHandler struct {
	storage storage.Repository
}

func New(storage storage.Repository) HTTPHandler {
	return HTTPHandler{storage}
}

func (h HTTPHandler) Index(rw http.ResponseWriter, rq *http.Request) {

	switch rq.Method {
	case http.MethodPost:
		h.POST(rw, rq)
	case http.MethodGet:
		h.GET(rw, rq)
	default:
		rw.WriteHeader(http.StatusBadRequest)
	}
}
