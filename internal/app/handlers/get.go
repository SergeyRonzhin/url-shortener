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

	url, exist := h.storage.Get(id)

	if !exist {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Add("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
