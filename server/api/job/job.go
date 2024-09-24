package job

import (
	"github.com/gin-gonic/gin"
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx-server/server/handler/job"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/ginx-contribs/ginx/pkg/resp"
)

func NewJobAPI(jHandler *job.Handler) *JobAPI {
	return &JobAPI{jHandler: jHandler}
}

type JobAPI struct {
	jHandler *job.Handler
}

// List
// @Summary      List
// @Description  list jobs by page
// @Tags         job
// @Accept       json
// @Produce      json
// @Param        JobPageOption  query  types.JobPageOption  true "JobPageOption"
// @Success      200  {object}  types.Response{data=types.JobPageList}
// @Router       /job/list [GET]
func (j *JobAPI) List(ctx *gin.Context) {
	var opts types.JobPageOption
	if err := ginx.ShouldValidateQuery(ctx, &opts); err != nil {
		return
	}

	result, err := j.jHandler.List(ctx, opts.Page, opts.Size, opts.Search)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Data(result).JSON()
	}
}

// Info
// @Summary      Info
// @Description  get job info
// @Tags         job
// @Accept       json
// @Produce      json
// @Param        JobNameOptions  query  types.JobNameOptions true "JobNameOptions"
// @Success      200  {object}  types.Response{data=types.JobInfo}
// @Router       /job/info [GET]
func (j *JobAPI) Info(ctx *gin.Context) {
	var opts types.JobNameOptions
	if err := ginx.ShouldValidateQuery(ctx, &opts); err != nil {
		return
	}

	one, err := j.jHandler.GetOne(ctx, opts.Name)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).Data(one).JSON()
	}
}

// Start
// @Summary      Start
// @Description  start the job
// @Tags         job
// @Accept       json
// @Produce      json
// @Param        JobNameOptions  query  types.JobNameOptions true "JobNameOptions"
// @Success      200  {object}  types.Response
// @Router       /job/start [POST]
func (j *JobAPI) Start(ctx *gin.Context) {
	var opts types.JobNameOptions
	if err := ginx.ShouldValidateQuery(ctx, &opts); err != nil {
		return
	}
	err := j.jHandler.Start(ctx, opts.Name)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).JSON()
	}
}

// Stop
// @Summary      Stop
// @Description  stop the job
// @Tags         job
// @Accept       json
// @Produce      json
// @Param        JobNameOptions  query  types.JobNameOptions true "JobNameOptions"
// @Success      200  {object}  types.Response
// @Router       /job/stop [POST]
func (j *JobAPI) Stop(ctx *gin.Context) {
	var opts types.JobNameOptions
	if err := ginx.ShouldValidateQuery(ctx, &opts); err != nil {
		return
	}
	err := j.jHandler.Stop(ctx, opts.Name)
	if err != nil {
		resp.Fail(ctx).Error(err).JSON()
	} else {
		resp.Ok(ctx).JSON()
	}
}
