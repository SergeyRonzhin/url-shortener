package handlers

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/service"
	"go.uber.org/zap"
)

type HTTPHandler struct {
	options   *config.Options
	logger    *zap.SugaredLogger
	shortener service.URLShortener
}

func New(options *config.Options, logger *zap.SugaredLogger, shortener service.URLShortener) HTTPHandler {
	return HTTPHandler{options, logger, shortener}
}
