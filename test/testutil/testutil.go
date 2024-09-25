package testutil

import (
	"context"
	"fmt"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/server"
	"github.com/ginx-contribs/ginx-server/server/conf"
	"github.com/ginx-contribs/ginx-server/server/svc"
	"github.com/ginx-contribs/ginx-server/server/types"
	"log/slog"
	"os"
	"time"
)

type TestServer struct {
	Tc      types.Context
	Sc      svc.Context
	Server  *ginx.Server
	Cleanup func()
}

// NewTestServer return a server for testing
func NewTestServer(ctx context.Context) (testServer *TestServer, err error) {
	appConf, err := ReadConf()
	if err != nil {
		return
	}
	// initialize database
	slog.Debug(fmt.Sprintf("connecting to %s(%s)", appConf.DB.Driver, appConf.DB.Address))
	db, err := server.NewDBClient(ctx, appConf.DB)
	if err != nil {
		return
	}

	// initialize redis client
	slog.Debug(fmt.Sprintf("connecting to redis(%s)", appConf.Redis.Address))
	redisClient, err := server.NewRedisClient(ctx, appConf.Redis)
	if err != nil {
		return
	}

	// initialize email client
	slog.Debug(fmt.Sprintf("establish email client(%s:%d)", appConf.Email.Host, appConf.Email.Port))
	emailClient, err := server.NewEmailClient(ctx, appConf.Email)
	if err != nil {
		return
	}

	tc := types.Context{
		AppConf: appConf,
		Ent:     db,
		Redis:   redisClient,
		Email:   emailClient,
	}

	// initialize ginx server
	svr, err := server.NewHttpServer(ctx, appConf, tc)
	if err != nil {
		return
	}
	tc.Router = svr.RouterGroup().Group(appConf.Server.BasePath)
	slog.Debug("setup api router")

	// initialize api router
	sc, err := server.Setup(tc)
	if err != nil {
		return
	}

	return &TestServer{
		Tc:     tc,
		Sc:     sc,
		Server: svr,
		Cleanup: func() {
			db.Close()
			redisClient.Close()
		},
	}, nil
}

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
}

const configFile = "testdata/conf.toml"

func ReadConf() (*conf.App, error) {
	appConf, err := conf.ReadFrom(configFile)
	if err != nil {
		return nil, err
	}
	return &appConf, err
}

// ReadDBConf returns the test configuration
func ReadDBConf() (conf.DB, error) {
	appConf, err := conf.ReadFrom(configFile)
	if err != nil {
		return conf.DB{}, err
	}
	return appConf.DB, err
}

func NewTimer() *Timer {
	return &Timer{}
}

// Timer is helper to calculate cost-time
type Timer struct {
	start time.Time
}

func (t *Timer) Start() {
	t.start = time.Now()
}

func (t *Timer) Stop() time.Duration {
	return time.Now().Sub(t.start)
}

func (t *Timer) Reset() {
	t.start = time.Time{}
}

func NewRound() *Round {
	return &Round{}
}

type Round struct {
	r int64
}

func (r *Round) Round() int64 {
	rr := r.r
	r.r++
	return rr
}

func (r *Round) Reset() {
	r.r = 0
}
