package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Options struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func New() (*Options, error) {

	o := &Options{}
	err := env.Parse(o)

	if err != nil {
		return nil, err
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

	return o, nil
}
