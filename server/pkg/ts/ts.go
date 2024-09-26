package ts

import "time"

// local time zone
var local = time.Local

// Location return local time zone
func Location() *time.Location {
	if local == nil {
		return time.Local
	}
	return local
}

// In updates local time zone
func In(loc *time.Location) {
	local = loc
}

// Zero return zero value for time.Time
func Zero() time.Time {
	return time.Time{}
}

// Now returns not time
func Now() time.Time {
	return time.Now().In(Location())
}

// Unix returns the unix timestamp for now
func Unix() int64 {
	return Now().Unix()
}

// FromUnix return time.Time from unix timestamp
func FromUnix(ts int64) time.Time {
	return time.Unix(ts, 0).In(Location())
}

// UnixMicro returns the unix micro timestamp for
func UnixMicro() int64 {
	return Now().UnixMicro()
}

// FromUnixMicro returns time.Time from the unix timestamp
func FromUnixMicro(ts int64) time.Time {
	return time.UnixMicro(ts).In(Location())
}
