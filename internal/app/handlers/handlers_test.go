package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/SergeyRonzhin/url-shortener/internal/app/service"
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	options *config.Options
)

func TestMain(m *testing.M) {
	var err error
	options, err = config.New()

	if err != nil {
		panic(err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestPOST(t *testing.T) {
	type request struct {
		method      string
		contentType string
		body        string
	}

	type expected struct {
		code        int
		contentType string
		emptyResult bool
	}

	tests := []struct {
		name     string
		rq       request
		expected expected
	}{
		{
			name: "successed returns short link",
			rq: request{
				method:      http.MethodPost,
				contentType: "text/plain",
				body:        "https://google.com",
			},
			expected: expected{
				code:        http.StatusCreated,
				contentType: "text/plain",
				emptyResult: false,
			},
		},
		{
			name: "returns 400 status for empty body in request",
			rq: request{
				method:      http.MethodPost,
				contentType: "text/plain",
				body:        "",
			},
			expected: expected{
				code:        http.StatusBadRequest,
				contentType: "",
				emptyResult: true,
			},
		},
		{
			name: "returns 400 status for unsupported content type",
			rq: request{
				method:      http.MethodPost,
				contentType: "application/xml",
				body:        "<url>https://google.com</url>",
			},
			expected: expected{
				code:        http.StatusBadRequest,
				contentType: "",
				emptyResult: true,
			},
		},
	}

	store := storage.NewMemStorage()
	logger, err := logger.New(options)

	if err != nil {
		panic(err)
	}

	httpHandler := New(options, logger, service.New(logger, options, store))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rq := httptest.NewRequest(test.rq.method, "/", strings.NewReader(test.rq.body))
			rq.Header.Set("Content-Type", test.rq.contentType)

			recorder := httptest.NewRecorder()

			httpHandler.POST(recorder, rq)

			rs := recorder.Result()

			assert.Equal(t, test.expected.code, rs.StatusCode)
			assert.Equal(t, test.expected.contentType, rs.Header.Get("Content-Type"))

			defer rs.Body.Close()

			result, err := io.ReadAll(rs.Body)

			require.NoError(t, err)

			if !test.expected.emptyResult {
				assert.NotEmpty(t, string(result))
			} else {
				assert.Empty(t, string(result))
			}
		})
	}
}

func TestGET(t *testing.T) {
	type request struct {
		method string
		param  string
	}

	type expected struct {
		code     int
		location string
	}

	tests := []struct {
		name     string
		rq       request
		expected expected
	}{
		{
			name: "successed redirect to original url",
			rq: request{
				method: http.MethodGet,
				param:  "QWerTy",
			},
			expected: expected{
				code:     http.StatusTemporaryRedirect,
				location: "https://google.com",
			},
		},
		{
			name: "returns 400 status for empty short link in request",
			rq: request{
				method: http.MethodGet,
				param:  "",
			},
			expected: expected{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
		{
			name: "returns 400 status for unknown short link",
			rq: request{
				method: http.MethodGet,
				param:  "test",
			},
			expected: expected{
				code:     http.StatusBadRequest,
				location: "",
			},
		},
	}

	store := storage.NewMemStorage()
	logger, err := logger.New(options)

	if err != nil {
		panic(err)
	}

	url := storage.URL{
		UUID:     uuid.NewString(),
		Short:    "QWerTy",
		Original: "https://google.com",
	}

	store.Add(context.TODO(), url)

	httpHandler := New(options, logger, service.New(logger, options, store))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rq := httptest.NewRequest(test.rq.method, "/"+test.rq.param, nil)
			recorder := httptest.NewRecorder()

			httpHandler.GET(recorder, rq)

			rs := recorder.Result()
			defer rs.Body.Close()

			assert.Equal(t, test.expected.code, rs.StatusCode)
			assert.Equal(t, test.expected.location, rs.Header.Get("Location"))
		})
	}
}

func TestShorten(t *testing.T) {
	type request struct {
		method      string
		contentType string
		body        string
	}

	type expected struct {
		code        int
		contentType string
		emptyResult bool
	}

	tests := []struct {
		name     string
		rq       request
		expected expected
	}{
		{
			name: "successed returns short link",
			rq: request{
				method:      http.MethodPost,
				contentType: "application/json",
				body:        `{ "url": "https://practicum.yandex.ru" }`,
			},
			expected: expected{
				code:        http.StatusCreated,
				contentType: "application/json",
				emptyResult: false,
			},
		},
		{
			name: "returns 400 status for empty url in request",
			rq: request{
				method:      http.MethodPost,
				contentType: "application/json",
				body:        `{ "url": "" }`,
			},
			expected: expected{
				code:        http.StatusBadRequest,
				contentType: "",
				emptyResult: true,
			},
		},
		{
			name: "returns 400 status for unsupported content type",
			rq: request{
				method:      http.MethodPost,
				contentType: "application/xml",
				body:        "<url>https://google.com</url>",
			},
			expected: expected{
				code:        http.StatusBadRequest,
				contentType: "",
				emptyResult: true,
			},
		},
	}

	store := storage.NewMemStorage()
	logger, err := logger.New(options)

	if err != nil {
		panic(err)
	}

	httpHandler := New(options, logger, service.New(logger, options, store))

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rq := httptest.NewRequest(test.rq.method, "/api/shorten", strings.NewReader(test.rq.body))
			rq.Header.Set("Content-Type", test.rq.contentType)

			recorder := httptest.NewRecorder()

			httpHandler.Shorten(recorder, rq)

			rs := recorder.Result()

			assert.Equal(t, test.expected.code, rs.StatusCode)
			assert.Equal(t, test.expected.contentType, rs.Header.Get("Content-Type"))

			defer rs.Body.Close()

			dataRs := ShortenRs{}

			if test.expected.contentType == "application/json" {

				err := json.NewDecoder(rs.Body).Decode(&dataRs)
				require.NoError(t, err)
			}

			if !test.expected.emptyResult {
				assert.NotEmpty(t, dataRs.Result)
			} else {
				assert.Empty(t, dataRs.Result)
			}
		})
	}
}
