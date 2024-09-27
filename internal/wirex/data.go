package wirex

import (
	"context"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/ginx-contribs/dbx"
	"github.com/ginx-contribs/ginx-server/ent"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/logx"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

// NewEntDB initialize database with ent
func NewEntDB(ctx context.Context, dbConf conf.DB, logger *logx.Logger) (*ent.Client, error) {
	sqldb, err := dbx.Open(dbx.Options{
		Driver:             dbConf.Driver,
		Address:            dbConf.Address,
		User:               dbConf.User,
		Password:           dbConf.Password,
		Database:           dbConf.Database,
		Params:             dbConf.Params,
		MaxIdleConnections: dbConf.MaxIdleConnections,
		MaxOpenConnections: dbConf.MaxOpenConnections,
		MaxLifeTime:        dbConf.MaxLifeTime.Duration(),
		MaxIdleTime:        dbConf.MaxIdleTime.Duration(),
	})
	if err != nil {
		return nil, err
	}
	entClient := ent.NewClient(
		ent.Log(logger.LogLogger(slog.LevelDebug).Println),
		ent.Driver(entsql.OpenDB(dbConf.Driver, sqldb)),
	)
	// migrate database
	if err := entClient.Schema.Create(ctx); err != nil {
		return nil, err
	}

	return entClient, err
}

// NewRedisClient initialize redis connection
func NewRedisClient(ctx context.Context, redisConf conf.Redis) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisConf.Address,
		Password:     redisConf.Password,
		ReadTimeout:  redisConf.ReadTimeout.Duration(),
		WriteTimeout: redisConf.WriteTimeout.Duration(),
	})
	pingResult := redisClient.Ping(ctx)
	if pingResult.Err() != nil {
		return nil, pingResult.Err()
	}
	return redisClient, nil
}
