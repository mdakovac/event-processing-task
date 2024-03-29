package cache_service

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache = cache.Cache

func CreateCache(defaultExpiration time.Duration, cleanupInterval time.Duration) *Cache {
	c := cache.New(defaultExpiration, cleanupInterval)
	return c
}
