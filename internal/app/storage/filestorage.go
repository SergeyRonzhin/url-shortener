package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/SergeyRonzhin/url-shortener/internal/app/config"
)

type FileStorage struct {
	urls    []URL
	options *config.Options
	mu      sync.Mutex
	file    *os.File
}

func NewFileStorage(options *config.Options) (*FileStorage, error) {

	file, err := os.OpenFile(options.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}

	scan := bufio.NewScanner(file)
	urls := make([]URL, 0)

	for scan.Scan() {
		url := URL{}
		err = json.Unmarshal(scan.Bytes(), &url)

		if err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	return &FileStorage{
		urls:    urls,
		options: options,
		file:    file,
	}, err
}

func (s *FileStorage) Add(url URL) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls = append(s.urls, url)
	data, err := json.Marshal(url)

	if err != nil {
		return err
	}

	writer := bufio.NewWriter(s.file)

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

	return err
}

func (s *FileStorage) Batch(urls []URL) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls = append(s.urls, urls...)
	writer := bufio.NewWriter(s.file)

	for _, url := range urls {

		data, err := json.Marshal(url)

		if err != nil {
			return err
		}

		_, err = writer.Write(data)

		if err != nil {
			return err
		}

		_, err = writer.WriteString("\n")

		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func (s *FileStorage) GetShortURL(original string) (bool, string) {

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, url := range s.urls {
		if url.Original == original {
			return true, url.Short
		}
	}

	return false, ""
}

func (s *FileStorage) GetOriginalURL(short string) (bool, string) {

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, url := range s.urls {
		if url.Short == short {
			return true, url.Original
		}
	}

	return false, ""
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}
