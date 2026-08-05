package commands

import (
	"strings"
	"testing"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

func findCmd(root *cobra.Command, path string) *cobra.Command {
	cur := root
	for _, part := range strings.Fields(path) {
		var next *cobra.Command
		for _, sub := range cur.Commands() {
			if sub.Name() == part {
				next = sub
				break
			}
		}
		if next == nil {
			return nil
		}
		cur = next
	}
	return cur
}

func assertHints(t *testing.T, root *cobra.Command, path, readOnly, destructive string) {
	t.Helper()
	cmd := findCmd(root, path)
	if cmd == nil {
		t.Fatalf("command %q not found", path)
	}
	if got := cmd.Annotations[ophis.AnnotationReadOnly]; got != readOnly {
		t.Errorf("%q: readOnlyHint = %q, want %q", path, got, readOnly)
	}
	if got := cmd.Annotations[ophis.AnnotationDestructive]; got != destructive {
		t.Errorf("%q: destructiveHint = %q, want %q", path, got, destructive)
	}
}

// TestApplyMCPAnnotations_Buckets pins the classification of the real tree. A host in
// read-only mode keeps only the readOnlyHint=true tools, so a regression here silently
// changes what an agent may run unattended.
func TestApplyMCPAnnotations_Buckets(t *testing.T) {
	root := NewRootCommand("test", "none", "none")

	reads := []string{
		"status", "stats", "log", "doctor",
		"clients list", "clients find",
		"rewrites list", "filters list", "services list", "services blocked",
		"dhcp status", "dhcp leases", "dhcp interfaces",
		"dns check", "tls status", "parental status", "safebrowsing status",
		"safesearch status", "access list",
	}
	for _, p := range reads {
		assertHints(t, root, p, "true", "")
	}

	writes := []string{
		"clients add", "rewrites add", "filters add", "filters refresh",
		"dhcp add-lease", "services block", "services unblock",
		"parental enable", "parental disable",
		"safebrowsing enable", "safebrowsing disable",
		"status enable", "status disable",
	}
	for _, p := range writes {
		assertHints(t, root, p, "false", "")
	}

	destructive := []string{
		"clients delete", "rewrites delete", "filters remove",
		"dhcp remove-lease", "dns cache-clear",
	}
	for _, p := range destructive {
		assertHints(t, root, p, "false", "true")
	}
}

// TestApplyMCPAnnotations_EveryLeafAnnotated is the regression guard that matters: a new
// command must not land unannotated, because an unannotated tool is invisible to
// read-only hosts and its intent was never declared.
func TestApplyMCPAnnotations_EveryLeafAnnotated(t *testing.T) {
	root := NewRootCommand("test", "none", "none")

	var missing []string
	walkCommands(root, func(cmd *cobra.Command) {
		if !cmd.Runnable() || cmd.Hidden || cmd.Name() == "help" {
			return
		}
		if _, ok := cmd.Annotations[ophis.AnnotationReadOnly]; !ok {
			missing = append(missing, cmd.CommandPath())
		}
	})
	if len(missing) > 0 {
		t.Errorf("commands missing MCP annotations: %v", missing)
	}
}

// TestApplyMCPAnnotations_UnknownVerbIsWrite pins the fail-safe default.
func TestApplyMCPAnnotations_UnknownVerbIsWrite(t *testing.T) {
	root := &cobra.Command{Use: "adguard-home"}
	grp := &cobra.Command{Use: "widgets"}
	grp.AddCommand(&cobra.Command{
		Use:  "frobnicate",
		RunE: func(_ *cobra.Command, _ []string) error { return nil },
	})
	root.AddCommand(grp)

	applyMCPAnnotations(root)

	assertHints(t, root, "widgets frobnicate", "false", "")
}

func TestSetMCPAnnotations_PreservesUnrelatedKeys(t *testing.T) {
	cmd := &cobra.Command{Use: "thing", Annotations: map[string]string{"custom": "kept"}}

	setMCPAnnotations(cmd, true, false)

	if cmd.Annotations["custom"] != "kept" {
		t.Errorf("unrelated annotation clobbered: %v", cmd.Annotations)
	}
	if cmd.Annotations[ophis.AnnotationOpenWorld] != "true" {
		t.Error("openWorldHint should always be true")
	}
	if _, ok := cmd.Annotations[ophis.AnnotationDestructive]; ok {
		t.Error("destructiveHint should be omitted when false")
	}
}
