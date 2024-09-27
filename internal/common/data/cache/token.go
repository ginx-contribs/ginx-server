package cache

import (
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

// TokenCache stores the token in the cache, and maintains them.
type TokenCache interface {
	// Get returns the given tokenId from the cache, just verifying if it exists
	Get(ctx context.Context, tokenId string) (string, error)
	// TTL returns the rest life-time of the given token
	TTL(ctx context.Context, tokenId string) (time.Duration, error)
	// Del removes the specified tokenId from the cache
	Del(ctx context.Context, tokenId string) error
	// Set sets the specified tokenId to the cache
	Set(ctx context.Context, tokenId, value string, expire time.Duration) error
}

func NewRedisTokenCache(prefix string, client *redis.Client) *RedisTokenCache {
	return &RedisTokenCache{prefix: prefix, redis: client}
}

// RedisTokenCache implements TokenCache interface for redis storage
type RedisTokenCache struct {
	prefix string
	redis  *redis.Client
}

func (r *RedisTokenCache) prefixKey(key string) string {
	return r.prefix + ":" + key
}

func (r *RedisTokenCache) Get(ctx context.Context, tokenId string) (string, error) {
	result := r.redis.Get(ctx, r.prefixKey(tokenId))
	return result.String(), result.Err()
}

func (r *RedisTokenCache) TTL(ctx context.Context, tokenId string) (time.Duration, error) {
	result := r.redis.TTL(ctx, r.prefixKey(tokenId))
	return result.Val(), result.Err()
}

func (r *RedisTokenCache) Del(ctx context.Context, tokenId string) error {
	// del tokenId from string
	result := r.redis.Del(ctx, r.prefixKey(tokenId))
	if err := result.Err(); errors.Is(err, redis.Nil) {
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func (r *RedisTokenCache) Set(ctx context.Context, tokenId, value string, expire time.Duration) error {
	result := r.redis.Set(ctx, r.prefixKey(tokenId), value, expire)
	return result.Err()
}
