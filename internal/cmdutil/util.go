package cmdutil

import (
	"fmt"
	"os"

	clierrors "github.com/jjuanrivvera/adguard-cli/internal/errors"
)

// Infof prints informational messages to stderr so stdout stays clean for data.
func Infof(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}

// Infoln prints an informational line to stderr.
func Infoln(msg string) {
	fmt.Fprintln(os.Stderr, msg)
}

// HandleError prints a formatted error to stderr.
func HandleError(err error) {
	fmt.Fprintln(os.Stderr, clierrors.FormatError(err))
}
