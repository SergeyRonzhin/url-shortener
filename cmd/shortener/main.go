package main

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/server"

	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewDevelopment()

	if err != nil {
		panic("cannot initialize zap logger")
	}

	defer logger.Sync()

	options := config.Options{}
	options.Init()

	if err := server.New(&options, logger.Sugar()).Run(); err != nil {
		panic(err)
	}
}
