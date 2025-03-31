package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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

	fmt.Printf("url: %s\n", url)

	if url == "" {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	shortLink := ""

	if result, key := h.storage.ContainsValue(url); result {
		shortLink = key
	} else {
		shortLink = generateShortLink(8)
		h.storage.Add(shortLink, url)
	}

	fmt.Printf("shortLink: %s\n", shortLink)

	rw.Header().Add("content-type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	fmt.Fprint(rw, h.options.Endpoint+shortLink)
}

func generateShortLink(length int) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)

	for i := range length {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
