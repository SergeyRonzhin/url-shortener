package main

import (
	"net/http"
)

func main() {
	runService()
}

func runService() {
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		panic(err.Error())
	}
}
