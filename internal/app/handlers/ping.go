package handlers

import (
	"net/http"
)

func (h HTTPHandler) Ping(rw http.ResponseWriter, rq *http.Request) {

	if err := h.shortener.Ping(rq.Context()); err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.Info("Connection to database successful!")
	rw.WriteHeader(http.StatusOK)
}
