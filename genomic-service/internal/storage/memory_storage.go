package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

// MemoryStorage implements Storage using RAM
type MemoryStorage struct {
	mu    sync.RWMutex
	store map[string][]byte
}

// NewMemoryStorage creates a new in-memory storage
func NewMemoryStorage() Storage {
	return &MemoryStorage{
		store: make(map[string][]byte),
	}
}

// Memory Storage Implementation
func (ms *MemoryStorage) Store(data []byte) (string, error) {
	// Generate hash for the data
	hash := sha256.Sum256(data)
	fileHash := hex.EncodeToString(hash[:])

	// Store data in memory
	ms.mu.Lock()
	ms.store[fileHash] = make([]byte, len(data))
	copy(ms.store[fileHash], data) // Create a copy of the data to prevent external modifications
	ms.mu.Unlock()

	return fileHash, nil
}

func (ms *MemoryStorage) Retrieve(fileHash string) ([]byte, error) {
	ms.mu.RLock()
	data, exists := ms.store[fileHash]
	ms.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("file not found: %s", fileHash)
	}

	// Return a copy of the data to prevent external modifications
	result := make([]byte, len(data))
	copy(result, data)
	return result, nil
}

func (ms *MemoryStorage) Delete(fileHash string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.store[fileHash]; !exists {
		return fmt.Errorf("file not found: %s", fileHash)
	}

	delete(ms.store, fileHash)
	return nil
}
