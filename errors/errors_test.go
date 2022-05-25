package errors_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gotech-labs/core/errors"
)

func TestServerError(t *testing.T) {
	errorType := "test_error"
	factory := TypedError("test_error")

	// assertion
	assert.Equal(t, Type(errorType), factory.Type())

	{
		err := factory.New("test error")

		// assert error function
		assert.Equal(t, Type(errorType), err.Type())
		assert.Equal(t, "test error", err.Error())
		assert.Nil(t, err.Unwrap())
		assert.False(t, err.Is(fmt.Errorf("another error")))
		assert.False(t, err.Is(ValidationError.New("another error")))
		assert.True(t, err.Is(factory.New("same error")))

		// assert error trace
		lines := strings.Split(fmt.Sprintf("%+v", err), "\n")
		assert.Len(t, lines, 3)
		assert.Equal(t, "[test_error] test error:", lines[0])
	}

	{
		cause := fmt.Errorf("db access error")
		err := factory.Wrap(cause, "test error")
		err.AddMetaData("id", "sid_xxxxxx")

		// assert error function
		assert.Equal(t, Type(errorType), err.Type())
		assert.Equal(t, "test error", err.Error())
		assert.Equal(t, cause, err.Unwrap())
		assert.False(t, err.Is(fmt.Errorf("another error")))
		assert.False(t, err.Is(ValidationError.New("another error")))
		assert.True(t, err.Is(factory.New("same error")))
		if b, err := json.Marshal(err); assert.NoError(t, err) {
			assert.JSONEq(t, `{"error":{"cause":"db access error","id":"sid_xxxxxx","message":"test error","type":"test_error"}}`, string(b))
		}

		// assert error trace
		lines := strings.Split(fmt.Sprintf("%+v", err), "\n")
		assert.Len(t, lines, 4)
		assert.Equal(t, "[test_error] test error: metadata = map[id:sid_xxxxxx]:", lines[0])
	}

	{
		cause := fmt.Errorf("db access error")
		wrapErrorFactory := TypedError("")
		err := factory.Wrapf(cause, "test error")
		wrapError := wrapErrorFactory.Wrap(err)

		// assert error function
		assert.Equal(t, Type(errorType), wrapError.Type())
		assert.Equal(t, "test error", wrapError.Error())
		assert.Equal(t, cause, wrapError.Unwrap())
		assert.False(t, wrapError.Is(fmt.Errorf("another error")))
		assert.False(t, wrapError.Is(ValidationError.New("another error")))
		assert.True(t, wrapError.Is(factory.New("same error")))
		if b, err := json.Marshal(err); assert.NoError(t, err) {
			assert.Equal(t, `{"error":{"cause":"db access error","message":"test error","type":"test_error"}}`, string(b))
		}

		// assert error trace
		lines := strings.Split(fmt.Sprintf("%+v", wrapError), "\n")
		assert.Len(t, lines, 6)
	}

	{
		notypeErrorFactory := TypedError("")
		err := notypeErrorFactory.Errorf("")

		// assert error function
		assert.Equal(t, Type("unknown"), err.Type())
		assert.Equal(t, "unknown error", err.Error())
		if b, err := json.Marshal(err); assert.NoError(t, err) {
			assert.Equal(t, `{"error":{"message":"unknown error","type":"unexpected_error"}}`, string(b))
		}
	}
}
