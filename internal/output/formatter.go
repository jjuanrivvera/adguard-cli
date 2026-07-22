package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"gopkg.in/yaml.v3"
)

type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

func ParseFormat(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "yaml":
		return FormatYAML
	default:
		return FormatTable
	}
}

// PrintJSON outputs data as indented JSON to stdout.
func PrintJSON(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// PrintYAML outputs data as YAML to stdout.
func PrintYAML(data any) error {
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	defer enc.Close()
	return enc.Encode(data)
}

// PrintTable outputs data as a formatted table to stdout.
func PrintTable(headers []string, rows [][]string) {
	// tablewriter v1 replaces the v0 setters with functional options; this reproduces the old
	// look: no borders or column separators, a header line, left-aligned unwrapped cells.
	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithRendition(tw.Rendition{
			Borders: tw.BorderNone,
			Settings: tw.Settings{
				Separators: tw.Separators{BetweenRows: tw.Off, BetweenColumns: tw.Off},
				Lines:      tw.Lines{ShowTop: tw.Off, ShowBottom: tw.Off, ShowHeaderLine: tw.On},
			},
		}),
		tablewriter.WithHeaderAutoWrap(tw.WrapNone),
		tablewriter.WithRowAutoWrap(tw.WrapNone),
		tablewriter.WithHeaderAlignment(tw.AlignLeft),
		tablewriter.WithRowAlignment(tw.AlignLeft),
		tablewriter.WithPadding(tw.Padding{Right: "  "}),
	)
	table.Header(headers)
	_ = table.Bulk(rows)
	_ = table.Render()
}

// Print dispatches to the correct formatter based on format string.
func Print(format Format, data any, headers []string, toRows func() [][]string) error {
	switch format {
	case FormatJSON:
		return PrintJSON(data)
	case FormatYAML:
		return PrintYAML(data)
	case FormatTable:
		PrintTable(headers, toRows())
		return nil
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
}
