package testutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRound(t *testing.T) {
	r := &Round{}
	for i := range 10 {
		assert.EqualValues(t, i, r.Round())
	}
}
