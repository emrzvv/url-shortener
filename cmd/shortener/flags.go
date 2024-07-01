package main

import (
	"flag"
	"github.com/emrzvv/url-shortener/cfg"
)

func parseFlags() {
	flag.StringVar(&cfg.Cfg.RunAddress, "a", "localhost:8080", "The address of the server")
	flag.StringVar(&cfg.Cfg.BaseAddress, "b", "http://localhost:8080", "Base address of shortened URLs")
	flag.Parse()
}
