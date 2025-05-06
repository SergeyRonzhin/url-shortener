package main

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/SergeyRonzhin/url-shortener/internal/app/server"
)

func main() {

	options, err := config.New()

	if err != nil {
		panic(err)
	}

	logger, err := logger.New(options)

	if err != nil {
		panic(err)
	}

	server, ctx, err := server.New(options, logger)

	if err != nil {
		panic(err)
	}

	server.Run(ctx)
}
