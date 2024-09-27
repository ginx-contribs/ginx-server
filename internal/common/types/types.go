package types

import (
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/ent"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx-server/pkg/mq"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/wneessen/go-mail"
)

var Provider = wire.NewSet(
	wire.FieldsOf(new(Injector), "Config"),
	wire.FieldsOf(new(Injector), "Router"),
	wire.FieldsOf(new(Injector), "EntDB"),
	wire.FieldsOf(new(Injector), "Redis"),
	wire.FieldsOf(new(Injector), "Email"),
	wire.FieldsOf(new(Injector), "MQ"),
	wire.FieldsOf(new(*conf.App), "Jwt"),
	wire.FieldsOf(new(*conf.App), "Email"),
)

// Injector holds all needed object for initializing app
type Injector struct {
	// app configuration
	Config *conf.App
	// root router for http internal
	Router *ginx.RouterGroup
	// ent db client
	EntDB *ent.Client
	// redis client
	Redis *redis.Client
	// email client
	Email *mail.Client
	// message queue
	MQ mq.Queue
}
