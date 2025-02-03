package errors

import (
	"encoding/json"
	"fmt"
)

type Error interface {
	error
	Type() Type
	Unwrap() error
	StackTrace() StackTrace
	buildStackTrace(traceLines ...string) []string
}

type StackTrace []string

type coreError struct {
	name        Type
	cause       error
	message     string
	attrs       []any
	stackFrames stackFrames
}

func (e *coreError) Type() Type {
	return e.name
}

func (e *coreError) Unwrap() error {
	next := AsError(e.cause)
	if next != nil {
		return next.Unwrap()
	}
	return e.cause
}

func (e *coreError) StackTrace() StackTrace {
	return StackTrace(e.buildStackTrace())
}

func (e *coreError) buildStackTrace(traceLines ...string) []string {
	if e.cause != nil {
		if next := AsError(e.cause); next != nil {
			traceLines = next.buildStackTrace(traceLines...)
		} else {
			traceLines = append(traceLines, fmt.Sprintf("Caused by: %v", e.cause))

		}
	}
	traceLines = append(traceLines, fmt.Sprintf("%s: %s", e.name, e.message))
	traceLines = append(traceLines, e.stackFrames.traceLines()...)
	return traceLines
}

func (e *coreError) Error() string {
	value := fmt.Sprintf("%s: %s", e.name, e.message)
	if e.cause != nil {
		value = fmt.Sprintf("%s (%s)", value, e.cause.Error())
	}
	return value
}

func (e *coreError) MarshalJSON() ([]byte, error) {
	v := struct {
		Msg        string     `json:"msg"`
		Cause      string     `json:"cause,omitempty"`
		StackTrace StackTrace `json:"stack_trace,omitempty"`
	}{
		Msg:        e.message,
		StackTrace: e.StackTrace(),
	}
	if e.cause != nil {
		v.Cause = e.cause.Error()
	}
	return json.Marshal(&v)
}

func (e *coreError) Format(s fmt.State, verb rune) {
	if e == nil {
		fmt.Fprintf(s, "<nil>")
		return
	}
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", *e)
			return
		}
		fallthrough
	case 's':
		fmt.Fprintf(s, "%s", e.message)
	case 'q':
		fmt.Fprintf(s, "%q", e.message)
	}
}

func newError(name Type, cause error, message string, args ...any) *coreError {
	return &coreError{
		name:        name,
		cause:       cause,
		message:     message,
		attrs:       args,
		stackFrames: callers(4).StackFrames(),
	}
}
