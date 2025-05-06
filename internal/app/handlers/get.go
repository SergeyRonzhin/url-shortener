package handlers

import (
	"net/http"
	"strings"
)

func (h HTTPHandler) GET(rw http.ResponseWriter, rq *http.Request) {
	id := strings.TrimPrefix(rq.URL.Path, "/")

	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	exist, link := h.shortener.GetOriginalLink(id)

	if !exist {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Add("Location", link)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
