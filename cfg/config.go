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

	runAddressFlag := flag.String("a", "localhost:8080", "The address of the server")
	baseAddressFlag := flag.String("b", "http://localhost:8080", "Base address of shortened URLs")

	if cfg.RunAddress == "" || cfg.BaseAddress == "" {
		flag.Parse()
	}

	if cfg.RunAddress == "" {
		cfg.RunAddress = *runAddressFlag
	}
	if cfg.BaseAddress == "" {
		cfg.BaseAddress = *baseAddressFlag
	}
	return nil
}

func (cfg *Config) LoadDefault() {
	cfg.RunAddress = "localhost:8080"
	cfg.BaseAddress = "http://localhost:8080"
}
