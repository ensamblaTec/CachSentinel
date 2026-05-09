package ports

import (
	"context"

	"github.com/ensamblatec/CachSentinel/internal/core/domain"
)

type CacheRepository interface {
	Get(ctx context.Context, key string) (*domain.CacheEntry, error)
	Set(ctx context.Context, key string, entry *domain.CacheEntry) error
	Delete(ctx context.Context, key string) error
}

type UpstreamFetcher interface {
	Fetch(ctx context.Context, key string) ([]byte, error)
}
