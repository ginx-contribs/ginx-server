package testutil

import (
	"github.com/ginx-contribs/ginx-server/server/conf"
	"log/slog"
	"os"
	"time"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
}

const configFile = "testdata/conf.toml"

func ReadConf() (conf.App, error) {
	appConf, err := conf.ReadFrom(configFile)
	if err != nil {
		return conf.App{}, err
	}
	return appConf, err
}

// ReadDBConf returns the test configuration
func ReadDBConf() (conf.DB, error) {
	appConf, err := conf.ReadFrom(configFile)
	if err != nil {
		return conf.DB{}, err
	}
	return appConf.DB, err
}

func NewTimer() *Timer {
	return &Timer{}
}

// Timer is helper to calculate cost-time
type Timer struct {
	start time.Time
}

func (t *Timer) Start() {
	t.start = time.Now()
}

func (t *Timer) Stop() time.Duration {
	return time.Now().Sub(t.start)
}

func (t *Timer) Reset() {
	t.start = time.Time{}
}

func NewRound() *Round {
	return &Round{}
}

type Round struct {
	r int64
}

func (r *Round) Round() int64 {
	rr := r.r
	r.r++
	return rr
}

func (r *Round) Reset() {
	r.r = 0
}
