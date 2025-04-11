package handlers

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/service"
)

type HTTPHandler struct {
	options   config.Options
	shortener service.URLShortener
}

func New(options config.Options, shortener service.URLShortener) HTTPHandler {
	return HTTPHandler{options, shortener}
}
