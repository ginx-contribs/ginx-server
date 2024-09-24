package types

import (
	"github.com/ginx-contribs/ginx-server/server/data/ent"
	"github.com/ginx-contribs/ginx/constant/status"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
)

var (
	ErrJobNotFound = statuserr.Errorf("cron job not found").SetCode(4_400_001).SetStatus(status.BadRequest)
)

type JobNameOptions struct {
	Name string `form:"name" binding:"required"`
}

type JobPageOption struct {
	Page   int    `form:"page"`
	Size   int    `form:"size"`
	Search string `form:"search"`
}

type JobInfo struct {
	Name  string `json:"name"`
	Entry int    `json:"entry"`
	Cron  string `json:"cron"`
	Next  int64  `json:"next"`
	Prev  int64  `json:"prev"`
}

type JobPageList struct {
	Total int       `json:"total"`
	List  []JobInfo `json:"list"`
}

func EntJobToJobInfo(j *ent.CronJob) JobInfo {
	return JobInfo{
		Entry: j.EntryID,
		Name:  j.Name,
		Cron:  j.Cron,
		Next:  j.Next,
		Prev:  j.Prev,
	}
}

func EntJobToJobInfoBatch(js []*ent.CronJob) []JobInfo {
	var jis []JobInfo
	for _, j := range js {
		jis = append(jis, EntJobToJobInfo(j))
	}
	return jis
}
