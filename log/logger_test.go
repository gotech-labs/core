package log_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	. "github.com/gotech-labs/core/log"
	"github.com/gotech-labs/core/system"
)

func TestLogger(t *testing.T) {
	type input struct {
		level  string
		logEvt func(*Logger) *zerolog.Event
	}
	for _, input := range []input{
		{
			level:  "debug",
			logEvt: func(l *Logger) *zerolog.Event { return l.Debug() },
		},
		{
			level:  "info",
			logEvt: func(l *Logger) *zerolog.Event { return l.Info() },
		},
		{
			level:  "warn",
			logEvt: func(l *Logger) *zerolog.Event { return l.Warn() },
		},
		{
			level:  "error",
			logEvt: func(l *Logger) *zerolog.Event { return l.Error() },
		},
	} {
		var (
			level  = input.level
			logEvt = input.logEvt
		)
		system.RunTest(t, fmt.Sprintf("%s message", level), func(t *testing.T) {
			var (
				buf      = bytes.NewBuffer(nil)
				expected = fmt.Sprintf(`{
					"level": "%s",
					"time": "2022-12-24T00:00:00+09:00",
					"message": "%s message"
				}`, level, level)
			)
			logEvt(New(buf)).Msgf("%s message", level)

			// assert log message
			assert.JSONEq(t, expected, buf.String())
		})
	}

	// panic error test (with default logger)
	system.RunTest(t, "panic message", func(t *testing.T) {
		var (
			buf      = bytes.NewBuffer(nil)
			expected = `{
				"level": "panic",
				"time": "2022-12-24T00:00:00+09:00",
				"message": "panic message"
			}`
		)
		SetupDefaultLogger(buf)

		assert.Panics(t, func() {
			Panic().Msg("panic message")
		})

		// assert log message
		assert.JSONEq(t, expected, buf.String())
	})
}
