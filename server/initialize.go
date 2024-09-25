package server

import (
	entsql "entgo.io/ent/dialect/sql"
	"errors"
	"github.com/dstgo/size"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/dbx"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/server/conf"
	"github.com/ginx-contribs/ginx-server/server/data/ent"
	authhandler "github.com/ginx-contribs/ginx-server/server/handler/auth"
	"github.com/ginx-contribs/ginx-server/server/handler/job"
	"github.com/ginx-contribs/ginx-server/server/mids"
	"github.com/ginx-contribs/ginx-server/server/svc"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/ginx-contribs/ginx/constant/methods"
	"github.com/ginx-contribs/ginx/contribs/requestid"
	"github.com/ginx-contribs/ginx/middleware"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/ginx-contribs/logx"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/wneessen/go-mail"
	"golang.org/x/net/context"
	"log/slog"
	"net/http/pprof"
	"strings"
	"time"

	// offline time zone database
	_ "time/tzdata"
)

// ContextProvider only use for wire injection
var ContextProvider = wire.NewSet(
	wire.FieldsOf(new(types.Context), "AppConf"),
	wire.FieldsOf(new(types.Context), "Ent"),
	wire.FieldsOf(new(types.Context), "Redis"),
	wire.FieldsOf(new(types.Context), "Router"),
	wire.FieldsOf(new(types.Context), "Email"),
	wire.FieldsOf(new(*conf.App), "Jwt"),
	wire.FieldsOf(new(*conf.App), "Email"),
)

// NewLogger returns a new app logger with the given options
func NewLogger(option conf.Log) (*logx.Logger, error) {

	writer, err := logx.NewWriter(&logx.WriterOptions{
		Filename: option.Filename,
	})
	if err != nil {
		return nil, err
	}
	handler, err := logx.NewHandler(writer, &logx.HandlerOptions{
		Level:       option.Level,
		Format:      option.Format,
		Prompt:      option.Prompt,
		Source:      option.Source,
		ReplaceAttr: nil,
		Color:       option.Color,
	})
	if err != nil {
		return nil, err
	}
	logger, err := logx.New(
		logx.WithHandlers(handler),
	)
	if err != nil {
		return nil, err
	}
	return logger, nil
}

func NewHttpServer(ctx context.Context, appConf *conf.App, tc types.Context) (*ginx.Server, error) {
	server := ginx.New(
		ginx.WithOptions(ginx.Options{
			Mode:               gin.ReleaseMode,
			Address:            appConf.Server.Address,
			ReadTimeout:        appConf.Server.ReadTimeout.Duration(),
			WriteTimeout:       appConf.Server.WriteTimeout.Duration(),
			IdleTimeout:        appConf.Server.IdleTimeout.Duration(),
			MaxMultipartMemory: appConf.Server.MultipartMax,
			MaxHeaderBytes:     int(size.MB * 2),
			MaxShutdownTimeout: time.Second * 5,
		}),
		// 404 handler
		ginx.WithNoRoute(middleware.NoRoute()),
		// 405 handler
		ginx.WithNoMethod(middleware.NoMethod(methods.Get, methods.Post, methods.Put, methods.Delete, methods.Options)),
		// global middlewares
		ginx.WithMiddlewares(
			// recovery handler
			middleware.Recovery(slog.Default(), nil),
			// request id
			requestid.RequestId(),
			// access logger
			middleware.Logger(slog.Default(), "request-log"),
			// jwt authentication
			mids.TokenAuthenticator(authhandler.NewTokenHandler(appConf.Jwt, tc.Redis)),
		),
	)

	// set validator for gin
	err := setupHumanizedValidator()
	if err != nil {
		return nil, err
	}

	// whether to enable pprof program profiling
	if appConf.Server.Pprof {
		server.Engine().GET("/pprof/profile", gin.WrapF(pprof.Profile))
		server.Engine().GET("/pprof/heap", gin.WrapH(pprof.Handler("heap")))
		server.Engine().GET("/pprof/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		slog.Info("pprof profiling enabled")
	}

	return server, nil
}

// override the default ginx validation error handler, see ginx.SetValidateHandler
func setupHumanizedValidator() error {
	v := validator.New()
	v.SetTagName("binding")
	englishValidator, err := ginx.EnglishValidator(v, validateParams)
	if err != nil {
		return err
	}
	ginx.SetValidator(englishValidator)
	ginx.SetValidateHandler(englishValidator.HandleError)
	return nil
}

func validateParams(ctx *gin.Context, val any, err error, translator ut.Translator) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var errorMsg []string
		for _, validateErr := range validationErrors {
			errorMsg = append(errorMsg, validateErr.Translate(translator))
		}
		verr := errors.New(strings.Join(errorMsg, ", "))
		resp.Fail(ctx).Code(types.ErrBadParams.Code).Error(verr).JSON()
		return
	}
	resp.Fail(ctx).Code(types.ErrBadParams.Code).ErrorMsg("params validate failed").JSON()
}

// NewDBClient initialize database with ent
func NewDBClient(ctx context.Context, dbConf conf.DB) (*ent.Client, error) {
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

// NewEmailClient initialize email client
func NewEmailClient(ctx context.Context, emailConf conf.Email) (*mail.Client, error) {
	client, err := mail.NewClient(emailConf.Host,
		mail.WithPort(emailConf.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(emailConf.Username),
		mail.WithPassword(emailConf.Password),
	)
	if err != nil {
		return nil, err
	}
	// test if smtp server is available
	err = client.DialWithContext(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// NewCronJob initialize cron jobs
func NewCronJob(ctx context.Context, tc types.Context, sc svc.Context) (*job.CronJob, error) {
	cj := sc.CronJob
	// hooks
	cj.BeforeHooks = append(cj.BeforeHooks, job.LogBefore(), job.UpdateBefore(sc.JobHandler))
	cj.AfterHooks = append(cj.AfterHooks, job.LogAfter())
	errs := []error{}
	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	for _, j := range cj.FutureJobs() {
		err := sc.JobHandler.Upsert(ctx, j)
		if err != nil {
			return nil, err
		}
	}
	return cj, nil
}
