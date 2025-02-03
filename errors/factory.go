package errors

type ErrorFactory interface {
	New(format string, args ...any) Error
	Wrap(cause error, format string, args ...any) Error
	Is(err error) bool
}

type errorFactory struct {
	name Type
}

func (ef *errorFactory) New(format string, args ...any) Error {
	return newError(ef.name, nil, format, args...)
}

func (ef *errorFactory) Wrap(cause error, format string, args ...any) Error {
	return newError(ef.name, cause, format, args...)
}

func (ef *errorFactory) Is(err error) bool {
	if e := AsError(err); e != nil {
		return ef.name == e.Type()
	}
	return false
}
