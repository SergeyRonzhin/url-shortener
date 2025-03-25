package handlers

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type HTTPHandlers struct {
	storage storage.Storage
}

func NewHTTPHandlers() HTTPHandlers {
	return HTTPHandlers{storage.NewStorage()}
}

func (h HTTPHandlers) GETHandler(rw http.ResponseWriter, rq *http.Request) {
	id := strings.TrimPrefix(rq.URL.Path, "/")

	url, exist := h.storage.Get(id)

	if !exist {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	rw.Header().Add("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

func (h HTTPHandlers) POSTHandler(rw http.ResponseWriter, rq *http.Request) {
	body, err := io.ReadAll(rq.Body)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	url := string(body)

	fmt.Printf("url: %s", url)

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

	fmt.Printf("shortLink: %s", shortLink)

	rw.Header().Add("content-type", "text/plain")
	rw.WriteHeader(http.StatusCreated)

	fmt.Fprint(rw, "http://localhost:8080/"+shortLink)
}

func generateShortLink(length int) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)

	for i := range length {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
