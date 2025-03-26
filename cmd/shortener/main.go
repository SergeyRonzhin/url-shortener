package main

import (
	"github.com/SergeyRonzhin/url-shortener/internal/app/server"
)

func main() {
	if err := server.InitServer().Run(); err != nil {
		panic(err)
	}
}
