package log

import (
	"context"
	"io"

	"github.com/gotech-labs/core/log/formats"
	"github.com/gotech-labs/core/log/levels"
)

type settings struct {
	out                  io.Writer
	format               formats.Format
	level                levels.Level
	withAttrsFunc        func() map[string]any
	withContextAttrsFunc func(context.Context) map[string]any
	pretty               bool
}

type option func(*settings)

func WithFormat(format formats.Format) func(options *settings) {
	return func(settings *settings) {
		settings.format = format
	}
}

func WithLevel(level levels.Level) func(options *settings) {
	return func(settings *settings) {
		settings.level = level
	}
}

func WithAttrs(withAttrsFunc func() map[string]any) func(options *settings) {
	return func(settings *settings) {
		settings.withAttrsFunc = withAttrsFunc
	}
}

func WithContextAttrs(withContextAttrsFunc func(context.Context) map[string]any) func(options *settings) {
	return func(settings *settings) {
		settings.withContextAttrsFunc = withContextAttrsFunc
	}
}
