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
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	url := string(body)

	if url == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink, err := h.shortener.GetShortLink(url)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
	}

	rw.Header().Add("content-type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	_, err = fmt.Fprint(rw, h.options.BaseURL+"/"+shortLink)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
	}
}
