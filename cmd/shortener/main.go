package main

import (
	"fmt"
	storage "github.com/emrzvv/url-shortener/internal/app/db"
	"github.com/emrzvv/url-shortener/internal/app/endpoint"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		return
	}
}

func run() error {
	s := &storage.Storage{}
	s.Init()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		endpoint.Shorten(w, r, s)
	})
	mux.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		endpoint.GetById(w, r, s)
	})
	address := ":8080"

	fmt.Println("Listening on " + address)
	if err := http.ListenAndServe(address, mux); err != nil {
		return err
	}

	return nil
}
