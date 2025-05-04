package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Options struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	LogLevel        string `env:"LOG_LEVEL"`
	LogEncoding     string `env:"LOG_ENCODING"`
}

func New() (*Options, error) {

	o := &Options{}
	err := env.Parse(o)

	if err != nil {
		return nil, err
	}

	serverAddress := flag.String("a", "localhost:8080", "Address for hosting service")
	baseURL := flag.String("b", "http://localhost:8080", "Base address for short links")
	pathToFile := flag.String("f", "storage.json", "Path to file storage")
	logLevel := flag.String("log_level", "info", "Log level")
	logEncoding := flag.String("log_encode", "json", "Log encoding")

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

	if o.LogLevel == "" {
		o.LogLevel = *logLevel
	}

	if o.LogEncoding == "" {
		o.LogEncoding = *logEncoding
	}

	return o, nil
}
