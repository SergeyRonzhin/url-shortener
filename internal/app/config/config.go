package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Options struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func (o *Options) Init() {

	err := env.Parse(o)

	if err != nil {
		fmt.Println(err)
	}

	serverAddress := flag.String("a", "localhost:8080", "Address for hosting service")
	baseURL := flag.String("b", "http://localhost:8080", "Base address for short links")
	pathToFile := flag.String("f", "C:\\Sergey\\temp_files\\shortener_storage.json", "Path to file storage")

	flag.Parse()

	if o.ServerAddress == "" {
		o.ServerAddress = *serverAddress
	}

	if o.BaseURL == "" {
		o.BaseURL = *baseURL
	}

	if o.FileStoragePath == "" {
		o.FileStoragePath = *pathToFile
	}
}
