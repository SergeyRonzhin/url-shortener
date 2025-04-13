package service

import (
	"math/rand"
	"time"
)

type Repository interface {
	Get(key string) (string, bool)
	Add(key string, value string)
	ContainsValue(value string) (bool, string)
}

type URLShortener struct {
	storage Repository
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func New(storage Repository) URLShortener {
	return URLShortener{storage}
}

func (s URLShortener) GetShortLink(url string) string {

	shortLink := ""

	if result, key := s.storage.ContainsValue(url); result {
		shortLink = key
	} else {
		shortLink = generateShortLink(8)
		s.storage.Add(shortLink, url)
	}

	return shortLink
}

func (s URLShortener) GetOriginalURL(shortLink string) (string, bool) {
	return s.storage.Get(shortLink)
}

func generateShortLink(length int) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)

	for i := range length {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}
