package errors

import "fmt"

type ErrorCode string

const (
	InvalidInput    ErrorCode = "INVALID_INPUT"
	ConfigError     ErrorCode = "CONFIG_ERROR"
	ConnectionError ErrorCode = "CONNECTION_ERROR"
	AuthError       ErrorCode = "AUTH_ERROR"
	NotFound        ErrorCode = "NOT_FOUND"
	APIError        ErrorCode = "API_ERROR"
)

type CLIError struct {
	Code    ErrorCode
	Message string
	Cause   error
	Hint    string
}

func (e *CLIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *CLIError) Unwrap() error {
	return e.Cause
}

func New(code ErrorCode, message string) *CLIError {
	return &CLIError{Code: code, Message: message}
}

func Wrap(code ErrorCode, message string, cause error) *CLIError {
	return &CLIError{Code: code, Message: message, Cause: cause}
}

func WithHint(err *CLIError, hint string) *CLIError {
	err.Hint = hint
	return err
}

func ConnectionFailed(url string, cause error) *CLIError {
	return WithHint(
		Wrap(ConnectionError, fmt.Sprintf("cannot connect to AdGuard Home at %s", url), cause),
		"Check that AdGuard Home is running and the URL is correct. Run 'adguard-home doctor' to diagnose.",
	)
}

func AuthFailed(cause error) *CLIError {
	return WithHint(
		Wrap(AuthError, "authentication failed", cause),
		"Check your credentials with 'adguard-home setup' or set ADGUARD_USERNAME and ADGUARD_PASSWORD environment variables.",
	)
}

func ConfigNotFound() *CLIError {
	return WithHint(
		New(ConfigError, "no configuration found"),
		"Run 'adguard-home setup' to configure your AdGuard Home instance.",
	)
}

func ClientNotFound(identifier string) *CLIError {
	return WithHint(
		New(NotFound, fmt.Sprintf("client %q not found", identifier)),
		"Use 'adguard-home clients list' to see all configured clients.",
	)
}

func FormatError(err error) string {
	cliErr, ok := err.(*CLIError)
	if !ok {
		return err.Error()
	}
	msg := fmt.Sprintf("Error [%s]: %s", cliErr.Code, cliErr.Message)
	if cliErr.Cause != nil {
		msg += fmt.Sprintf("\n  Cause: %v", cliErr.Cause)
	}
	if cliErr.Hint != "" {
		msg += fmt.Sprintf("\n  Hint: %s", cliErr.Hint)
	}
	return msg
}
