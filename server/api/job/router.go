package job

import (
	"github.com/ginx-contribs/ginx"
)

type Router struct {
	Job *JobAPI
}

func NewRouter(root *ginx.RouterGroup, job *JobAPI) Router {
	group := root.Group("/job")
	group.GET("/info", job.Info)
	group.GET("/list", job.List)
	group.POST("/start", job.Start)
	group.POST("/stop", job.Stop)

	return Router{Job: job}
}
