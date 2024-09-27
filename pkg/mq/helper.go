package mq

import "context"

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

type steps []func() (error, bool)

func (s *steps) Then(step func() (error, bool)) *steps {
	*s = append(*s, step)
	return s
}
