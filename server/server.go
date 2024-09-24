package server

import (
	"context"
	"fmt"
	_ "github.com/ginx-contribs/ent-sqlite"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/server/conf"
	"github.com/ginx-contribs/ginx-server/server/pkg/logh"
	"github.com/ginx-contribs/ginx-server/server/types"
	"log/slog"
)

// NewApp returns a new http app server
func NewApp(ctx context.Context, appConf *conf.App) (*ginx.Server, error) {
	// initialize database
	slog.Debug(fmt.Sprintf("connecting to %s(%s)", appConf.DB.Driver, appConf.DB.Address))
	db, err := NewDBClient(ctx, appConf.DB)
	if err != nil {
		return nil, err
	}

	// initialize redis client
	slog.Debug(fmt.Sprintf("connecting to redis(%s)", appConf.Redis.Address))
	redisClient, err := NewRedisClient(ctx, appConf.Redis)
	if err != nil {
		return nil, err
	}

	// initialize email client
	slog.Debug(fmt.Sprintf("establish email client(%s:%d)", appConf.Email.Host, appConf.Email.Port))
	emailClient, err := NewEmailClient(ctx, appConf.Email)
	if err != nil {
		return nil, err
	}

	tc := types.Context{
		AppConf: appConf,
		Ent:     db,
		Redis:   redisClient,
		Email:   emailClient,
	}

	// initialize ginx server
	server, err := NewHttpServer(ctx, appConf, tc)
	if err != nil {
		return nil, err
	}
	tc.Router = server.RouterGroup().Group(appConf.Server.BasePath)
	slog.Debug("setup api router")

	// initialize api router
	sc, err := setup(tc)
	if err != nil {
		return nil, err
	}

	queue := sc.MQ
	queue.Start(ctx)
	slog.Info("message queue is listening")

	// register cron job
	cronJob, err := NewCronJob(ctx, tc, sc)
	if err != nil {
		return nil, err
	}
	started := cronJob.Start()
	slog.Info(fmt.Sprintf("created %d cron jobs", started))

	// shutdown hook
	onShutdown := func(ctx context.Context) error {
		logh.ErrorNotNil("message queue closed failed", queue.Close())
		slog.Info(fmt.Sprintf("stopped %d jobs", cronJob.Stop()))
		// should close db and redis at the end
		logh.ErrorNotNil("db closed failed", db.Close())
		logh.ErrorNotNil("redis closed failed", redisClient.Close())
		return nil
	}
	server.OnShutdown = append(server.OnShutdown, onShutdown)

	return server, nil
}
