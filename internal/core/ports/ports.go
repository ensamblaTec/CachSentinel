package ports

import (
	"context"

	"github.com/ensamblatec/CachSentinel/internal/core/domain"
)

type CacheRepository[T any] interface {
	Get(ctx context.Context, key string) (*domain.CacheEntry[T], error)
	Set(ctx context.Context, key string, entry *domain.CacheEntry[T]) error
	Delete(ctx context.Context, key string) error
}

type UpstreamFetcher[T any] interface {
	Fetch(ctx context.Context, key string) (T, error)
}
