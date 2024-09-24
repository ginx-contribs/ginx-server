package job

import (
	"context"
	"errors"
	"github.com/ginx-contribs/ginx-server/server/data/ent"
	"github.com/ginx-contribs/ginx-server/server/data/repo"
	"github.com/ginx-contribs/ginx-server/server/types"
	"github.com/ginx-contribs/ginx/pkg/resp/statuserr"
	"github.com/robfig/cron/v3"
)

func NewJobHandler(jobrepo *repo.JobRepo, cronjob *CronJob) *Handler {
	return &Handler{jobRepo: jobrepo, cronjob: cronjob}
}

type Handler struct {
	jobRepo *repo.JobRepo
	cronjob *CronJob
}

// List returns a list of jobs by page
func (h *Handler) List(ctx context.Context, page int, size int, search string) (types.JobPageList, error) {
	list, total, err := h.jobRepo.ListByPage(ctx, page, size, search)
	if err != nil {
		return types.JobPageList{}, statuserr.InternalError(err)
	}
	var jobs []types.JobInfo
	for _, j := range list {
		entry := h.cronjob.cron.Entry(cron.EntryID(j.EntryID))
		info := types.JobInfo{
			Entry: j.EntryID,
			Name:  j.Name,
			Cron:  j.Cron,
			Prev:  j.Prev,
			Next:  j.Next,
		}
		// if the entry is already in the future list
		if (entry != cron.Entry{}) {
			// revise Prev Next timestamp
			if entry.Prev.UnixMicro() < 0 {
				info.Prev = j.Prev
			} else {
				info.Prev = entry.Prev.UnixMicro()
			}
			if entry.Next.UnixMicro() < 0 {
				info.Next = j.Next
			} else {
				info.Next = entry.Next.UnixMicro()
			}
		}
		jobs = append(jobs, info)
	}
	return types.JobPageList{
		Total: total,
		List:  jobs,
	}, err
}

func (h *Handler) Upsert(ctx context.Context, job FutureJob) error {
	return h.jobRepo.UpsertOne(ctx, &ent.CronJob{
		Name:    job.Name(),
		Cron:    job.Cron(),
		EntryID: int(job.ID),
		Prev:    job.Prev.UnixMicro(),
		Next:    job.Next.UnixMicro(),
	})
}

func (h *Handler) GetOne(ctx context.Context, name string) (types.JobInfo, error) {
	one, err := h.jobRepo.QueryOne(ctx, name)
	if ent.IsNotFound(err) {
		return types.JobInfo{}, types.ErrJobNotFound
	} else if err != nil {
		return types.JobInfo{}, statuserr.InternalError(err)
	}
	entry := h.cronjob.cron.Entry(cron.EntryID(one.EntryID))
	info := types.JobInfo{
		Name:  one.Name,
		Entry: one.EntryID,
		Cron:  one.Cron,
		Next:  one.Next,
		Prev:  one.Prev,
	}
	// if the entry is already in the future list
	if (entry != cron.Entry{}) {
		if entry.Prev.UnixMicro() < 0 {
			info.Prev = one.Prev
		} else {
			info.Prev = entry.Prev.UnixMicro()
		}
		if entry.Next.UnixMicro() < 0 {
			info.Next = one.Next
		} else {
			info.Next = entry.Next.UnixMicro()
		}
	}
	return info, nil
}

// Stop removes the job from the future scheduled jobs list
func (h *Handler) Stop(ctx context.Context, name string) error {
	job, e := h.cronjob.GetJob(name)
	if !e {
		return errors.New("job not found")
	}
	// remove the job from the future scheduler
	h.cronjob.DelJob(name)
	// update information
	err := h.jobRepo.UpsertOne(ctx, &ent.CronJob{
		Cron: job.Cron(),
		// set the entryId = -1 to indicate this job is stopped
		EntryID: -1,
		Prev:    job.Prev.UnixMicro(),
		Next:    -1,
	})
	if err != nil {
		return err
	}
	return nil
}

// Start add job into future scheduled jobs list
func (h *Handler) Start(ctx context.Context, name string) error {
	err := h.cronjob.ContinueJob(name)
	if err != nil {
		return err
	}
	job, _ := h.cronjob.GetJob(name)

	// update information
	err = h.jobRepo.UpsertOne(ctx, &ent.CronJob{
		Cron:    job.Cron(),
		EntryID: int(job.ID),
		Prev:    job.Prev.UnixMicro(),
		Next:    job.Next.UnixMicro(),
	})
	if err != nil {
		return err
	}
	return nil
}
