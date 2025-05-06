package service

import (
	"math/rand"
	"time"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
	"github.com/google/uuid"
)

type Repository interface {
	Add(url storage.URL) error
	Batch(urls []storage.URL) error
	GetShortURL(original string) (bool, string)
	GetOriginalURL(short string) (bool, string)
	Close() error
}

type URLShortener struct {
	logger  *logger.Logger
	options *config.Options
	storage Repository
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func New(logger *logger.Logger, options *config.Options, storage Repository) URLShortener {
	return URLShortener{logger, options, storage}
}

func (s URLShortener) GetShortLink(originalURL string) (bool, string, error) {

	exists, shortURL := s.storage.GetShortURL(originalURL)

	if exists {
		return true, s.getShortWithBaseURL(shortURL), nil
	}

	shortURL = generateShortURL(8)

	url := storage.URL{
		UUID:     uuid.NewString(),
		Original: originalURL,
		Short:    shortURL,
	}

	err := s.storage.Add(url)
	return false, s.getShortWithBaseURL(url.Short), err
}

func (s URLShortener) GetShortLinks(originalURLs map[string]string) (map[string]string, error) {

	urls := make([]storage.URL, 0, len(originalURLs))
	result := make(map[string]string, len(originalURLs))

	for uuid, original := range originalURLs {

		exists, shortURL := s.storage.GetShortURL(original)

		if exists {
			result[uuid] = s.getShortWithBaseURL(shortURL)
			continue
		}

		url := storage.URL{
			UUID:     uuid,
			Original: original,
			Short:    generateShortURL(8),
		}

		urls = append(urls, url)
		result[uuid] = s.getShortWithBaseURL(url.Short)
	}

	var err error

	if len(urls) != 0 {
		err = s.storage.Batch(urls)
	}

	return result, err
}

func (s URLShortener) GetOriginalLink(shortLink string) (bool, string) {
	return s.storage.GetOriginalURL(shortLink)
}

func generateShortURL(length int) string {

	rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, length)

	for i := range length {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

func (s URLShortener) getShortWithBaseURL(short string) string {
	return s.options.BaseURL + "/" + short
}
