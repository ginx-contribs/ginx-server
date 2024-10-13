package logh

import "log/slog"

// NoError log error if the error is not nil
func NoError(msg string, err error) {
	if err != nil {
		slog.Error(msg, slog.Any("error", err))
	}
}
