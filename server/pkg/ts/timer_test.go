package ts

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimer_Start_Reset(t *testing.T) {
	timer := NewTimer()
	timer.Begin()
	startAt := timer.startAt
	time.Sleep(time.Second)
	timer.Begin()
	if !assert.EqualValues(t, startAt, timer.startAt) {
		return
	}
	timer.Reset()
	if !assert.EqualValues(t, startAt, timer.startAt) {
		return
	}
	t.Log(startAt, timer.startAt)
}

func TestTimer_Stop_Reset(t *testing.T) {
	timer := NewTimer()
	timer.Begin()
	time.Sleep(time.Second)
	timer.Stop()
	stopAt := timer.stopAt
	time.Sleep(time.Millisecond * 500)
	timer.Stop()
	if !assert.EqualValues(t, stopAt, timer.stopAt) {
		return
	}
	timer.Reset()
	if !assert.NotEqualValues(t, stopAt, timer.stopAt) {
		return
	}
	t.Log(stopAt, timer.stopAt)
}

func TestTimer_Duration(t *testing.T) {
	timer := NewTimer()
	timer.Begin()
	time.Sleep(time.Second)
	timer.Stop()
	if !assert.EqualValues(t, timer.Duration(), timer.stopAt.Sub(timer.startAt)) {
		return
	}
	t.Log(timer.Duration(), timer.stopAt.Sub(timer.startAt))
}
