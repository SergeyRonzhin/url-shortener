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
	DatabaseDsn     string `env:"DATABASE_DSN"`
	MigrationsPath  string `env:"MIGRATIONS_PATH"`
}

func New() (*Options, error) {

	o := &Options{}
	err := env.Parse(o)

	if err != nil {
		return nil, err
	}

	serverAddress := flag.String("a", "localhost:8080", "Address for hosting service")
	baseURL := flag.String("b", "http://localhost:8080", "Base address for short links")
	pathToFile := flag.String("f", "", "Path to file storage")
	logLevel := flag.String("log_level", "info", "Log level")
	logEncoding := flag.String("log_encode", "json", "Log encoding")
	databaseDsn := flag.String("d", "", "Connection string to database")
	migrationsPath := flag.String("m", ".", "Path to migration files")

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

	if o.DatabaseDsn == "" {
		o.DatabaseDsn = *databaseDsn
	}

	if o.MigrationsPath == "" {
		o.MigrationsPath = *migrationsPath
	}

	return o, nil
}
