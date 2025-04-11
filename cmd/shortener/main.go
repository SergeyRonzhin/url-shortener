package main

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/server"
)

func main() {

	options := config.Options{}
	options.Init()

	if err := server.New(options).Run(); err != nil {
		panic(err)
	}
}
