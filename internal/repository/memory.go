package repository

import (
	"errors"
	"sync"
)

var (
	ErrNotFound      = errors.New("url não encontrada")
	ErrAlreadyExists = errors.New("código já existe")
)

type MemoryRepository struct {
	mu   sync.RWMutex
	urls map[string]string // código -> url original
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		urls: make(map[string]string),
	}
}

func (r *MemoryRepository) Save(code, url string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.urls[code]; exists {
		return ErrAlreadyExists
	}
	r.urls[code] = url
	return nil
}

func (r *MemoryRepository) Find(code string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	url, ok := r.urls[code]
	if !ok {
		return "", ErrNotFound
	}
	return url, nil
}