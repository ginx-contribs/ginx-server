package handler

import (
	"github.com/ginx-contribs/ginx-server/server/handler/auth"
	"github.com/ginx-contribs/ginx-server/server/handler/email"
	"github.com/ginx-contribs/ginx-server/server/handler/job"
	"github.com/ginx-contribs/ginx-server/server/handler/user"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	// auth handlers
	auth.NewAuthHandler,
	auth.NewTokenHandler,
	auth.NewVerifyCodeHandler,

	// email handlers
	email.NewEmailHandler,

	// user handlers
	user.NewUserHandler,

	// job
	job.NewCronJob,
	job.NewJobHandler,
)
