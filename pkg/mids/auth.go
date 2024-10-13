package mids

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/internal/common/route"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx-server/pkg/token"
	"github.com/ginx-contribs/ginx/constant/headers"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"strings"
)

// TokenAuthenticator authenticates each request if is valid.
func TokenAuthenticator(verify func(ctx context.Context, token string) (token.Token, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// check if is public api
		metadata := ginx.MetaFromCtx(ctx)
		if !metadata.Contains(route.Private) {
			ctx.Next()
			return
		}

		// parse token string from header
		header := ctx.Request.Header.Get(headers.Authorization)
		pair := strings.Split(header, " ")
		if len(pair) != 2 || pair[0] != "Bearer" {
			resp.Fail(ctx).Status(status.Unauthorized).Error(types.ErrCredentialInvalid).JSON()
			ctx.Abort()
			return
		}
		tokenString := pair[1]

		// verify token if is valid
		tokenInfo, err := verify(ctx, tokenString)
		if err == nil {
			// stores token info into context
			route.SetTokenInfo(ctx, &tokenInfo)
			ctx.Next()
		} else {
			ctx.Abort()
			// check if is needed to refresh
			if errors.Is(err, token.ErrTokenNeedsRefresh) {
				resp.Fail(ctx).Error(types.ErrTokenNeedsRefresh).JSON()
			} else if errors.Is(err, token.ErrAccessTokenExpired) {
				resp.Fail(ctx).Error(types.ErrCredentialExpired).JSON()
			} else { // invalid token
				ctx.Error(err)
				resp.Fail(ctx).Error(types.ErrCredentialInvalid).JSON()
			}
		}
	}
}
