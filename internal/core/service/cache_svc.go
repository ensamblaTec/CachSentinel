package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/ensamblatec/CachSentinel/internal/core/domain"
	"github.com/ensamblatec/CachSentinel/internal/core/ports"
	"golang.org/x/sync/singleflight"
)

type CacheService struct {
	repo    ports.CacheRepository
	fetcher ports.UpstreamFetcher
	cfg     domain.Config
	sf      singleflight.Group
	logger  *slog.Logger
}

func NewCacheService(repo ports.CacheRepository, fetcher ports.UpstreamFetcher, cfg domain.Config) *CacheService {
	return &CacheService{
		repo:    repo,
		fetcher: fetcher,
		cfg:     cfg,
		logger:  slog.Default(),
	}
}

func (service *CacheService) Execute(ctx context.Context, key string) ([]byte, error) {
	entry, err := service.repo.Get(ctx, key)

	if err == nil && entry != nil {
		if time.Now().Before(entry.ExpiresAt) {
			service.evaluatePrediction(key, entry)
			return entry.Value, nil
		}

		service.logger.Info("cache_expired", "key", key)
		if err := service.repo.Delete(ctx, key); err != nil {
			service.logger.Error("cache_delete_error", "key", key, "err", err)
		}
	}

	val, err, _ := service.sf.Do(key, func() (any, error) {
		service.logger.Info("cache_miss", "key", key)
		data, fetchErr := service.refreshAndStore(ctx, key)

		if fetchErr != nil {
			if entry != nil {
				service.logger.Error("upstream_fetch_error_with_stale_cache", "key", key, "err", fetchErr)
				entry.ExpiresAt = time.Now().Add(30 * time.Second)
				_ = service.repo.Set(ctx, key, entry)
				return entry.Value, nil
			}
			return nil, fetchErr
		}
		return data, nil
	})

	if err != nil {
		return nil, err
	}

	return val.([]byte), nil
}

func (service *CacheService) evaluatePrediction(key string, entry *domain.CacheEntry) {
	timeLeft := time.Until(entry.ExpiresAt)

	if timeLeft < (service.cfg.DefaultTTL/5) && entry.HitCount > service.cfg.HitThreshold {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err, _ := service.sf.Do(key, func() (any, error) {
				return service.refreshAndStore(ctx, key)
			})

			if err != nil {
				service.logger.Error("predictive_refresh_error", "key", key, "err", err)
			} else {
				service.logger.Info("predictive_refresh_success", "key", key)
			}
		}()
	}
}

func (service *CacheService) refreshAndStore(ctx context.Context, key string) ([]byte, error) {
	data, err := service.fetcher.Fetch(ctx, key)
	if err != nil {
		return nil, err
	}

	entry := &domain.CacheEntry{
		Value:     data,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(service.cfg.DefaultTTL),
		HitCount:  1,
	}

	if err := service.repo.Set(ctx, key, entry); err != nil {
		service.logger.Error("store_error", "err", err)
	}

	return data, nil
}
