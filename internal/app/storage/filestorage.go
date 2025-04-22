package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
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

func NewFileStorage(options *config.Options) (*FileStorage, error) {

	file, err := os.OpenFile(options.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	defer func() {
		err = file.Close()

		if err != nil {
			return
		}
	}()

	scan := bufio.NewScanner(file)
	urls := make(map[string]URL)

	for scan.Scan() {
		url := URL{}
		err = json.Unmarshal(scan.Bytes(), &url)

		if err != nil {
			return nil, err
		}

		urls[url.ShortURL] = url
	}

	return &FileStorage{URLs: urls, Size: len(urls), options: options}, nil
}

func (s *FileStorage) Get(key string) (string, bool) {
	s.mu.Lock()
	url, exist := s.URLs[key]
	s.mu.Unlock()

	return url.OriginalURL, exist
}

func (s *FileStorage) Add(key string, value string) error {
	s.mu.Lock()

	url := URL{
		UUID:        strconv.Itoa(s.Size + 1),
		ShortURL:    key,
		OriginalURL: value,
	}

	s.URLs[url.ShortURL] = url

	data, err := json.Marshal(url)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(s.options.FileStoragePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	defer func() {
		err = file.Close()

		if err != nil {
			return
		}
	}()

	writer := bufio.NewWriter(file)

	_, err = writer.Write(data)

	if err != nil {
		return err
	}

	_, err = writer.WriteString("\n")

	if err != nil {
		return err
	}

	err = writer.Flush()

	if err != nil {
		return err
	}

	s.mu.Unlock()

	return nil
}

func (s *FileStorage) ContainsValue(value string) (bool, string) {
	for key, url := range s.URLs {
		if url.OriginalURL == value {
			return true, key
		}
	}

	return false, ""
}
