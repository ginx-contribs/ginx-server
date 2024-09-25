package user

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/server/handler/user"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/ginx-contribs/ginx/pkg/resp"
)

func NewUserAPI(userHandler *user.UserHandler) *UserAPI {
	return &UserAPI{userHandler: userHandler}
}

type UserAPI struct {
	userHandler *user.UserHandler
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
	tokenInfo := types.MustGetTokenInfo(ctx)
	userInfo, err := u.userHandler.FindByUID(ctx, tokenInfo.Claims.UserId)
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
	var opt types.ULIDOptions
	if err := ginx.ShouldValidateQuery(ctx, &opt); err != nil {
		return
	}
	userInfo, err := u.userHandler.FindByUID(ctx, opt.Uid)
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
// @Param        UserSearchOption   query   types.UserSearchOption  true  "UserSearchOption"
// @Success      200  {object}  types.Response{data=types.UserSearchResult}
// @Router       /user/list [GET]
func (u UserAPI) List(ctx *gin.Context) {
	var page types.UserSearchOption
	if err := ginx.ShouldValidateQuery(ctx, &page); err != nil {
		return
	}

	userInfoList, err := u.userHandler.ListUserByPage(ctx, page.Page, page.Size, page.Search)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Data(userInfoList).JSON()
	}
}
