package adapter

import (
	"context"
	"errors"
	"sync"

	"github.com/ensamblatec/CachSentinel/internal/core/domain"
)

type MemoryStore[T any] struct {
	data sync.Map
}

func NewMemoryStore[T any]() *MemoryStore[T] {
	return &MemoryStore[T]{}
}

func (memory *MemoryStore[T]) Get(ctx context.Context, key string) (*domain.CacheEntry[T], error) {
	val, ok := memory.data.Load(key)
	if !ok {
		return nil, errors.New("key_not_found")
	}

	entry := val.(*domain.CacheEntry[T])
	entry.HitCount++
	return entry, nil
}

func (memory *MemoryStore[T]) Set(ctx context.Context, key string, entry *domain.CacheEntry[T]) error {
	memory.data.Store(key, entry)
	return nil
}

func (memory *MemoryStore[T]) Delete(ctx context.Context, key string) error {
	memory.data.Delete(key)
	return nil
}
