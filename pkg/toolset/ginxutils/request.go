package ginxutils

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx-server/internal/common/route"
	systype "github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx-server/pkg/token"
	"github.com/ginx-contribs/ginx/pkg/resp"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
)

// GetLoginUserToken return current user token
func GetLoginUserToken(ctx *gin.Context) (*token.Token, bool) {
	tokenInfo, e, err := route.GetTokenInfo(ctx)
	if !e {
		resp.Fail(ctx).Error(systype.ErrCredentialInvalid).JSON()
		return nil, false
	}
	if err != nil {
		resp.Fail(ctx).Error(statuserr.InternalError(err)).JSON()
		return nil, false
	}
	return tokenInfo, true
}
