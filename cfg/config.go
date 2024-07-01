package cfg

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	RunAddress  string `env:"SERVER_ADDRESS"`
	BaseAddress string `env:"BASE_URL"`
}

var Cfg = &Config{}

func (cfg *Config) Load() error {
	err := env.Parse(cfg)
	if err != nil {
		return err
	}
	if cfg.RunAddress == "" {
		flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "The address of the server")
	}
	if cfg.BaseAddress == "" {
		flag.StringVar(&cfg.BaseAddress, "b", "http://localhost:8080", "Base address of shortened URLs")
	}
	flag.Parse()
	return nil
}
