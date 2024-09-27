package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/handler"
	types2 "github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx/pkg/resp"
)

type UserAPI struct {
	UserHandler handler.UserHandler
}

// Me
// @Summary      Me
// @Description  return user information for current user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.Response{data=types.UserInfo}
// @Router       /user/me [GET]
func (u UserAPI) Me(ctx *gin.Context) {
	tokenInfo := types2.MustGetTokenInfo(ctx)
	userInfo, err := u.UserHandler.FindByUID(ctx, tokenInfo.Claims.UserId)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Data(userInfo).JSON()
	}
}

// Info
// @Summary      Info
// @Description  get user information by given uid
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        uid  query  string  true "uid"
// @Success      200  {object}  types.Response{data=types.UserInfo}
// @Router       /user/info [GET]
func (u UserAPI) Info(ctx *gin.Context) {
	var opt types2.UidOptions
	if err := ginx.ShouldValidateQuery(ctx, &opt); err != nil {
		return
	}
	userInfo, err := u.UserHandler.FindByUID(ctx, opt.Uid)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Data(userInfo).JSON()
	}
}

// List
// @Summary      List
// @Description  list user info by page
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        SearchUserOptions   query   types.SearchUserOptions  true  "SearchUserOptions"
// @Success      200  {object}  types.Response{data=types.UserSearchResult}
// @Router       /user/list [GET]
func (u UserAPI) List(ctx *gin.Context) {
	var page types2.SearchUserOptions
	if err := ginx.ShouldValidateQuery(ctx, &page); err != nil {
		return
	}

	userInfoList, err := u.UserHandler.ListUserByPage(ctx, page.Page, page.Size, page.Search)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Data(userInfoList).JSON()
	}
}
