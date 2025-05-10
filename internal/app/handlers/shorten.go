package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ShortenRq struct {
	URL string `json:"url"`
}

type ShortenRs struct {
	Result string `json:"result"`
}

func (h HTTPHandler) Shorten(rw http.ResponseWriter, rq *http.Request) {

	contentType := rq.Header.Get("content-type")

	if !strings.Contains(contentType, "application/json") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	dataRq := ShortenRq{}
	err := json.NewDecoder(rq.Body).Decode(&dataRq)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if dataRq.URL == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	exists, link, err := h.shortener.GetShortLink(rq.Context(), dataRq.URL)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rs, err := json.Marshal(ShortenRs{link})

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("content-type", "application/json")

	if exists {
		rw.WriteHeader(http.StatusConflict)
	} else {
		rw.WriteHeader(http.StatusCreated)
	}

	_, err = rw.Write(rs)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
	}
}
