package data

import (
	"github.com/ginx-contribs/ginx-server/internal/common/data/cache"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	// cache
	cache.NewRedisTokenCache,
	wire.Bind(new(cache.TokenCache), new(*cache.RedisTokenCache)),
	cache.NewRedisCaptchaCache,
	wire.Bind(new(cache.CaptchaCache), new(*cache.RedisCodeCache)),
)
