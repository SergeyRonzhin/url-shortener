package storage

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
)

type URL struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type FileStorage struct {
	URLs    map[string]URL
	Size    int
	options *config.Options
	mu      sync.Mutex
}

func NewFileStorage(options *config.Options) FileStorage {

	data, err := os.ReadFile(options.FileStoragePath)

	if err != nil {
		// todo:
	}

	lines := []string{}

	if len(data) != 0 {
		lines = strings.Split(string(data), "\n")
	}

	urls := make(map[string]URL, len(lines))

	for _, line := range lines {
		url := URL{}

		if line == "" {
			continue
		}

		err = json.Unmarshal([]byte(line), &url)

		if err != nil {
			// todo:
		}

		urls[url.ShortURL] = url
	}

	return FileStorage{
		URLs:    urls,
		Size:    len(urls),
		options: options,
	}
}

func (s *FileStorage) Get(key string) (string, bool) {
	s.mu.Lock()
	url, exist := s.URLs[key]
	s.mu.Unlock()

	return url.OriginalURL, exist
}

func (s *FileStorage) Add(key string, value string) {
	s.mu.Lock()
	// s.URLs[key] = value
	s.mu.Unlock()
}

func (s *FileStorage) ContainsValue(value string) (bool, string) {
	for key, url := range s.URLs {
		if url.OriginalURL == value {
			return true, key
		}
	}

	return false, ""
}
