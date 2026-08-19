package logging

import "log/slog"

func Error(err error) slog.Attr {
	return slog.Any("error", err)
}

func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}
