package storage

import (
	"sync"
)

type MemoryStorage struct {
	links map[string]string
	mu    sync.Mutex
}

func NewMemoryStorage() MemoryStorage {
	return MemoryStorage{links: make(map[string]string)}
}

func (s *MemoryStorage) Get(key string) (string, bool) {
	s.mu.Lock()
	value, exist := s.links[key]
	s.mu.Unlock()

	return value, exist
}

func (s *MemoryStorage) Add(key string, value string) error {
	s.mu.Lock()
	s.links[key] = value
	s.mu.Unlock()

	return nil
}

func (s *MemoryStorage) ContainsValue(value string) (bool, string) {
	for key, url := range s.links {
		if url == value {
			return true, key
		}
	}

	return false, ""
}
