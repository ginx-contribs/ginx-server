package ids

import (
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
)

// UUID creates a new random UUID and returns it as a string or panics.
func UUID() string {
	return uuid.NewString()
}

// ULID returns an ULID with the current time in Unix milliseconds and
// monotonically increasing entropy for the same millisecond.
func ULID() string {
	return ulid.Make().String()
}
