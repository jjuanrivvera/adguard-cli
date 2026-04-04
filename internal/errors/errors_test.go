package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLIError_Error(t *testing.T) {
	err := New(ConnectionError, "cannot connect")
	assert.Equal(t, "cannot connect", err.Error())
}

func TestCLIError_ErrorWithCause(t *testing.T) {
	cause := fmt.Errorf("dial tcp timeout")
	err := Wrap(ConnectionError, "cannot connect", cause)
	assert.Equal(t, "cannot connect: dial tcp timeout", err.Error())
}

func TestCLIError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("underlying error")
	err := Wrap(APIError, "api call failed", cause)
	assert.Equal(t, cause, err.Unwrap())
}

func TestWithHint(t *testing.T) {
	err := New(ConfigError, "no config found")
	err = WithHint(err, "Run 'adguard-home setup' first")
	assert.Equal(t, "Run 'adguard-home setup' first", err.Hint)
}

func TestConnectionFailed(t *testing.T) {
	cause := fmt.Errorf("connection refused")
	err := ConnectionFailed("http://192.168.0.1:8001", cause)
	assert.Equal(t, ConnectionError, err.Code)
	assert.Contains(t, err.Error(), "192.168.0.1:8001")
	assert.Contains(t, err.Hint, "adguard-home doctor")
}

func TestAuthFailed(t *testing.T) {
	err := AuthFailed(fmt.Errorf("401"))
	assert.Equal(t, AuthError, err.Code)
	assert.Contains(t, err.Hint, "adguard-home setup")
}

func TestConfigNotFound(t *testing.T) {
	err := ConfigNotFound()
	assert.Equal(t, ConfigError, err.Code)
	assert.Contains(t, err.Hint, "adguard-home setup")
}

func TestClientNotFound(t *testing.T) {
	err := ClientNotFound("192.168.0.99")
	assert.Equal(t, NotFound, err.Code)
	assert.Contains(t, err.Error(), "192.168.0.99")
	assert.Contains(t, err.Hint, "clients list")
}

func TestFormatError_CLIError(t *testing.T) {
	err := WithHint(
		Wrap(ConnectionError, "cannot connect", fmt.Errorf("timeout")),
		"Check your network",
	)
	output := FormatError(err)
	assert.Contains(t, output, "Error [CONNECTION_ERROR]")
	assert.Contains(t, output, "cannot connect")
	assert.Contains(t, output, "Cause: timeout")
	assert.Contains(t, output, "Hint: Check your network")
}

func TestFormatError_StandardError(t *testing.T) {
	err := fmt.Errorf("regular error")
	output := FormatError(err)
	assert.Equal(t, "regular error", output)
}

func TestFormatError_NoCause(t *testing.T) {
	err := New(InvalidInput, "bad input")
	output := FormatError(err)
	assert.Contains(t, output, "Error [INVALID_INPUT]: bad input")
	assert.NotContains(t, output, "Cause:")
}

func TestFormatError_NoHint(t *testing.T) {
	err := Wrap(APIError, "api failed", fmt.Errorf("500"))
	output := FormatError(err)
	assert.NotContains(t, output, "Hint:")
}
