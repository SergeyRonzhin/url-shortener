package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type originalURL struct {
	UUID string `json:"correlation_id"`
	URL  string `json:"original_url"`
}

type shortURL struct {
	UUID string `json:"correlation_id"`
	URL  string `json:"short_url"`
}

func (h HTTPHandler) ShortenBatch(rw http.ResponseWriter, rq *http.Request) {

	contentType := rq.Header.Get("content-type")

	if !strings.Contains(contentType, "application/json") {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURLs := make([]originalURL, 0)
	err := json.NewDecoder(rq.Body).Decode(&originalURLs)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(originalURLs) == 0 {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	urlsMap := make(map[string]string, len(originalURLs))

	for _, u := range originalURLs {
		urlsMap[u.UUID] = u.URL
	}

	shortLinks, err := h.shortener.GetShortLinks(rq.Context(), urlsMap)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURLs := make([]shortURL, 0, len(shortLinks))

	for key, value := range shortLinks {
		shortURLs = append(shortURLs, shortURL{key, value})
	}

	rs, err := json.Marshal(shortURLs)

	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Set("content-type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(rs)
}
