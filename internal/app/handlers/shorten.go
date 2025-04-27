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

	shortLink, err := h.shortener.GetShortLink(dataRq.URL)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	dataRs := ShortenRs{
		Result: h.options.BaseURL + "/" + shortLink,
	}

	rs, err := json.Marshal(dataRs)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(rs)
}
