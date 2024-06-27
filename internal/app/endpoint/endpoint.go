package endpoint

import (
	"fmt"
	storage "github.com/emrzvv/url-shortener/internal/app/db"
	"github.com/emrzvv/url-shortener/internal/app/service"
	"io"
	"net/http"
)

// TODO: middleware? / (general request validation) -> /validate (body validation) -> /shorten (response with shortened url)
func Shorten(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	if r.Method != http.MethodPost ||
		r.Header.Get("Content-Type") != "text/plain" ||
		r.Body == nil ||
		r.Body == http.NoBody {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "bad request 400")
	} else if body, err := io.ReadAll(r.Body); err == nil && service.IsUrlValid(string(body)) {
		w.WriteHeader(http.StatusCreated)
		url := service.GenerateShortenUrl(6)
		for _, ok := s.Get(url); ok; url = service.GenerateShortenUrl(6) {
		}
		s.Set(url, string(body))
		w.Write([]byte("http://localhost:8080/" + url)) // TODO: get self address from config or smth
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid url 400") // TODO: dry
	}
}

func GetById(w http.ResponseWriter, r *http.Request, s *storage.Storage) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "bad request 400")
	} else {
		if id := r.PathValue("id"); !service.IsIdValid(id) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad request 400") // TODO: dry
		} else if value, ok := s.Get(id); ok {
			w.Header().Add("Location", value)
			w.WriteHeader(http.StatusTemporaryRedirect)
			fmt.Fprintf(w, "temporary redirect 307")
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad request 400") // TODO: dry
		}
	}
}
