package main

import (
	"net/http"

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

	switch rq.Method {
	case http.MethodPost:
		httpHandlers.POSTHandler(rw, rq)
	case http.MethodGet:
		httpHandlers.GETHandler(rw, rq)
	default:
		rw.WriteHeader(http.StatusBadRequest)
	}
}
