package main

import (
	"github.com/emrzvv/url-shortener/cfg"
	storage "github.com/emrzvv/url-shortener/internal/app/db"
	"github.com/emrzvv/url-shortener/internal/app/endpoint"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {
	var err error

	err = cfg.Cfg.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err = run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	s := &storage.InMemoryDB{}
	s.Init()
	router := chi.NewRouter()
	router.Use(chimw.Logger)
	router.Route("/", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			endpoint.Shorten(w, r, s)
		})
		r.Route("/{id}", func(r chi.Router) { // TODO: regex constraint :[0-9a-zA-Z]{6} but returns 404 instead 400
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				endpoint.GetByID(w, r, s)
			})
		})
	})

	log.Println("Listening on " + cfg.Cfg.RunAddress)
	if err := http.ListenAndServe(cfg.Cfg.RunAddress, router); err != nil {
		return err
	}

	return nil
}
