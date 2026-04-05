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
				// Exclude destructive commands from MCP exposure
				CmdSelector: ophis.ExcludeCmdsContaining("reset", "update"),
				// Exclude sensitive inherited flags
				InheritedFlagSelector: ophis.ExcludeFlags("instance"),
			},
		},
	})
}
