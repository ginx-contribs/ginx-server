package logh

import "log/slog"

// ErrorNotNil log error if the error is not nil
func ErrorNotNil(msg string, err error) {
	if err != nil {
		slog.Error(msg, slog.Any("error", err))
	}
}
