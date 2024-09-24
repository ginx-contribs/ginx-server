package repo

import (
	"context"
	"github.com/ginx-contribs/ginx-server/server/data/ent"
	"github.com/ginx-contribs/ginx-server/server/data/ent/cronjob"
)

func NewJobRepo(client *ent.Client) *JobRepo {
	return &JobRepo{Ent: client}
}

type JobRepo struct {
	Ent *ent.Client
}

func (j *JobRepo) Clear(ctx context.Context) (int, error) {
	return j.Ent.CronJob.Delete().Exec(ctx)
}

// UpsertOne creates a new Job if it is not existing, otherwise it wille update it.
func (j *JobRepo) UpsertOne(ctx context.Context, job *ent.CronJob) error {
	first, err := j.Ent.CronJob.Query().Where(cronjob.Name(job.Name)).First(ctx)
	if ent.IsNotFound(err) {
		_, err := j.Ent.CronJob.Create().SetCronJob(job).Save(ctx)
		return err
	} else if err != nil {
		return err
	}
	_, err = j.Ent.CronJob.UpdateOneID(first.ID).
		SetEntryID(job.EntryID).
		SetPrev(job.Prev).
		SetNext(job.Next).Save(ctx)
	if err != nil {
		return err
	}
	return nil
}

// QueryOne returns a job with the given name
func (j *JobRepo) QueryOne(ctx context.Context, name string) (*ent.CronJob, error) {
	return j.Ent.CronJob.Query().Where(cronjob.Name(name)).First(ctx)
}

func (j *JobRepo) FindByEntryId(ctx context.Context, ids ...int) ([]*ent.CronJob, error) {
	return j.Ent.CronJob.Query().Where(cronjob.EntryIDIn(ids...)).All(ctx)
}

// ListByPage returns a list of jobs
func (j *JobRepo) ListByPage(ctx context.Context, page int, size int, search string) ([]*ent.CronJob, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	query := j.Ent.CronJob.Query().
		Offset((page - 1) * size).
		Limit(size)

	if search != "" {
		query = query.Where(cronjob.NameEqualFold(search))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	list, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
