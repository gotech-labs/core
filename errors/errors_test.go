package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	//. "github.com/gotech-labs/core/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsError(t *testing.T) {
	err := ValidationError.New("illegal code")
	assert.Nil(t, AsError(nil))
	assert.Nil(t, AsError(errors.New("golang error")))
	assert.Equal(t, err, AsError(err))
}

func TestUnwrap(t *testing.T) {
	for _, in := range []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "test - nil",
			err:      nil,
			expected: nil,
		},
		{
			name:     "test - golang error",
			err:      errors.New("mysql open error"),
			expected: nil,
		},
		{
			name:     "test - new error",
			err:      ValidationError.New("user name is a required field"),
			expected: nil,
		},
		{
			name:     "test - wrap error",
			err:      InitializationError.Wrap(errors.New("mysql open error"), "database access error"),
			expected: errors.New("mysql open error"),
		},
		{
			name: "test - nested error",
			err: UnexpectedError.Wrap(
				InitializationError.Wrap(errors.New("mysql open error"), "database access error"), "nested error message"),
			expected: errors.New("mysql open error"),
		},
	} {
		assert.Equal(t, in.expected, Unwrap(in.err))
	}
}

func TestErrorType(t *testing.T) {
	assert.False(t, ValidationError.Is(nil))
	assert.False(t, ValidationError.Is(errors.New("mysql open error")))
	assert.True(t, ValidationError.Is(ValidationError.New("user name is a required field")))
}

/*
func TestDefinedErrorType(t *testing.T) {
	testData := testData()
	assert.True(t, ValidationError.Is(testData.err))
	assert.True(t, InitializationError.Is(testData.wrapErr))
	assert.True(t, UnexpectedError.Is(testData.nestErr))
	assert.False(t, UnexpectedError.Is(testData.cause))
}

func TestErrorUnwrap(t *testing.T) {
	testData := testData()
	assert.Nil(t, testData.err.Unwrap())
	assert.Equal(t, testData.cause, testData.wrapErr.Unwrap())
	assert.Equal(t, testData.cause, testData.nestErr.Unwrap())
}

func TestErrorMessage(t *testing.T) {
	testData := testData()
	assert.Equal(t, "invalid parameter: user name is a required field", testData.err.Error())
	assert.Equal(t, "initialization error: database access error (mysql open error)", testData.wrapErr.Error())
	assert.Equal(t, "unexpected error: nested database error (initialization error: database access error (mysql open error))", testData.nestErr.Error())
}
*/

func TestErrorFormat(t *testing.T) {
	var err *coreError
	assert.Equal(t, "<nil>", fmt.Sprint(err))

	err = &coreError{
		name:        Type("validation error"),
		cause:       nil,
		message:     "user name is a required field",
		attrs:       nil,
		stackFrames: nil,
	}
	assert.Equal(t, "user name is a required field", fmt.Sprint(err))
	assert.Equal(t, "user name is a required field", fmt.Sprintf("%+s", err))
	assert.Equal(t, "\"user name is a required field\"", fmt.Sprintf("%+q", err))
	assert.Equal(t, "{name:validation error cause:<nil> message:user name is a required field attrs:[] stackFrames:[]}", fmt.Sprintf("%+v", err))
}

func TestErrorTrace(t *testing.T) {
	for _, in := range []struct {
		name     string
		err      Error
		expected []string
	}{
		{
			name: "test - new error json",
			err:  ValidationError.New("user name is a required field"),
			expected: []string{
				"invalid parameter: user name is a required field",
				"    at ${CURRENT_DIR}/errors_test.go:xxx (TestErrorTrace)",
			},
		},
		{
			name: "test - wrap error json",
			err:  InitializationError.Wrap(errors.New("mysql open error"), "database access error"),
			expected: []string{
				"Caused by: mysql open error",
				"initialization error: database access error",
				"    at ${CURRENT_DIR}/errors_test.go:xxx (TestErrorTrace)",
			},
		},
		{
			name: "test - nested error json",
			err: UnexpectedError.Wrap(
				InitializationError.Wrap(errors.New("mysql open error"), "database access error"), "nested error message"),
			expected: []string{
				"Caused by: mysql open error",
				"initialization error: database access error",
				"    at ${CURRENT_DIR}/errors_test.go:xxx (TestErrorTrace)",
				"unexpected error: nested error message",
				"    at ${CURRENT_DIR}/errors_test.go:xxx (TestErrorTrace)",
			},
		},
	} {
		stackTrace := in.err.StackTrace()
		require.Len(t, stackTrace, len(in.expected))
		for i, trace := range in.expected {
			actual := cutLineNumber(stackTrace[i])
			expected := strings.ReplaceAll(trace, "${CURRENT_DIR}", currentDir)
			assert.Equal(t, expected, actual)
		}
	}
}

func TestErrorJSON(t *testing.T) {
	for _, in := range []struct {
		name     string
		err      Error
		expected string
	}{
		{
			name: "test - new error json",
			err:  ValidationError.New("user name is a required field"),
			expected: `{
				"msg": "user name is a required field",
				"stack_trace": [
					"invalid parameter: user name is a required field",
					"    at ${CURRENT_DIR}/errors_test.go:xxx (TestErrorJSON)"
				]
			}`,
		},
		{
			name: "test - wrap error json",
			err:  InitializationError.Wrap(errors.New("mysql open error"), "database access error"),
			expected: `{
				"cause": "mysql open error",
				"msg": "database access error",
				"stack_trace": [
					"Caused by: mysql open error",
					"initialization error: database access error",
					"    at ${CURRENT_DIR}/errors_test.go:xxx (TestErrorJSON)"
				]
			}`,
		},
		{
			name: "test - nested error json",
			err: UnexpectedError.Wrap(
				InitializationError.Wrap(errors.New("mysql open error"), "database access error"), "nested error message"),
			expected: `{
				"cause": "initialization error: database access error (mysql open error)",
				"msg": "nested error message",
				"stack_trace": [
					"Caused by: mysql open error",
					"initialization error: database access error",
					"    at ${CURRENT_DIR}/errors_test.go:xxx (TestErrorJSON)",
					"unexpected error: nested error message",
					"    at ${CURRENT_DIR}/errors_test.go:xxx (TestErrorJSON)"
				]
			}`,
		},
	} {
		data, e := json.MarshalIndent(in.err, "", "  ")
		require.NoError(t, e)
		actual := cutLineNumber(string(data))
		expected := strings.ReplaceAll(in.expected, "${CURRENT_DIR}", currentDir)
		assert.JSONEq(t, expected, actual, in.name)
	}
}

func cutLineNumber(value string) string {
	values := strings.Split(value, ".go")
	for i, v := range values {
		pos := strings.Index(v, " (")
		if !strings.HasPrefix(v, ":") || pos == -1 {
			continue
		}
		values[i] = ":xxx" + v[pos:]
	}
	return strings.Join(values, ".go")
}

var currentDir, _ = os.Getwd()
