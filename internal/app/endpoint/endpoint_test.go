package endpoint

import (
	"context"
	"github.com/emrzvv/url-shortener/cfg"
	storage "github.com/emrzvv/url-shortener/internal/app/db"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type expectedS struct {
	code    int
	headers map[string]string
	body    string
}

type testCase struct {
	name           string
	method         string
	url            string
	urlParams      map[string]string
	headers        map[string]string
	headersToCheck []string
	body           string
	doStorageOps   func(s storage.Storage)
	cleanStorage   bool
	expected       expectedS
}

func (test *testCase) run(t *testing.T, s storage.Storage) func(func(http.ResponseWriter, *http.Request, storage.Storage)) {
	return func(handler func(http.ResponseWriter, *http.Request, storage.Storage)) {
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
			if test.cleanStorage {
				defer s.Clear()
			}
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			for k, v := range test.urlParams {
				rctx.URLParams.Add(k, v)
			}

			request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, rctx))

			handler(w, request, s)
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

func TestShorten(t *testing.T) {
	s := &storage.InMemoryDB{}
	s.Init()
	cfg.Cfg.LoadDefault()
	tests := []testCase{
		//{
		//	name:   "#1 check not allowed method",
		//	method: http.MethodPut,
		//	url:    "/",
		//	headers: map[string]string{
		//		"Content-Type": "text/plain; charset=utf-8",
		//	},
		//	headersToCheck: []string{},
		//	body:           "https://ya.ru",
		//	doStorageOps:   func(s storage.Storage) {},
		//	expected: expectedS{
		//		code:    400,
		//		headers: map[string]string{},
		//		body:    "",
		//	},
		//},
		{
			name:   "#2 check wrong content-type header",
			method: http.MethodPost,
			url:    "/",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			headersToCheck: []string{},
			body:           "https://ya.ru",
			doStorageOps:   func(s storage.Storage) {},
			expected: expectedS{
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
			doStorageOps:   func(s storage.Storage) {},
			expected: expectedS{
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
			doStorageOps:   func(s storage.Storage) {},
			expected: expectedS{
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
			doStorageOps:   func(s storage.Storage) {},
			expected: expectedS{
				code: 201,
				headers: map[string]string{
					"Content-Type": "text/plain; charset=utf-8",
				},
				body: "^http://localhost:8080/[0-9a-zA-Z]{6}$",
			},
		},
	}
	for _, test := range tests {
		test.run(t, s)(Shorten)
	}
}

func TestGetByID(t *testing.T) {
	s := &storage.InMemoryDB{}
	s.Init()
	cfg.Cfg.LoadDefault()

	tests := []testCase{
		//{
		//	name:           "#1 wrong method",
		//	method:         http.MethodPost,
		//	url:            "/AAAAAA/",
		//	headers:        map[string]string{},
		//	headersToCheck: []string{},
		//	body:           "",
		//	doStorageOps:   func(s storage.Storage) {},
		//	expected: expectedS{
		//		code:    400,
		//		headers: map[string]string{},
		//		body:    "",
		//	},
		//},
		{
			name:   "#2 invalid id",
			method: http.MethodGet,
			url:    "/{id}/",
			urlParams: map[string]string{
				"id": "00000000000",
			},
			headers:        map[string]string{},
			headersToCheck: []string{},
			body:           "",
			doStorageOps:   func(s storage.Storage) {},
			expected: expectedS{
				code:    400,
				headers: map[string]string{},
				body:    "",
			},
		},
		{
			name:   "#3 unknown id",
			method: http.MethodGet,
			url:    "/{id}/",
			urlParams: map[string]string{
				"id": "000000",
			},
			headers:        map[string]string{},
			headersToCheck: []string{},
			body:           "",
			doStorageOps: func(s storage.Storage) {
				s.Set("111111", "https://ya.ru")
			},
			cleanStorage: true,
			expected: expectedS{
				code:    400,
				headers: map[string]string{},
				body:    "",
			},
		},
		{
			name:   "#4 default",
			method: http.MethodGet,
			url:    "/{id}/",
			urlParams: map[string]string{
				"id": "111111",
			},
			headers:        map[string]string{},
			headersToCheck: []string{"Location"},
			body:           "",
			doStorageOps: func(s storage.Storage) {
				s.Set("111111", "https://ya.ru")
				s.Set("222222", "https://vk.com")
				s.Set("333333", "https://t.me")
			},
			cleanStorage: true,
			expected: expectedS{
				code: 307,
				headers: map[string]string{
					"Location": "https://ya.ru",
				},
				body: "",
			},
		},
	}
	for _, test := range tests {
		test.run(t, s)(GetByID)
	}
}
