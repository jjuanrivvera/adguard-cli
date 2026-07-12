package commands

import (
	"strings"
	"testing"
)

func TestScanSecretLine(t *testing.T) {
	// A long paste (well past canonical mode's 1024-char MAX_CANON) reads intact — the whole point.
	long := strings.Repeat("A", 2000)
	if got, err := scanSecretLine(strings.NewReader(long + "\r")); err != nil || got != long {
		t.Errorf("long line: len=%d err=%v", len(got), err)
	}
	if got, _ := scanSecretLine(strings.NewReader("ab\x7fc\n")); got != "ac" { // backspace edits
		t.Errorf("backspace: %q", got)
	}
	if _, err := scanSecretLine(strings.NewReader("\x03")); err == nil { // Ctrl-C cancels
		t.Error("Ctrl-C should cancel")
	}
	if got, _ := scanSecretLine(strings.NewReader("xyz")); got != "xyz" { // EOF returns buffered
		t.Errorf("EOF: %q", got)
	}
}
