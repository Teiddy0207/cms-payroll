package inmemcache

import (
	"cal-salary/core/logger"
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
)

type InMemoryCache struct {
	client *bigcache.BigCache
}

func NewInMemoryCache() *InMemoryCache {
	client, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		logger.Error("Failed to create in-memory cache: " + err.Error())
		return nil
	}
	return &InMemoryCache{
		client: client,
	}
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value []byte) error {
	return c.client.Set(key, value)
}

func (c *InMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.client.Get(key)
}
