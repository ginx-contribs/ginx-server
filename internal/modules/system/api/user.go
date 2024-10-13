package api

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/handler"
	"github.com/ginx-contribs/ginx-server/internal/modules/system/types"
	"github.com/ginx-contribs/ginx-server/pkg/toolset/ginxutils"
	"github.com/ginx-contribs/ginx/pkg/resp"
)

type UserAPI struct {
	UserHandler handler.UserHandler
}

// Profile
// @Summary      Profile
// @Description  return user information for current user
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.Response{data=types.UserInfo}
// @Router       /user/profile [GET]
func (u UserAPI) Profile(ctx *gin.Context) {
	token, ok := ginxutils.GetLoginUserToken(ctx)
	if !ok {
		return
	}
	uid := token.Claims.Payload.(types.TokenPayload).UserId
	userInfo, err := u.UserHandler.FindByUID(ctx, uid)
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
// @Router       /user/:uid [GET]
func (u UserAPI) Info(ctx *gin.Context) {
	var opt types.UidOptions
	if err := ginx.ShouldValidateURI(ctx, &opt); err != nil {
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
// @Router       /users [GET]
func (u UserAPI) List(ctx *gin.Context) {
	var page types.SearchUserOptions
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
