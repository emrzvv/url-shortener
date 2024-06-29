package endpoint

import (
	storage "github.com/emrzvv/url-shortener/internal/app/db"
	"github.com/emrzvv/url-shortener/internal/app/service"
	"io"
	"net/http"
	"strings"
)

// TODO: middleware? / (general request validation) -> /validate (body validation) -> /shorten (response with shortened url)
func Shorten(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	if r.Method != http.MethodPost ||
		r.Header.Get("Content-Type") != "text/plain; charset=utf-8" ||
		r.Body == http.NoBody {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil || !service.IsURLValid(string(body)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	url := service.GenerateShortenURL(6)
	for _, ok := s.Get(url); ok; url = service.GenerateShortenURL(6) {
	}
	s.Set(url, string(body))
	w.Write([]byte("http://localhost:8080/" + url)) // TODO: get self address from config or smth
}

func GetByID(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//TODO: use id := r.PathValue("id")
	id := strings.Split(r.URL.Path, "/")[1]
	if id == "" || !service.IsIDValid(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	value, ok := s.Get(id)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Header().Set("Location", value)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
