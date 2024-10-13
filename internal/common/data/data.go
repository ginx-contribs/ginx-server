package data

import (
	"github.com/ginx-contribs/ginx-server/internal/common/data/cache"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	cache.NewRedisCaptchaCache,
	wire.Bind(new(cache.CaptchaCache), new(*cache.RedisCodeCache)),
)
