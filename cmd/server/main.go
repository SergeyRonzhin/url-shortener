package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

type Car struct {
	ID    int
	Brand string
	Model string
}

var cars = []Car{
	{ID: 1, Brand: "Renault", Model: "Logan"},
	{ID: 2, Brand: "Renault", Model: "Duster"},
	{ID: 3, Brand: "BMW", Model: "X6"},
	{ID: 4, Brand: "BMW", Model: "M5"},
	{ID: 5, Brand: "VW", Model: "Passat"},
	{ID: 6, Brand: "VW", Model: "Jetta"},
	{ID: 7, Brand: "Audi", Model: "A4"},
	{ID: 8, Brand: "Audi", Model: "Q7"},
}

func main() {
	http.ListenAndServe(":8080", CarRouter())
}

func CarRouter() chi.Router {

	r := chi.NewRouter()

	r.Route("/cars", func(r chi.Router) {
		r.Get("/", getCars)
		r.Route("/{brand}", func(r chi.Router) {
			r.Get("/", getCarsByBrand)
			r.Get("/{model}", getCarByModel)
		})
	})

	return r
}

func getCars(rw http.ResponseWriter, rq *http.Request) {
	result := make([]string, len(cars))

	for _, c := range cars {
		result = append(result, c.Brand+" "+c.Model)
	}

	rw.Write([]byte(strings.Join(result, ", ")))
}

func getCarsByBrand(rw http.ResponseWriter, rq *http.Request) {
	brand := chi.URLParam(rq, "brand")

	if brand == "" {
		http.Error(rw, "brand param is missed", http.StatusBadRequest)
		return
	}

	result := []string{}

	for _, c := range cars {
		if strings.EqualFold(c.Brand, brand) {
			result = append(result, c.Brand+" "+c.Model)
		}
	}

	rw.Write([]byte(strings.Join(result, ", ")))
}

func getCarByModel(rw http.ResponseWriter, rq *http.Request) {
	brand := chi.URLParam(rq, "brand")
	model := chi.URLParam(rq, "model")

	if brand == "" {
		http.Error(rw, "brand param is missed", http.StatusBadRequest)
		return
	}

	if model == "" {
		http.Error(rw, "model param is missed", http.StatusBadRequest)
		return
	}

	for _, c := range cars {
		if strings.EqualFold(c.Brand+" "+c.Model, brand+" "+model) {
			rw.Write([]byte(c.Brand + " " + c.Model))
			return
		}
	}

	http.Error(rw, fmt.Sprintf("unknown car: %s", brand+" "+model), http.StatusNotFound)
}

func testRequest(t *testing.T, ts *httptest.Server, method string, path string) (*http.Response, string) {
	rq, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	rs, err := ts.Client().Do(rq)
	require.NoError(t, err)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	require.NoError(t, err)

	return rs, string(body)
}
