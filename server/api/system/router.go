package system

import (
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/server/types/route"
)

type Router struct {
	System *SystemAPI
}

func NewRouter(root *ginx.RouterGroup, systemAPI *SystemAPI) Router {
	// test api
	root.MGET("/ping", ginx.M{route.Public}, systemAPI.Ping)
	root.MGET("/pong", ginx.M{route.Private}, systemAPI.Pong)

	return Router{System: systemAPI}
}
