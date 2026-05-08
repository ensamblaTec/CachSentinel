package domain

import "time"

type CacheEntry[T any] struct {
	Value      T         `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	HitCount   int64     `json:"hit_count"`
	LastAccess time.Time `json:"last_access"`
}

type Config struct {
	DefaultTTL       time.Duration
	PredictiveWindow time.Duration
	HitThreshold     int64
}
