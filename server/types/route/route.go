package route

import (
	"github.com/ginx-contribs/ginx"
	"github.com/ginx-contribs/ginx/contribs/ratelimit/counter"
	"time"
)

const AuthKey = "auth"

// Private metadata means that api needs to be user authenticated
var Private = ginx.V{Key: AuthKey, Val: 0}

// Public means that api no need to be authenticated
var Public = ginx.V{Key: AuthKey, Val: 1}

const CountKey = "count"

// CountLimit metadata means that api need to rate limit by number of requests
func CountLimit(limit int, duration time.Duration) ginx.V {
	return ginx.V{Key: CountKey, Val: counter.Limiter{Limit: limit, Window: duration}}
}
