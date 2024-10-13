package wirex

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx-server/internal/common/types"
	"github.com/ginx-contribs/ginx-server/pkg/mids"
	"github.com/ginx-contribs/ginx/contribs/requestid"
	"github.com/ginx-contribs/ginx/middleware"
	"log/slog"
	"time"
)

// GlobalMiddlewares initialize all needed global middlewares, order is important.
func GlobalMiddlewares(injector types.Injector) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		Recovery(),
		RequestID(),
		AccessLogger(),
		TokenVerify(injector),
		RequestCache(injector),
	}
}

// Recovery return recovery middleware
func Recovery() gin.HandlerFunc {
	return middleware.Recovery(slog.Default(), nil)
}

// RequestID returns request-id middleware
func RequestID() gin.HandlerFunc {
	return requestid.RequestId()
}

// AccessLogger return access logger middleware
func AccessLogger() gin.HandlerFunc {
	return middleware.Logger(slog.Default(), "request-log")
}

// RequestCache return Cache middleware
func RequestCache(injector types.Injector) gin.HandlerFunc {
	return middleware.CacheRedis("ginx-cache.", time.Second*2, injector.Redis)
}

// TokenVerify return jwt token authenticate middleware
func TokenVerify(injector types.Injector) gin.HandlerFunc {
	return mids.TokenAuthenticator(injector.Token.VerifyAccess)
}
