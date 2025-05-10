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

	original := string(body)

	if original == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, link, err := h.shortener.GetShortLink(rq.Context(), original)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
	}

	rw.Header().Add("content-type", "text/plain")

	if exists {
		rw.WriteHeader(http.StatusConflict)
	} else {
		rw.WriteHeader(http.StatusCreated)
	}

	_, err = fmt.Fprint(rw, link)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
	}
}
