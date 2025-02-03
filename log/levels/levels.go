package levels

import "log/slog"

type Level slog.Level

const (
	Debug = Level(slog.LevelDebug)
	Info  = Level(slog.LevelInfo)
	Warn  = Level(slog.LevelWarn)
	Error = Level(slog.LevelError)
)
