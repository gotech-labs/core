package errors

import (
	"errors"
)

var (
	ValidationError             = Define("invalid parameter")
	IllegalArgumentError        = Define("illegal argument")
	InitializationError         = Define("initialization error")
	AuthenticationRequiredError = Define("authentication required")
	UnexpectedError             = Define("unexpected error")
)

func Define(name string) ErrorFactory {
	return &errorFactory{
		name: Type(name),
	}
}

func Unwrap(err error) error {
	target := AsError(err)
	if target == nil {
		return errors.Unwrap(err)
	}
	return target.Unwrap()
}

func AsError(err error) Error {
	if err == nil {
		return nil
	}
	var e *coreError
	if errors.As(err, &e) {
		return e
	}
	return nil
}

type Type string
