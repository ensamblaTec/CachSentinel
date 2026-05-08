package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/ensamblatec/CachSentinel/internal/core/domain"
	"github.com/ensamblatec/CachSentinel/internal/core/ports"
	"golang.org/x/sync/singleflight"
)

type CacheService[T any] struct {
	repo    ports.CacheRepository[T]
	fetcher ports.UpstreamFetcher[T]
	cfg     domain.Config
	sf      singleflight.Group
	logger  *slog.Logger
}

func NewCacheService[T any](repo ports.CacheRepository[T], fetcher ports.UpstreamFetcher[T], cfg domain.Config) *CacheService[T] {
	return &CacheService[T]{
		repo:    repo,
		fetcher: fetcher,
		cfg:     cfg,
		logger:  slog.Default(),
	}
}

func (service *CacheService[T]) Execute(ctx context.Context, key string) (T, error) {
	entry, err := service.repo.Get(ctx, key)

	if err == nil && entry != nil {
		service.evaluatePrediction(key, entry)
		return entry.Value, nil
	}

	val, err, _ := service.sf.Do(key, func() (any, error) {
		return service.refreshAndStore(ctx, key)
	})

	if err != nil {
		return *new(T), err
	}

	return val.(T), nil
}

func (service *CacheService[T]) evaluatePrediction(key string, entry *domain.CacheEntry[T]) {
	timeLeft := time.Until(entry.ExpiresAt)

	if timeLeft < (service.cfg.DefaultTTL/5) && entry.HitCount > service.cfg.HitThreshold {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			service.logger.Info("predictive_refresh_triggered", "key", key)
			service.refreshAndStore(ctx, key)
		}()
	}
}

func (service *CacheService[T]) refreshAndStore(ctx context.Context, key string) (T, error) {
	data, err := service.fetcher.Fetch(ctx, key)
	if err != nil {
		return *new(T), err
	}

	entry := &domain.CacheEntry[T]{
		Value:     data,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(service.cfg.DefaultTTL),
	}

	if err := service.repo.Set(ctx, key, entry); err != nil {
		service.logger.Error("store_error", "err", err)
	}

	return data, nil
}
