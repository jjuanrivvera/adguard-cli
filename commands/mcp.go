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
				// Exclude destructive + setup/meta commands from MCP exposure. "update" covers the
				// CLI self-updater (`update`) + `server check-update`; "upgrade" covers `server
				// upgrade` (server-side update); "setup"/"config" are interactive/meta commands an
				// agent should not drive (instance selection is fixed at server startup).
				CmdSelector: ophis.ExcludeCmdsContaining("reset", "update", "upgrade", "setup", "config"),
				// Exclude sensitive inherited flags
				InheritedFlagSelector: ophis.ExcludeFlags("instance"),
			},
		},
	})
}
