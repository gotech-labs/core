package log_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	. "github.com/gotech-labs/core/log"
	"github.com/gotech-labs/core/log/formats"
	"github.com/gotech-labs/core/log/levels"
	"github.com/gotech-labs/core/testing/runner"
	"github.com/stretchr/testify/assert"
)

func TestGlobalLogger(t *testing.T) {
	runner.RunTest(t, "logging test", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		SetGlobalLogger(buf, WithLevel(levels.Debug), WithFormat(formats.JSON))

		DebugWithContext(context.Background(), "debug message")
		InfoWithContext(context.Background(), "info message")
		WarnWithContext(context.Background(), "warn message")
		ErrorWithContext(context.Background(), "error message")
		expected := strings.Join([]string{
			`{"time":"` + runner.TestingTimeStr + `","level":"DEBUG","msg":"debug message"}`,
			`{"time":"` + runner.TestingTimeStr + `","level":"INFO","msg":"info message"}`,
			`{"time":"` + runner.TestingTimeStr + `","level":"WARN","msg":"warn message"}`,
			`{"time":"` + runner.TestingTimeStr + `","level":"ERROR","msg":"error message"}`,
		}, "\n") + "\n"
		assert.Equal(t, expected, buf.String())
	})
}

func TestLogger(t *testing.T) {
	testingTimeStr := runner.TestingTime.Format(time.RFC3339Nano)
	// output log test
	runner.RunTest(t, "logging test", func(t *testing.T) {
		callLoggingFunction := func(level levels.Level) string {
			buf := bytes.NewBuffer(nil)
			l := New(buf, WithLevel(level), WithFormat(formats.Text))
			l.DebugWithContext(context.Background(), "debug message")
			l.InfoWithContext(context.Background(), "info message")
			l.WarnWithContext(context.Background(), "warn message")
			l.ErrorWithContext(context.Background(), "error message")
			return buf.String()
		}
		{
			actual := callLoggingFunction(levels.Debug)
			expected := strings.Join([]string{
				`time=` + testingTimeStr + ` level=DEBUG msg="debug message"`,
				`time=` + testingTimeStr + ` level=INFO msg="info message"`,
				`time=` + testingTimeStr + ` level=WARN msg="warn message"`,
				`time=` + testingTimeStr + ` level=ERROR msg="error message"`,
			}, "\n") + "\n"
			assert.Equal(t, expected, actual)
		}
		{
			actual := callLoggingFunction(levels.Info)
			expected := strings.Join([]string{
				`time=` + runner.TestingTimeStr + ` level=INFO msg="info message"`,
				`time=` + runner.TestingTimeStr + ` level=WARN msg="warn message"`,
				`time=` + runner.TestingTimeStr + ` level=ERROR msg="error message"`,
			}, "\n") + "\n"
			assert.Equal(t, expected, actual)
		}
		{
			actual := callLoggingFunction(levels.Warn)
			expected := strings.Join([]string{
				`time=` + runner.TestingTimeStr + ` level=WARN msg="warn message"`,
				`time=` + runner.TestingTimeStr + ` level=ERROR msg="error message"`,
			}, "\n") + "\n"
			assert.Equal(t, expected, actual)
		}
		{
			actual := callLoggingFunction(levels.Error)
			expected := strings.Join([]string{
				`time=` + runner.TestingTimeStr + ` level=ERROR msg="error message"`,
			}, "\n") + "\n"
			assert.Equal(t, expected, actual)
		}
	})
}
