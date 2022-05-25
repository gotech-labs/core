package errors

import (
	"bytes"
	"fmt"

	"golang.org/x/xerrors"
)

// TypedErrorBuilder - typed error builder
type TypedErrorBuilder interface {
	New(messages ...string) Error
	Errorf(format string, args ...interface{}) Error
	Wrap(cause error, messages ...string) Error
	Wrapf(cause error, format string, args ...interface{}) Error
	Type() Type
	Is(interface{}) bool
}

type typedErrorBuilder struct {
	errorType Type
}

func (ef *typedErrorBuilder) Type() Type {
	return ef.errorType
}

func (ef *typedErrorBuilder) New(messages ...string) Error {
	return ef.build(nil, messages...)
}

func (ef *typedErrorBuilder) Errorf(format string, args ...interface{}) Error {
	return ef.build(nil, fmt.Sprintf(format, args...))
}

func (ef *typedErrorBuilder) Wrap(cause error, messages ...string) Error {
	return ef.build(cause, messages...)
}

func (ef *typedErrorBuilder) Wrapf(cause error, format string, args ...interface{}) Error {
	return ef.build(cause, fmt.Sprintf(format, args...))
}

func (ef *typedErrorBuilder) Is(obj interface{}) bool {
	if obj == nil {
		return false
	}
	if err, ok := obj.(error); ok {
		if e := AsError(err); e != nil {
			return ef.Type() == e.Type()
		}
	}
	return false
}

func (ef *typedErrorBuilder) build(cause error, messages ...string) *typedError {
	buf := new(bytes.Buffer)
	for _, msg := range messages {
		buf.WriteString(msg)
	}
	return &typedError{
		cause:     cause,
		errorType: ef.errorType,
		message:   buf.String(),
		frame:     xerrors.Caller(2),
	}
}
