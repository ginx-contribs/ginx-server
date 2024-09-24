package ts

import "time"

var Loc = time.Local

func In(loc *time.Location) {
	Loc = loc
}

// Now returns not time
func Now() time.Time {
	return time.Now().In(Loc)
}

// UnixMicro returns the unix timestamp representation of now
func UnixMicro() int64 {
	return Now().UnixMicro()
}

// FromUnixMicro returns *time.Time from the unix timestamp
func FromUnixMicro(ts int64) time.Time {
	return time.UnixMicro(ts).In(Loc)
}
