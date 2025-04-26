package main

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/server"

	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewDevelopment()

	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	o, err := config.New()

	if err != nil {
		panic(err)
	}

	if err := server.New(o, logger.Sugar()).Run(); err != nil {
		panic(err)
	}
}
