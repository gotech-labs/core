package errors

import (
	"golang.org/x/xerrors"
)

var (
	ValidationError             = TypedError("invalid_parameter")
	AuthenticationRequiredError = TypedError("authentication_required")
	InitializationError         = TypedError("initialization_error")
	UnexpectedError             = TypedError("unexpected_error")
)

// TypedError - create a new typed error builder
func TypedError(errorType string) TypedErrorBuilder {
	return &typedErrorBuilder{
		errorType: Type(errorType),
	}
}

// Error - error interface
type Error interface {
	error
	Is(error) bool
	Unwrap() error
	AddMetaData(field string, data interface{}) Error
	Type() Type

	generateErrorMap(errMap map[string]interface{})
}

func AsError(err error) Error {
	if err == nil {
		return nil
	}
	var e Error
	if xerrors.As(err, &e) {
		return e
	}
	return nil
}
