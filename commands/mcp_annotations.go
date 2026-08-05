package commands

import (
	"strconv"

	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

// MCP tool annotations (cliwright GOAL.md §3b).
//
// Hosts running a read-only MCP session allow a tool only when readOnlyHint is strictly
// true, and treat a missing annotation as a write rather than as unknown. With nothing
// annotated the host finds no allowed tool and drops the entire server, so every command
// has to declare itself.
//
// Classification is by leaf verb and is fail-safe: only the verbs isReadVerb recognises are
// advertised read-only. Everything else — known mutations (add, enable, block, refresh) and
// any verb added later that this file has not seen — is a write. A miss costs an approval
// prompt, never a silently auto-approved mutation.
//
// The setup/config/update commands need no special handling here: newMCPCmd already keeps
// them off the MCP surface entirely. If that exclusion is ever relaxed they fall to the
// write default, which is the safe direction.

// isReadVerb reports whether a leaf command name only reads state.
func isReadVerb(name string) bool {
	switch name {
	case "list", "find", "status", "stats", "check", "log", "doctor",
		"leases", "interfaces", "blocked":
		return true
	}
	return false
}

// isDestructiveVerb reports whether a leaf command name is a mutation that cannot be undone
// from the CLI — the removed entry is gone and has to be re-created by hand.
func isDestructiveVerb(name string) bool {
	switch name {
	case "delete", "remove", "remove-lease", "cache-clear":
		return true
	}
	return false
}

// applyMCPAnnotations stamps hints on every command in the tree. Call it once, after the
// tree is fully assembled — ophis reads cmd.Annotations at run time, in registerTools.
func applyMCPAnnotations(root *cobra.Command) {
	walkCommands(root, func(cmd *cobra.Command) {
		if !cmd.Runnable() || cmd.Hidden || cmd.Name() == "help" {
			return
		}
		switch {
		case isDestructiveVerb(cmd.Name()):
			setMCPAnnotations(cmd, false, true)
		case isReadVerb(cmd.Name()):
			setMCPAnnotations(cmd, true, false)
		default:
			setMCPAnnotations(cmd, false, false)
		}
	})
}

// walkCommands visits every command below root, leaves first.
func walkCommands(root *cobra.Command, visit func(*cobra.Command)) {
	var walk func(*cobra.Command)
	walk = func(c *cobra.Command) {
		for _, sub := range c.Commands() {
			walk(sub)
			visit(sub)
		}
	}
	walk(root)
}

// setMCPAnnotations merges MCP hints into cmd.Annotations, preserving unrelated keys.
// Every annotated command talks to an AdGuard Home instance, so openWorldHint is always
// true.
//
// destructiveHint is written only when true; leaving it off lets the host apply the spec
// default (true when readOnlyHint is false), which fails safe.
func setMCPAnnotations(cmd *cobra.Command, readOnly, destructive bool) {
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[ophis.AnnotationReadOnly] = strconv.FormatBool(readOnly)
	cmd.Annotations[ophis.AnnotationOpenWorld] = "true"
	if destructive {
		cmd.Annotations[ophis.AnnotationDestructive] = "true"
	}
}
