package token

import (
	"context"
	"github.com/jellydator/ttlcache/v2"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// Cache is responsible for persist jwt token.
type Cache interface {
	// Get returns the given tokenId from the cache, just verifying if it exists
	Get(ctx context.Context, prefix, tokenId string) (string, bool, error)
	// TTL returns the rest life-time of the given token
	TTL(ctx context.Context, prefix, tokenId string) (time.Duration, bool, error)
	// Del removes the specified tokenId from the cache
	Del(ctx context.Context, prefix, tokenId string) error
	// Set sets the specified tokenId to the cache
	Set(ctx context.Context, prefix, tokenId, value string, expire time.Duration) error
}

func prefixKey(prefix, key string) string {
	return prefix + ":" + key
}

func NewRedisTokenCache(client *redis.Client) *RedisTokenCache {
	return &RedisTokenCache{redis: client}
}

// RedisTokenCache implements TokenCache interface for redis storage
type RedisTokenCache struct {
	redis *redis.Client
}

func (r *RedisTokenCache) Get(ctx context.Context, prefix, tokenId string) (string, bool, error) {
	result := r.redis.Get(ctx, prefixKey(prefix, tokenId))
	if errors.Is(result.Err(), redis.Nil) {
		return result.String(), false, nil
	}
	return result.String(), result.Err() == nil, nil
}

func (r *RedisTokenCache) TTL(ctx context.Context, prefix, tokenId string) (time.Duration, bool, error) {
	result := r.redis.TTL(ctx, prefixKey(prefix, tokenId))
	if errors.Is(result.Err(), redis.Nil) {
		return result.Val(), false, nil
	}
	return result.Val(), result.Err() == nil, result.Err()
}

func (r *RedisTokenCache) Set(ctx context.Context, prefix, tokenId, value string, expire time.Duration) error {
	result := r.redis.Set(ctx, prefixKey(prefix, tokenId), value, expire)
	return result.Err()
}

func (r *RedisTokenCache) Del(ctx context.Context, prefix, tokenId string) error {
	// del tokenId from string
	result := r.redis.Del(ctx, prefixKey(prefix, tokenId))
	if err := result.Err(); errors.Is(err, redis.Nil) {
		return nil
	}
	return result.Err()
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		memStore: ttlcache.NewCache(),
	}
}

// MemoryCache implement Cache by ttlcache.Cache in memory
type MemoryCache struct {
	memStore *ttlcache.Cache
}

func (m *MemoryCache) Get(ctx context.Context, prefix, tokenId string) (string, bool, error) {
	key := prefixKey(prefix, tokenId)
	value, err := m.memStore.Get(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return "", false, nil
	}
	return value.(string), err == nil, err
}

func (m *MemoryCache) TTL(ctx context.Context, prefix, tokenId string) (time.Duration, bool, error) {
	key := prefixKey(prefix, tokenId)
	_, duration, err := m.memStore.GetWithTTL(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return duration, false, nil
	}
	return duration, err == nil, err
}

func (m *MemoryCache) Del(ctx context.Context, prefix, tokenId string) error {
	key := prefixKey(prefix, tokenId)
	err := m.memStore.Remove(key)
	if errors.Is(err, ttlcache.ErrNotFound) {
		return nil
	}
	return err
}

func (m *MemoryCache) Set(ctx context.Context, prefix, tokenId, value string, expire time.Duration) error {
	key := prefixKey(prefix, tokenId)
	return m.memStore.SetWithTTL(key, value, expire)
}
