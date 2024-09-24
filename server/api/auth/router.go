package auth

import (
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/server/types/route"
)

type Router struct {
	auth *AuthAPI
}

func NewRouter(group *ginx.RouterGroup, auth *AuthAPI) Router {

	// auth
	authGroup := group.MGroup("/auth", ginx.M{route.Public})
	authGroup.POST("/login", auth.Login)
	authGroup.POST("/register", auth.Register)
	authGroup.POST("/reset", auth.ResetPassword)
	authGroup.POST("/refresh", auth.Refresh)
	authGroup.POST("/code", auth.VerifyCode)

	return Router{auth: auth}
}
