package main

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/server"
)

func main() {
	if err := server.New().Run(); err != nil {
		panic(err)
	}
}
