package middlewares

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
)

type Middlewares struct {
	logger *logger.Logger
}

func New(logger *logger.Logger) Middlewares {
	return Middlewares{logger}
}
