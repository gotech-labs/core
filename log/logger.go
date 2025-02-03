package log

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/gotech-labs/core/log/formats"
	"github.com/gotech-labs/core/log/levels"
	"github.com/gotech-labs/core/system"
)

var (
	globalLogger = New(os.Stdout)
)

func New(out io.Writer, options ...option) Logger {
	// init settings
	settings := &settings{
		out:                  out,
		format:               formats.JSON,
		level:                levels.Debug,
		withAttrsFunc:        func() map[string]any { return nil },
		withContextAttrsFunc: func(context.Context) map[string]any { return nil },
		pretty:               true,
	}
	// bind options
	for _, option := range options {
		option(settings)
	}
	// create slog handler
	handler := settings.format.NewHandler(out, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.Level(settings.level),
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			switch attr.Key {
			case slog.TimeKey:
				attr.Value = slog.TimeValue(system.CurrentTime())
			}
			return attr
		},
	})
	// setup global attributes
	keyValues := settings.withAttrsFunc()
	attrs := make([]any, 0, len(keyValues))
	for k, v := range keyValues {
		attrs = append(attrs, k, v)
	}
	// create logger
	return &logger{
		internal: slog.New(handler).With(attrs...),
		settings: settings,
	}
}

type Logger interface {
	Debug(msg string, args ...any)
	DebugWithContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoWithContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnWithContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorWithContext(ctx context.Context, msg string, args ...any)
}

type logger struct {
	internal *slog.Logger
	settings *settings
}

func (l *logger) Debug(msg string, args ...any) {
	l.log(context.Background(), slog.LevelDebug, msg, args...)
}

func (l *logger) DebugWithContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelDebug, msg, args...)
}

func (l *logger) Info(msg string, args ...any) {
	l.log(context.Background(), slog.LevelInfo, msg, args...)
}

func (l *logger) InfoWithContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelInfo, msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.log(context.Background(), slog.LevelWarn, msg, args...)
}

func (l *logger) WarnWithContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelWarn, msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.log(context.Background(), slog.LevelError, msg, args...)
}

func (l *logger) ErrorWithContext(ctx context.Context, msg string, args ...any) {
	l.log(ctx, slog.LevelError, msg, args...)
}

func (l *logger) log(ctx context.Context, level slog.Level, msg string, args ...any) {
	if ctx != context.Background() && l.settings.withAttrsFunc != nil {
		attrs := l.settings.withContextAttrsFunc(ctx)
		if len(attrs) > 0 {
			for k, v := range attrs {
				args = append(args, k, v)
			}
		}
	}
	l.internal.Log(ctx, level, msg, args...)
}

func SetGlobalLogger(out io.Writer, options ...option) {
	globalLogger = New(out, options...)
}

func Debug(msg string, args ...any) {
	globalLogger.Debug(msg, args...)
}

func DebugWithContext(ctx context.Context, msg string, args ...any) {
	globalLogger.DebugWithContext(ctx, msg, args...)
}

func Info(msg string, args ...any) {
	globalLogger.Info(msg, args...)
}

func InfoWithContext(ctx context.Context, msg string, args ...any) {
	globalLogger.InfoWithContext(ctx, msg, args...)
}

func Warn(msg string, args ...any) {
	globalLogger.Warn(msg, args...)
}

func WarnWithContext(ctx context.Context, msg string, args ...any) {
	globalLogger.WarnWithContext(ctx, msg, args...)
}

func Error(msg string, args ...any) {
	globalLogger.Error(msg, args...)
}

func ErrorWithContext(ctx context.Context, msg string, args ...any) {
	globalLogger.ErrorWithContext(ctx, msg, args...)
}
