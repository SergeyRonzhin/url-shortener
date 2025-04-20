package middlewares

import (
	"go.uber.org/zap"
)

type Middlewares struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) Middlewares {
	return Middlewares{logger}
}
