package api

import (
	"github.com/ginx-contribs/ginx-server/server/api/auth"
	"github.com/ginx-contribs/ginx-server/server/api/job"
	"github.com/ginx-contribs/ginx-server/server/api/system"
	"github.com/ginx-contribs/ginx-server/server/api/user"
	"github.com/google/wire"
)

// RegisterRouter
// @title	                        Lobby HTTP API
// @version		                    v0.0.0-Beta
// @description                     This is swagger generated api documentation, know more information about lobby on GitHub.
// @contact.name                    ginx-contribs
// @contact.url                     https://github.com/ginx-contribs/ginx-server
// @BasePath	                    /api/
// @license.name                    MIT LICENSE
// @license.url                     https://mit-license.org/
// @securityDefinitions.apikey      BearerAuth
// @in                              header
// @name                            Authorization
//
//go:generate swag init --ot yaml --generatedTime -g api.go -d ./,../types,../pkg --output ./ && swag fmt -g api.go -d ./

type Router struct {
	Auth   auth.Router
	System system.Router
	User   user.Router
	Job    job.Router
}

var Provider = wire.NewSet(
	// auth router
	auth.NewAuthAPI,
	auth.NewRouter,
	// system router
	system.NewSystemAPI,
	system.NewRouter,
	// user router
	user.NewUserAPI,
	user.NewRouter,

	// job router
	job.NewJobAPI,
	job.NewRouter,

	// build Router struct
	wire.Struct(new(Router), "*"),
)
