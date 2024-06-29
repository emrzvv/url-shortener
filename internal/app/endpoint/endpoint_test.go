package endpoint

import (
	storage "github.com/emrzvv/url-shortener/internal/app/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShorten(t *testing.T) {
	s := &storage.Storage{}
	s.Init()
	type expected struct {
		code    int
		headers map[string]string
		body    string
	}

	tests := []struct {
		name           string
		method         string
		url            string
		headers        map[string]string
		headersToCheck []string
		body           string
		doStorageOps   func(s *storage.Storage)
		expected       expected
	}{
		{
			name:   "#1 check not allowed method",
			method: http.MethodPut,
			url:    "/",
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			headersToCheck: []string{},
			body:           "https://ya.ru",
			doStorageOps: func(s *storage.Storage) {
				return
			},
			expected: expected{
				code:    400,
				headers: map[string]string{},
				body:    "",
			},
		},
		{
			name:   "#2 check wrong content-type header",
			method: http.MethodPost,
			url:    "/",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			headersToCheck: []string{},
			body:           "https://ya.ru",
			doStorageOps: func(s *storage.Storage) {
				return
			},
			expected: expected{
				code:    400,
				headers: map[string]string{},
				body:    "",
			},
		},
		{
			name:   "#3 check empty body",
			method: http.MethodPost,
			url:    "/",
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			headersToCheck: []string{},
			body:           "",
			doStorageOps: func(s *storage.Storage) {
				return
			},
			expected: expected{
				code:    400,
				headers: map[string]string{},
				body:    "",
			},
		},
		{
			name:   "#4 check invalid url in body",
			method: http.MethodPost,
			url:    "/",
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			headersToCheck: []string{},
			body:           "httpts:///y2o3ria.2348kjnru.slkdf//skl",
			doStorageOps: func(s *storage.Storage) {
				return
			},
			expected: expected{
				code:    400,
				headers: map[string]string{},
				body:    "",
			},
		},
		{
			name:   "#5 check short url generation",
			method: http.MethodPost,
			url:    "/",
			headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
			headersToCheck: []string{"Content-Type"},
			body:           "https://ya.ru",
			doStorageOps: func(s *storage.Storage) {
				return
			},
			expected: expected{
				code: 201,
				headers: map[string]string{
					"Content-Type": "text/plain; charset=utf-8",
				},
				body: "^https?://localhost:8080/[0-9a-zA-Z]{6}$",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var requestBody io.Reader
			if test.body != "" {
				requestBody = strings.NewReader(test.body)
			} else {
				requestBody = http.NoBody
			}
			request := httptest.NewRequest(test.method, test.url, requestBody)
			for k, v := range test.headers {
				request.Header.Add(k, v)
			}
			test.doStorageOps(s)

			w := httptest.NewRecorder()
			Shorten(w, request, s)
			result := w.Result()

			require.Equal(t, test.expected.code, result.StatusCode)
			for _, header := range test.headersToCheck {
				require.Equal(t, test.expected.headers[header], result.Header.Get(header))
			}
			defer result.Body.Close()
			resultBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			if test.expected.body != "" {
				assert.Regexp(t, test.expected.body, string(resultBody))
			}
		})
	}
}
