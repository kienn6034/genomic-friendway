package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	storages := map[string]Storage{
		"memory": NewMemoryStorage(),
	}

	for name, storage := range storages {
		t.Run(name, func(t *testing.T) {
			// Test data
			testData := []byte("test gene data")

			// Test Store
			hash, err := storage.Store(testData)
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)

			// Test Retrieve
			retrieved, err := storage.Retrieve(hash)
			assert.NoError(t, err)
			assert.Equal(t, testData, retrieved)

			// Test Delete
			err = storage.Delete(hash)
			assert.NoError(t, err)

			// Verify deletion
			_, err = storage.Retrieve(hash)
			assert.Error(t, err)
		})
	}
}
