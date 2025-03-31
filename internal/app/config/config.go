package config

import "flag"

type Options struct {
	Host     string
	Endpoint string
}

func (o *Options) Init() {
	flag.StringVar(&o.Host, "a", ":8080", "Address for hosting service")
	flag.StringVar(&o.Endpoint, "b", "http://localhost:8080/", "Base address for short links")
}
