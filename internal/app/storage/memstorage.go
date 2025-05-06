package storage

import (
	"sync"
)

type MemStorage struct {
	urls []URL
	mu   sync.Mutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{urls: make([]URL, 0)}
}

func (s *MemStorage) Add(u URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls = append(s.urls, u)

	return nil
}

func (s *MemStorage) Batch(urls []URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls = append(s.urls, urls...)

	return nil
}

func (s *MemStorage) GetShortURL(original string) (bool, string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, url := range s.urls {
		if url.Original == original {
			return true, url.Short
		}
	}

	return false, ""
}

func (s *MemStorage) GetOriginalURL(short string) (bool, string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, url := range s.urls {
		if url.Short == short {
			return true, url.Original
		}
	}

	return false, ""
}

func (s *MemStorage) Close() error {
	return nil
}
