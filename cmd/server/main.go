package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ensamblatec/CachSentinel/internal/core/domain"
	"github.com/ensamblatec/CachSentinel/internal/core/service"
	"github.com/ensamblatec/CachSentinel/internal/infrastructure/adapter"
	"github.com/ensamblatec/CachSentinel/internal/infrastructure/api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := domain.Config{
		DefaultTTL:   1 * time.Minute,
		HitThreshold: 5,
	}

	repo := adapter.NewMemoryStore[any]()
	fetcher := &adapter.HTTPFetcher{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Client:  &http.Client{Timeout: 10 * time.Second},
	}

	cacheSvc := service.NewCacheService[any](repo, fetcher, cfg)
	handler := api.NewProxyHandler(cacheSvc)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	portServer := 8080
	go func() {
		logger.Info("server_started", "port", portServer)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server_failed", "err", err)
			os.Exit(1)
		}
	}()

	<-done
	logger.Info("server_stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("shutdown_failed", "err", err)
	}
	logger.Info("server_exited_cleanly")
}
