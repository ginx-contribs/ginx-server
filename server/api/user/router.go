package user

import (
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/server/types/route"
)

type Router struct {
	User *UserAPI
}

func NewRouter(root *ginx.RouterGroup, userApi *UserAPI) Router {

	userGroup := root.MGroup("/user", ginx.M{route.Private})
	userGroup.MGET("/me", ginx.M{route.Private}, userApi.Me)
	userGroup.MGET("/info", ginx.M{route.Private}, userApi.Info)
	userGroup.MGET("/list", ginx.M{route.Private}, userApi.List)

	return Router{User: userApi}
}
