package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
	serv := httptest.NewServer(CarRouter())
	defer serv.Close()

	var testTable = []struct {
		url    string
		want   string
		status int
	}{
		{"/cars/renault/Logan", "Renault Logan", http.StatusOK},
		{"/cars/audi/a4", "Audi A4", http.StatusOK},
		// проверим на ошибочный запрос
		{"/cars/audi/a6", "unknown car: audi a6\n", http.StatusNotFound},
		{"/cars/BMW/M5", "BMW M5", http.StatusOK},
		{"/cars/renault/X6", "unknown car: renault X6\n", http.StatusNotFound},
		{"/cars/Vw/Passat", "VW Passat", http.StatusOK},
	}

	for _, test := range testTable {
		rs, body := testRequest(t, serv, http.MethodGet, test.url)
		defer rs.Body.Close()

		assert.Equal(t, test.status, rs.StatusCode)
		assert.Equal(t, test.want, body)
	}
}
