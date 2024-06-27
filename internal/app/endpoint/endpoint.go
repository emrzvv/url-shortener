package endpoint

import (
	storage "github.com/emrzvv/url-shortener/internal/app/db"
	"github.com/emrzvv/url-shortener/internal/app/service"
	"io"
	"net/http"
)

// TODO: middleware? / (general request validation) -> /validate (body validation) -> /shorten (response with shortened url)
func Shorten(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	if r.Method != http.MethodPost ||
		r.Header.Get("Content-Type") != "text/plain; charset=utf-8" ||
		r.Body == nil ||
		r.Body == http.NoBody {
		w.WriteHeader(http.StatusBadRequest)
	} else if body, err := io.ReadAll(r.Body); err == nil && service.IsURLValid(string(body)) {
		w.WriteHeader(http.StatusCreated)
		url := service.GenerateShortenURL(6)
		for _, ok := s.Get(url); ok; url = service.GenerateShortenURL(6) {
		}
		s.Set(url, string(body))
		w.Write([]byte("http://localhost:8080/" + url)) // TODO: get self address from config or smth
	} else {
		w.WriteHeader(http.StatusBadRequest) // TODO: dry
	}
}

func GetByID(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		if id := r.PathValue("id"); !service.IsIDValid(id) {
			w.WriteHeader(http.StatusBadRequest)
		} else if value, ok := s.Get(id); ok {
			w.Header().Set("Location", value)
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
