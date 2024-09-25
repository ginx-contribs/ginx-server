package data

import (
	"github.com/ginx-contribs/ginx-server/server/data/cache"
	"github.com/ginx-contribs/ginx-server/server/data/mq"
	"github.com/ginx-contribs/ginx-server/server/data/repo"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	// cache
	cache.NewRedisTokenCache,
	wire.Bind(new(cache.TokenCache), new(*cache.RedisTokenCache)),
	cache.NewRedisCodeCache,
	wire.Bind(new(cache.VerifyCodeCache), new(*cache.RedisCodeCache)),

	// user
	repo.NewUserRepo,

	// mq
	mq.NewStreamQueue,
	wire.Bind(new(mq.Queue), new(*mq.StreamQueue)),
)
