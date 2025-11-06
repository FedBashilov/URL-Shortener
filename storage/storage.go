package storage

import (
	"fmt"
	"sync"
)

type URLStorage struct {
	mu   sync.RWMutex
	urls map[string]string
}

var Storage = &URLStorage{
	urls: make(map[string]string),
}

func (s *URLStorage) Save(longURL string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	shortCode := fmt.Sprintf("url%d", len(s.urls)+1)
	s.urls[shortCode] = longURL

	return shortCode
}

func (s *URLStorage) Get(shortCode string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	longURL, exists := s.urls[shortCode]
	return longURL, exists
}
