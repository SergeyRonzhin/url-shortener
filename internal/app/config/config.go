package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Options struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	BaseURL       string `env:"BASE_URL"`
}

func (o *Options) Init() {

	err := env.Parse(o)

	if err != nil {
		fmt.Println(err)
	}

	if o.ServerAddress == "" {
		flag.StringVar(&o.ServerAddress, "a", "localhost:8080", "Address for hosting service")
	}

	if o.BaseURL == "" {
		flag.StringVar(&o.BaseURL, "b", "http://localhost:8080", "Base address for short links")
	}
}
