package wirex

import (
	"context"
	"github.com/dstgo/size"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/internal/common/types"
	"github.com/ginx-contribs/ginx-server/internal/conf"
	"github.com/ginx-contribs/ginx/constant/methods"
	"github.com/ginx-contribs/ginx/middleware"
	"github.com/ginx-contribs/ginx/pkg/resp"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"log/slog"
	"net/http/pprof"
	"strings"
	"time"
)

// NewHttpServer return new http server with given configuration
func NewHttpServer(ctx context.Context, appConf *conf.App, injector types.Injector) (*ginx.Server, error) {
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
		ginx.WithMiddlewares(GlobalMiddlewares(injector)...),
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
