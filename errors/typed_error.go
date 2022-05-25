package errors

import (
	"bytes"
	"encoding/json"
	"fmt"

	"golang.org/x/xerrors"
)

// Type - error type
type Type string

// typedError - typed error object
type typedError struct {
	cause     error
	errorType Type
	message   string
	metadata  map[string]interface{}
	frame     xerrors.Frame
}

func (e *typedError) Error() string {
	next := AsError(e.cause)
	if next != nil {
		return next.Error()
	}
	if len(e.message) > 0 {
		return e.message
	}
	if e.cause != nil {
		return e.cause.Error()
	}
	return unknownErrorMessage
}

func (e *typedError) Is(err error) bool {
	if er := AsError(err); er != nil {
		return e.Type() == er.Type()
	}
	return false
}

func (e *typedError) Unwrap() error {
	next := AsError(e.cause)
	if next != nil {
		return next.Unwrap()
	}
	return e.cause
}

func (e *typedError) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

func (e *typedError) FormatError(p xerrors.Printer) error {
	var buf bytes.Buffer
	if len(e.errorType) > 0 {
		buf.WriteString(fmt.Sprintf("[%s] ", e.errorType))
	}
	if len(e.message) > 0 {
		buf.WriteString(e.message)
	}
	if len(e.metadata) != 0 {
		buf.WriteString(fmt.Sprintf(": metadata = %+v", e.metadata))
	}
	p.Print(buf.String())
	e.frame.Format(p)
	return e.cause
}

func (e *typedError) AddMetaData(key string, value interface{}) Error {
	if e.metadata == nil {
		e.metadata = make(map[string]interface{}, 1)
	}
	e.metadata[key] = value
	return e
}

func (e *typedError) Type() Type {
	next := AsError(e.cause)
	if next != nil {
		return next.Type()
	}
	if len(e.errorType) > 0 {
		return e.errorType
	}
	return Type("unknown")
}

func (e *typedError) generateErrorMap(errMap map[string]interface{}) {
	if len(e.errorType) > 0 {
		errMap["type"] = e.errorType
	}
	if len(e.message) > 0 {
		message := e.message
		if msg, ok := errMap["message"]; ok {
			value := fmt.Sprintf("%v", msg)
			if len(value) > 0 {
				message = message + ", " + value
			}
		}
		errMap["message"] = message
	}
	if len(e.metadata) > 0 {
		for k, v := range e.metadata {
			errMap[k] = v
		}
	}
	if e.cause != nil {
		next := AsError(e.cause)
		if next != nil {
			next.generateErrorMap(errMap)
		} else {
			errMap["cause"] = e.cause.Error()
		}
	}
}

func (e *typedError) MarshalJSON() ([]byte, error) {
	errMap := make(map[string]interface{})
	e.generateErrorMap(errMap)
	if _, ok := errMap["type"]; !ok {
		errMap["type"] = UnexpectedError.Type()
	}
	if _, ok := errMap["message"]; !ok {
		errMap["message"] = unknownErrorMessage
		if cause, ok := errMap["cause"]; ok {
			errMap["message"] = cause
		}
	}
	return json.Marshal(&struct {
		Error map[string]interface{} `json:"error"`
	}{
		Error: errMap,
	})
}

const (
	unknownErrorMessage = "unknown error"
)
