package main

import (
	"context"
	"os"

	"github.com/gotech-labs/core/log"
	"github.com/gotech-labs/core/log/formats"
	"github.com/gotech-labs/core/log/levels"
)

func main() {
	type projectKey struct{}
	ctx := context.WithValue(context.Background(), projectKey{}, "test")

	// using global logger (default text)
	log.Debug("debug message")
	log.DebugWithContext(ctx, "debug message")
	log.Info("info message")
	log.InfoWithContext(ctx, "info message")
	log.Warn("warn message")
	log.WarnWithContext(ctx, "warn message")
	log.Error("error message")
	log.ErrorWithContext(ctx, "error message")

	// using custom logger
	l := log.New(os.Stdout,
		log.WithFormat(formats.JSON),
		log.WithLevel(levels.Info),
		log.WithAttrs(func() map[string]any {
			return map[string]any{
				"env": "DEV",
			}
		}),
		log.WithContextAttrs(func(ctx context.Context) map[string]any {
			return map[string]any{
				"project": ctx.Value(projectKey{}),
			}
		}))
	l.Debug("debug message")                 // no output log message
	l.DebugWithContext(ctx, "debug message") // no output log message
	l.Info("info message", "id", 123)
	l.InfoWithContext(ctx, "info message", "id", 123)
	l.Warn("warn message", "user", &User{ID: 123, Name: "taro"})
	l.WarnWithContext(ctx, "warn message", "user", &User{ID: 123, Name: "taro"})
	l.Error("error message", "cause", "unknown error")
	l.ErrorWithContext(ctx, "error message", "cause", "unknown error")
}

type User struct {
	ID   int
	Name string
}
