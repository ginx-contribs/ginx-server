package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/pkg/token"
	"github.com/ginx-contribs/ginx/contribs/ratelimit/counter"
	"github.com/pkg/errors"
	"time"
)

const AuthKey = "auth"

// Private metadata means that api needs to be user authenticated
var Private = ginx.V{Key: AuthKey, Val: 0}

// Public means that api no need to be authenticated
var Public = ginx.V{Key: AuthKey, Val: 1}

const CountKey = "count"

// CountLimit metadata means that api need to rate limit by number of requests
func CountLimit(limit int, duration time.Duration) ginx.V {
	return ginx.V{Key: CountKey, Val: counter.Limiter{Limit: limit, Window: duration}}
}

const tokenKey = "auth.token.context.info.key"

// SetTokenInfo stores token information into context
func SetTokenInfo(ctx *gin.Context, token *token.Token) {
	ctx.Set(tokenKey, token)
}

// GetTokenInfo returns token information from context
func GetTokenInfo(ctx *gin.Context) (*token.Token, error) {
	value, exists := ctx.Get(tokenKey)
	if !exists {
		return nil, errors.New("there is no token in context")
	}
	if gotToken, ok := value.(*token.Token); !ok {
		return nil, fmt.Errorf("expected %T, got %T", &token.Token{}, value)
	} else if gotToken == nil {
		return nil, errors.New("nil token in context")
	} else {
		return gotToken, nil
	}
}

// MustGetTokenInfo returns token information from context, panic if err != nil
func MustGetTokenInfo(ctx *gin.Context) *token.Token {
	info, err := GetTokenInfo(ctx)
	if err != nil {
		panic(err)
	}
	return info
}
