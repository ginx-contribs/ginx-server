package ts

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnix(t *testing.T) {
	unix := Unix()
	time := FromUnix(unix)
	if !assert.EqualValues(t, unix, time.Unix()) {
		return
	}
	t.Log(unix, time.Unix())
}

func TestUnixMicro(t *testing.T) {
	unix := UnixMicro()
	time := FromUnixMicro(unix)
	if !assert.EqualValues(t, unix, time.UnixMicro()) {
		return
	}
	t.Log(unix, time.Unix())
}
