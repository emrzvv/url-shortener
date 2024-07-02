package endpoint

import (
	"fmt"
	"github.com/emrzvv/url-shortener/cfg"
	"github.com/emrzvv/url-shortener/internal/app/service"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func Shorten(w http.ResponseWriter, r *http.Request, s service.URLShortenerService) {
	if r.Header.Get("Content-Type") != "text/plain; charset=utf-8" ||
		r.Body == http.NoBody {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shorten, shortenErr := s.GenerateShortURL(string(body))
	if shortenErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", cfg.Cfg.BaseAddress, shorten)))
}

func GetByID(w http.ResponseWriter, r *http.Request, s service.URLShortenerService) {
	id := chi.URLParam(r, "id")

	origin, err := s.GetOriginURLByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", origin)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
