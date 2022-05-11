package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"

	"github.com/gotech-labs/core/system"
)

var (
	defaultLogger = New(os.Stdout)
)

func New(out io.Writer) *Logger {
	return &Logger{
		internal: zerolog.New(out),
	}
}

func Debug() *zerolog.Event {
	return defaultLogger.Debug()
}

func Info() *zerolog.Event {
	return defaultLogger.Info()
}

func Warn() *zerolog.Event {
	return defaultLogger.Warn()
}

func Error() *zerolog.Event {
	return defaultLogger.Error()
}

func Panic() *zerolog.Event {
	return defaultLogger.Panic()
}

func SetupDefaultLogger(out io.Writer) {
	defaultLogger = New(out)
}

type Logger struct {
	internal zerolog.Logger
}

func (l *Logger) Debug() *zerolog.Event {
	return withParams(l.internal.Debug())
}

func (l *Logger) Info() *zerolog.Event {
	return withParams(l.internal.Info())
}

func (l *Logger) Warn() *zerolog.Event {
	return withParams(l.internal.Warn())
}

func (l *Logger) Error() *zerolog.Event {
	return withParams(l.internal.Error())
}

func (l *Logger) Panic() *zerolog.Event {
	return withParams(l.internal.Panic())
}

func withParams(evt *zerolog.Event) *zerolog.Event {
	return evt.Time("time", system.CurrentTime())
}
