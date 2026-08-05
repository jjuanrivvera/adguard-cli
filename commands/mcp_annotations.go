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
// Classification is by leaf verb and is fail-safe: only the explicit read allowlist is
// advertised read-only. Everything else — known mutations (add, enable, block, refresh)
// and any verb added later that this file has not seen — is a write. A miss costs an
// approval prompt, never a silently auto-approved mutation.
//
// The setup/config/update commands need no special handling here: newMCPCmd already keeps
// them off the MCP surface entirely. If that exclusion is ever relaxed they fall to the
// write default, which is the safe direction.

// readVerbs are leaf command names that only read state.
var readVerbs = map[string]bool{
	"list":       true,
	"find":       true,
	"status":     true,
	"stats":      true,
	"check":      true,
	"log":        true,
	"doctor":     true,
	"leases":     true,
	"interfaces": true,
	"blocked":    true,
}

// destructiveVerbs are mutations that cannot be undone from the CLI — the removed entry is
// gone and has to be re-created by hand.
var destructiveVerbs = map[string]bool{
	"delete":       true,
	"remove":       true,
	"remove-lease": true,
	"cache-clear":  true,
}

// applyMCPAnnotations stamps hints on every command in the tree. Call it once, after the
// tree is fully assembled — ophis reads cmd.Annotations at run time, in registerTools.
func applyMCPAnnotations(root *cobra.Command) {
	walkCommands(root, func(cmd *cobra.Command) {
		if !cmd.Runnable() || cmd.Hidden || cmd.Name() == "help" {
			return
		}
		switch {
		case destructiveVerbs[cmd.Name()]:
			setMCPAnnotations(cmd, false, true)
		case readVerbs[cmd.Name()]:
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
