package server

import (
	"context"
	"fmt"
	_ "github.com/ginx-contribs/ent-sqlite"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/internal/common/types"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx-server/internal/modules"
	"github.com/ginx-contribs/ginx-server/internal/wirex"
	"github.com/ginx-contribs/ginx-server/pkg/logh"
	"github.com/ginx-contribs/ginx-server/pkg/mq"
	"github.com/ginx-contribs/logx"
	"log/slog"

	// offline time zone database
	_ "time/tzdata"
)

// NewHTTPServer returns a new http app server
func NewHTTPServer(ctx context.Context, appConf *conf.App, logger *logx.Logger) (*ginx.Server, error) {
	// initialize database
	slog.Debug(fmt.Sprintf("connecting to %s(%s)", appConf.DB.Driver, appConf.DB.Address))
	db, err := wirex.NewEntDB(ctx, appConf.DB, logger)
	if err != nil {
		return nil, err
	}
	// initialize redis client
	slog.Debug(fmt.Sprintf("connecting to redis(%s)", appConf.Redis.Address))
	redisClient, err := wirex.NewRedisClient(ctx, appConf.Redis)
	if err != nil {
		return nil, err
	}
	// initialize email client
	slog.Debug(fmt.Sprintf("establish email client(%s:%d)", appConf.Email.Host, appConf.Email.Port))
	emailClient, err := wirex.NewEmailSender(ctx, appConf.Email)
	if err != nil {
		return nil, err
	}
	// initialize token resolver
	tokenResolver, err := wirex.NewTokenResolver(ctx, appConf.Jwt, redisClient)
	if err != nil {
		return nil, err
	}
	// initialize message queue
	queue := mq.NewStreamQueue(ctx, redisClient)
	// build injector
	injector := types.Injector{
		Config: appConf,
		EntDB:  db,
		Redis:  redisClient,
		Token:  tokenResolver,
		Email:  emailClient,
		MQ:     queue,
	}
	// initialize ginx server
	server, err := wirex.NewHttpServer(ctx, appConf, injector)
	if err != nil {
		return nil, err
	}
	injector.Router = server.RouterGroup().Group(appConf.Server.BasePath)
	slog.Debug("setup api router")

	// inject modules
	mods, err := wirex.Inject(injector)
	if err != nil {
		return nil, err
	}
	// initialize modules
	modManager := modules.NewModuleManager(&mods)
	err = modManager.Init(injector)
	if err != nil {
		return nil, err
	}

	// hooks before start
	onStart := func(ctx context.Context) error {
		queue.Start(ctx)
		slog.Info("message queue is listening")
		return nil
	}
	server.BeforeStarting = append(server.BeforeStarting, onStart)
	// hooks before shutdown
	onShutdown := func(ctx context.Context) error {
		logh.NoError("modules closed failed", modManager.Close())
		// datasource should be closed at last
		logh.NoError("message queue closed failed", queue.Close())
		logh.NoError("db closed failed", db.Close())
		logh.NoError("redis closed failed", redisClient.Close())
		return nil
	}
	server.OnShutdown = append(server.OnShutdown, onShutdown)

	return server, nil
}
