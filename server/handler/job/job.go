package job

import (
	"errors"
	"fmt"
	"github.com/ginx-contribs/ginx-server/test/testutil"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/robfig/cron/v3"
	"io"
	"log/slog"
	"maps"
	"slices"
	"sync/atomic"
)

// Job is representation of a cron jobs
type Job interface {
	// Name returns the name of the jobs
	Name() string
	// Cron returns the cron expression for the jobs
	Cron() string
	// Cmd returns the command to be executed by the jobs
	Cmd() func() (attrs []any, err error)
}

// FutureJob is representation of a future cron jobs
type FutureJob struct {
	Job
	cron.Entry
}

func NewCronJob() *CronJob {
	c := cron.New(cron.WithLogger(discardLogger{out: io.Discard}))
	cj := &CronJob{
		cron:   c,
		record: cmap.New[Job](),
		future: cmap.New[FutureJob](),
	}
	return cj
}

// CronJob is a simple cron jobs manager
type CronJob struct {
	cron   *cron.Cron
	record cmap.ConcurrentMap[string, Job]
	future cmap.ConcurrentMap[string, FutureJob]

	BeforeHooks []BeforeHook
	AfterHooks  []AfterHook
}

// AddJob add a new job to the future list
func (c *CronJob) AddJob(job Job) error {
	_, e := c.GetJob(job.Name())
	if e {
		return errors.New("same job already")
	}
	entryId, err := c.cron.AddJob(job.Cron(), c.newWrapper(job))
	if err != nil {
		return err
	}
	c.record.Set(job.Name(), job)
	c.future.Set(job.Name(), FutureJob{Entry: cron.Entry{ID: entryId}, Job: job})
	return nil
}

// GetJob return a job from the future list
func (c *CronJob) GetJob(name string) (FutureJob, bool) {
	job, e := c.future.Get(name)
	if !e {
		return FutureJob{}, false
	}
	entry := c.cron.Entry(job.ID)
	if (entry != cron.Entry{}) {
		job.Entry = entry
	}
	return job, true
}

func (c *CronJob) FutureJobs() []FutureJob {
	return slices.Collect(maps.Values(c.future.Items()))
}

// DelJob remove a job from the future list, but it will not remove from the record
func (c *CronJob) DelJob(name string) {
	job, e := c.GetJob(name)
	if !e {
		return
	}
	c.cron.Remove(job.ID)
	c.future.Remove(name)
	slog.Info(fmt.Sprintf("remove job %v from future", name))
}

// ContinueJob put a job existing in the record into the future list
func (c *CronJob) ContinueJob(name string) error {
	job, e := c.record.Get(name)
	if !e {
		return errors.New("job not found")
	}
	entryId, err := c.cron.AddJob(job.Cron(), c.newWrapper(job))
	if err != nil {
		return err
	}
	c.future.Set(job.Name(), FutureJob{Entry: cron.Entry{ID: entryId}, Job: job})
	slog.Info(fmt.Sprintf("continue job %v to future", name))
	return nil
}

// Start starts the cron schedule running
func (c *CronJob) Start() int {
	c.cron.Start()
	return len(c.cron.Entries())
}

// Stop stops and waits for the all scheduled jobs to complete
func (c *CronJob) Stop() int {
	live := len(c.cron.Entries())
	c.cron.Stop()
	return live
}

func (c *CronJob) newWrapper(job Job) *JobWrapper {
	return &JobWrapper{c: c, job: job}
}

// JobWrapper is a wrapper for Job
type JobWrapper struct {
	c     *CronJob
	job   Job
	round atomic.Int64
}

func (r *JobWrapper) Run() {
	// it must be existing in there
	job, _ := r.c.GetJob(r.job.Name())

	// execute hooks before cmd
	for _, hook := range r.c.BeforeHooks {
		hook(job, r.round.Load())
	}

	timer := testutil.NewTimer()
	timer.Start()

	// execute command
	attrs, err := r.job.Cmd()()

	elapsed := timer.Stop()

	// execute hooks after cmd
	for _, hook := range r.c.AfterHooks {
		hook(job, r.round.Load(), elapsed, err, attrs...)
	}

	r.round.Add(1)
}

type discardLogger struct {
	out    io.Writer
	prefab string
}

func (c discardLogger) Info(msg string, keysAndValues ...interface{}) {}

func (c discardLogger) Error(err error, msg string, keysAndValues ...interface{}) {}
