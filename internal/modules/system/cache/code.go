package cache

import (
	"errors"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"time"
)

var (
	ErrCodeRepeated = errors.New("code repeated")
)

// CaptchaCache is responsible for storing captcha code into cache
type CaptchaCache interface {
	// Set storage code into cache with ttl, return false if retry failed.
	Set(ctx context.Context, usage types.Usage, code, to string, ttl, retry time.Duration) (bool, error)
	// Get returns the specified code
	Get(ctx context.Context, usage types.Usage, code string) (string, error)
	// Del remove the specified code from cache
	Del(ctx context.Context, usage types.Usage, code string) error
}

var _ CaptchaCache = (*RedisCaptchaCache)(nil)

func NewRedisCaptchaCache(cache *redis.Client) *RedisCaptchaCache {
	return &RedisCaptchaCache{cache: cache}
}

// RedisCaptchaCache implements CaptchaCache with redis cache
type RedisCaptchaCache struct {
	cache *redis.Client
}

func (r *RedisCaptchaCache) Set(ctx context.Context, usage types.Usage, code, to string, ttl, retry time.Duration) (bool, error) {

	// check retry ttl
	retryRes, err := r.cache.Get(ctx, to).Result()
	if !errors.Is(err, redis.Nil) && err != nil {
		return false, statuserr.InternalError(err)
	} else if err == nil && retryRes != "" {
		return false, nil
	}

	// check if is repeated
	get, err := r.Get(ctx, usage, code)
	if !errors.Is(err, redis.Nil) && err != nil {
		return false, err
	} else if get != "" {
		return false, ErrCodeRepeated
	}

	codeKey := usage.Name() + ":" + code

	// set verify code
	if _, err := r.cache.Set(ctx, codeKey, to, ttl).Result(); err != nil {
		return false, err
	}

	// set retry ttl
	if err := r.cache.Set(ctx, to, to, retry).Err(); err != nil {
		return false, err
	}

	return true, nil
}

func (r *RedisCaptchaCache) Get(ctx context.Context, usage types.Usage, code string) (string, error) {
	codeKey := usage.Name() + ":" + code
	return r.cache.Get(ctx, codeKey).Result()
}

func (r *RedisCaptchaCache) Del(ctx context.Context, usage types.Usage, code string) error {
	codeKey := usage.Name() + ":" + code
	return r.cache.Del(ctx, codeKey).Err()
}
