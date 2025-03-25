package main

import (
	"net/http"
	"strings"

	"github.com/SergeyRonzhin/url-shortener/internal/app/handlers"
)

var (
	httpHandlers = handlers.NewHTTPHandlers()
)

func main() {
	runService()
}

func runService() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		panic(err)
	}
}

func mainHandler(rw http.ResponseWriter, rq *http.Request) {

	contentType := rq.Header.Get("content-type")

	if !strings.Contains(contentType, "text/plain") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if rq.Method == http.MethodPost {
		httpHandlers.POSTHandler(rw, rq)
		return
	}

	if rq.Method == http.MethodGet {
		httpHandlers.GETHandler(rw, rq)
		return
	}

	rw.WriteHeader(http.StatusBadRequest)
}
