package output

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected Format
	}{
		{"json", FormatJSON},
		{"yaml", FormatYAML},
		{"table", FormatTable},
		{"", FormatTable},
		{"unknown", FormatTable},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.expected, ParseFormat(tt.input))
		})
	}
}

func TestPrintJSON(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := map[string]string{"name": "test", "value": "123"}
	err := PrintJSON(data)
	require.NoError(t, err)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	var result map[string]string
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "test", result["name"])
	assert.Equal(t, "123", result["value"])
}

func TestPrintYAML(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := map[string]string{"key": "value"}
	err := PrintYAML(data)
	require.NoError(t, err)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	assert.Contains(t, buf.String(), "key: value")
}

func TestPrintTable(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintTable(
		[]string{"Name", "Status"},
		[][]string{
			{"Server1", "Running"},
			{"Server2", "Stopped"},
		},
	)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	output := buf.String()
	assert.Contains(t, output, "NAME")
	assert.Contains(t, output, "STATUS")
	assert.Contains(t, output, "Server1")
	assert.Contains(t, output, "Running")
	assert.Contains(t, output, "Server2")
	assert.Contains(t, output, "Stopped")
}

func TestPrint_JSON(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	data := []string{"a", "b", "c"}
	err := Print(FormatJSON, data, nil, nil)
	require.NoError(t, err)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	var result []string
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestPrint_Table(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := Print(FormatTable, nil,
		[]string{"Col1", "Col2"},
		func() [][]string {
			return [][]string{{"val1", "val2"}}
		},
	)
	require.NoError(t, err)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	assert.Contains(t, buf.String(), "val1")
}
