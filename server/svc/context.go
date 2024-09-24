package svc

import (
	"github.com/ginx-contribs/ginx-server/server/api"
	"github.com/ginx-contribs/ginx-server/server/data/mq"
	"github.com/ginx-contribs/ginx-server/server/data/repo"
	"github.com/ginx-contribs/ginx-server/server/handler/auth"
	"github.com/ginx-contribs/ginx-server/server/handler/email"
	"github.com/ginx-contribs/ginx-server/server/handler/job"
	"github.com/ginx-contribs/ginx-server/server/handler/user"
	"github.com/google/wire"
)

var Provider = wire.NewSet(
	wire.Struct(new(Context), "*"),
)

// Context holds all handler and repo instances, just for helper
type Context struct {
	// api
	ApiRouter api.Router

	// user
	UserHandler *user.UserHandler
	UserRepo    *repo.UserRepo

	// system
	AuthHandler *auth.AuthHandler

	// email
	EmailHandler *email.Handler

	// job
	JobHandler *job.Handler
	JobRepo    *repo.JobRepo
	CronJob    *job.CronJob

	// message queue
	MQ mq.Queue
}
