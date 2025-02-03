package formats

import (
	"io"
	"log/slog"
)

var (
	JSON Format = &format[*slog.JSONHandler]{
		name:        "JSON",
		handlerFunc: slog.NewJSONHandler,
	}
	Text Format = &format[*slog.TextHandler]{
		name:        "Text",
		handlerFunc: slog.NewTextHandler,
	}
)

type Format interface {
	NewHandler(io.Writer, *slog.HandlerOptions) slog.Handler
}

type format[T slog.Handler] struct {
	name        string
	handlerFunc func(out io.Writer, options *slog.HandlerOptions) T
}

func (f *format[T]) NewHandler(out io.Writer, options *slog.HandlerOptions) slog.Handler {
	return f.handlerFunc(out, options)
}
