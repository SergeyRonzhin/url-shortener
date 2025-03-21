package main

import (
	"net/http"
	"strings"

	"github.com/SergeyRonzhin/url-shortener/internal/app/handlers"
)

func main() {
	runService()
}

func runService() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err.Error())
	}
}

func mainHandler(rw http.ResponseWriter, rq *http.Request) {

	id := strings.TrimPrefix(rq.URL.Path, "/")

	if id == "" && rq.Method == http.MethodPost {
		handlers.POSTHandler(rw, rq)
		return
	}

	if id != "" && rq.Method == http.MethodGet {
		handlers.POSTHandler(rw, rq)
		return
	}

	http.Error(rw, "Bad request", http.StatusBadRequest)
}
