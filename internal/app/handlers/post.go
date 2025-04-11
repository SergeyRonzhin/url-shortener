package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (h HTTPHandler) POST(rw http.ResponseWriter, rq *http.Request) {

	contentType := rq.Header.Get("content-type")

	if !strings.Contains(contentType, "text/plain") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(rq.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	url := string(body)

	if url == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink := h.shortener.GetShortLink(url)

	rw.Header().Add("content-type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	_, err = fmt.Fprint(rw, h.options.BaseURL+"/"+shortLink)

	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusBadRequest)
	}
}
