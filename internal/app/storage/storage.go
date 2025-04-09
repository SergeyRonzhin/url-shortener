package storage

import (
	"sync"
)

type Storage struct {
	links map[string]string
	mu    sync.Mutex
}

func New() Storage {
	return Storage{links: make(map[string]string)}
}

func (s *Storage) Get(key string) (string, bool) {
	s.mu.Lock()
	value, exist := s.links[key]
	s.mu.Unlock()

	return value, exist
}

func (s *Storage) Add(key string, value string) {
	s.mu.Lock()
	s.links[key] = value
	s.mu.Unlock()
}

func (s *Storage) ContainsValue(value string) (bool, string) {
	for key, url := range s.links {
		if url == value {
			return true, key
		}
	}

	return false, ""
}
