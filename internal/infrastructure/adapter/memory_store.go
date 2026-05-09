package adapter

import (
	"context"
	"errors"
	"sync"

	"github.com/ensamblatec/CachSentinel/internal/core/domain"
)

type MemoryStore struct {
	data sync.Map
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (memory *MemoryStore) Get(ctx context.Context, key string) (*domain.CacheEntry, error) {
	val, ok := memory.data.Load(key)
	if !ok {
		return nil, errors.New("key_not_found")
	}

	entry := val.(*domain.CacheEntry)
	entry.HitCount++
	return entry, nil
}

func (memory *MemoryStore) Set(ctx context.Context, key string, entry *domain.CacheEntry) error {
	memory.data.Store(key, entry)
	return nil
}

func (memory *MemoryStore) Delete(ctx context.Context, key string) error {
	memory.data.Delete(key)
	return nil
}
