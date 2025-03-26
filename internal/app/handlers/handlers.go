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

type HTTPHandler struct {
	storage storage.Repository
}

func NewHTTPHandler(storage storage.Repository) HTTPHandler {
	return HTTPHandler{storage}
}

func (h HTTPHandler) Index(rw http.ResponseWriter, rq *http.Request) {

	switch rq.Method {
	case http.MethodPost:
		h.IndexPOST(rw, rq)
	case http.MethodGet:
		h.IndexGET(rw, rq)
	default:
		rw.WriteHeader(http.StatusBadRequest)
	}
}

func (h HTTPHandler) IndexGET(rw http.ResponseWriter, rq *http.Request) {
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

func (h HTTPHandler) IndexPOST(rw http.ResponseWriter, rq *http.Request) {

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
