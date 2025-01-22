package main

import (
	"github.com/emrzvv/url-shortener/cfg"
	storage "github.com/emrzvv/url-shortener/internal/app/db"
	"github.com/emrzvv/url-shortener/internal/app/endpoint"
	"github.com/emrzvv/url-shortener/internal/app/endpoint/middleware"
	"github.com/emrzvv/url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}

func run() error {
	config, configErr := cfg.LoadNewConfig()
	if configErr != nil {
		log.Fatal().Err(configErr).Msg("failed to load config")
	}

	db := storage.NewInMemoryDBStorage(make(map[string]string))
	urlShortenerService := service.NewURLShortenerService(db, config)

	router := chi.NewRouter()
	//router.Use(chimw.Logger)
	router.Use(middleware.Logger)
	router.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			endpoint.Shorten(w, r, urlShortenerService)
		})
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				endpoint.GetByID(w, r, urlShortenerService)
			})
		})
	})

	log.Info().Msgf("Listening on %s", config.RunAddress)
	if err := http.ListenAndServe(config.RunAddress, router); err != nil {
		return err
	}

	return nil
}
