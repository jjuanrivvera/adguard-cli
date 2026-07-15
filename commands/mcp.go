package commands

import (
	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

func newMCPCmd() *cobra.Command {
	return ophis.Command(&ophis.Config{
		ToolNamePrefix: "adguard",
		Selectors: []ophis.Selector{
			{
				// Exclude destructive commands from MCP exposure. "update" now covers the CLI
				// self-updater (`update`) + `server check-update`; "upgrade" covers `server upgrade`
				// (the server-side update, previously the top-level `update`).
				CmdSelector: ophis.ExcludeCmdsContaining("reset", "update", "upgrade"),
				// Exclude sensitive inherited flags
				InheritedFlagSelector: ophis.ExcludeFlags("instance"),
			},
		},
	})
}
