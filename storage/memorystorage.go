package storage

import (
	"context"
	"sync"
)

// limit sets the limit to 100 hashes
var limit int = 1000

// deleteSize defines how many hashes should be deleted after the limit is reached
var deleteSize int = 100

// NewMemoryStorage returnes a new memorystorage
func NewMemoryStorage() *memorystorage {
	list := []string{}
	return &memorystorage{list: list}
}

type memorystorage struct {
	list []string
	mux  sync.Mutex
}

// Set stores the given hash
func (ms *memorystorage) Set(ctx context.Context, hash string) error {
	ms.mux.Lock()
	defer ms.mux.Unlock()

	if len(ms.list) >= limit {
		ms.list = append(ms.list[deleteSize:], hash)
	} else {
		ms.list = append(ms.list, hash)
	}

	return nil
}

// Exists checks if the given hash exists
func (ms *memorystorage) Exists(ctx context.Context, hash string) (bool, error) {
	ms.mux.Lock()
	defer ms.mux.Unlock()

	exists := false
	for _, h := range ms.list {
		if h == hash {
			exists = true
		}
	}

	if exists {
		return true, nil
	}

	return false, nil
}
