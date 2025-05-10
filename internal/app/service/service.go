package service

import (
	"context"
	"math/rand"
	"time"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
	"github.com/SergeyRonzhin/url-shortener/internal/app/logger"
	"github.com/SergeyRonzhin/url-shortener/internal/app/storage"
	"github.com/google/uuid"
)

type Repository interface {
	Add(ctx context.Context, url storage.URL) error
	Batch(ctx context.Context, urls []storage.URL) error
	GetShortURL(ctx context.Context, original string) (bool, string)
	GetOriginalURL(ctx context.Context, short string) (bool, string)
	Close() error
	Ping(ctx context.Context) error
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

func (s URLShortener) GetShortLink(ctx context.Context, originalURL string) (bool, string, error) {

	exists, shortURL := s.storage.GetShortURL(ctx, originalURL)

	if exists {
		return true, s.getShortWithBaseURL(shortURL), nil
	}

	shortURL = generateShortURL(8)

	url := storage.URL{
		UUID:     uuid.NewString(),
		Original: originalURL,
		Short:    shortURL,
	}

	err := s.storage.Add(ctx, url)
	return false, s.getShortWithBaseURL(url.Short), err
}

func (s URLShortener) GetShortLinks(ctx context.Context, originalURLs map[string]string) (map[string]string, error) {

	urls := make([]storage.URL, 0, len(originalURLs))
	result := make(map[string]string, len(originalURLs))

	for uuid, original := range originalURLs {

		exists, shortURL := s.storage.GetShortURL(ctx, original)

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
		err = s.storage.Batch(ctx, urls)
	}

	return result, err
}

func (s URLShortener) GetOriginalLink(ctx context.Context, shortLink string) (bool, string) {
	return s.storage.GetOriginalURL(ctx, shortLink)
}

func (s URLShortener) Ping(ctx context.Context) error {
	return s.storage.Ping(ctx)
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
