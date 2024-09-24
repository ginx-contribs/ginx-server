package job

import (
	"context"
	"log/slog"
	"time"
)

type BeforeHook func(job FutureJob, round int64)

type AfterHook func(job FutureJob, round int64, elapsed time.Duration, err error, attrs ...any)

// LogBefore log job information before job execution
func LogBefore() BeforeHook {
	return func(job FutureJob, round int64) {
		slog.Info("cron job prepared",
			slog.Int64("round", round),
			slog.String("name", job.Name()),
			slog.Int("entry", int(job.ID)),
			slog.Time("prev", job.Prev),
			slog.Time("next", job.Next))
	}
}

// LogAfter log job information before job execution
func LogAfter() AfterHook {
	return func(job FutureJob, round int64, elapsed time.Duration, err error, attrs ...any) {
		baseAttrs := []any{
			slog.String("name", job.Name()),
			slog.Int64("round", round),
			slog.Duration("elapsed", elapsed),
			slog.Int("entry", int(job.ID)),
			slog.Time("prev", job.Prev),
			slog.Time("next", job.Next),
		}
		if err != nil {
			slog.Error("cron job failed", append(baseAttrs, slog.Any("error", err))...)
		} else {
			slog.Info("cron job finished", append(baseAttrs, slog.Group("result", attrs...))...)
		}
	}
}

// UpdateBefore update job information in the db
func UpdateBefore(h *Handler) BeforeHook {
	return func(job FutureJob, round int64) {
		err := h.Upsert(context.Background(), job)
		if err != nil {
			slog.Error("update job failed", slog.Any("error", err))
		}
	}
}
