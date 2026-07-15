package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

// newServerCmd groups commands that act on the AdGuard Home SERVER itself. Its subcommands
// were previously top-level `check-update` and `update`; they moved under `server` so the
// fleet-standard top-level `update` can mean "update this CLI binary" (see update.go).
func newServerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Manage the AdGuard Home server (upgrade, check for updates)",
	}
	cmd.AddCommand(newServerUpgradeCmd(), newServerCheckUpdateCmd())
	return cmd
}

// newServerCheckUpdateCmd was the top-level `check-update` command.
func newServerCheckUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check-update",
		Short: "Check if a new version of AdGuard Home is available",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}
			v, err := client.GetVersionInfo()
			if err != nil {
				return err
			}
			return output.Print(getFormat(), v,
				[]string{"Field", "Value"},
				func() [][]string {
					rows := [][]string{
						{"Current Version", v.Version},
					}
					if v.NewVersion != "" {
						rows = append(rows, []string{"New Version", v.NewVersion})
						rows = append(rows, []string{"Can Auto-Update", fmt.Sprintf("%t", v.CanAutoUpdate)})
					} else {
						rows = append(rows, []string{"Status", "Up to date"})
					}
					if v.Announcement != "" {
						rows = append(rows, []string{"Announcement", v.Announcement})
					}
					return rows
				},
			)
		},
	}
}

// newServerUpgradeCmd was the top-level `update` command (renamed to free `update` for the
// CLI self-updater). It triggers AdGuard Home's own server-side update over the REST API.
func newServerUpgradeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Update the AdGuard Home server to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}
			cmdutil.Infoln("Triggering AdGuard Home update...")
			if err := client.Update(); err != nil {
				return err
			}
			cmdutil.Infoln("Update triggered. AdGuard Home will restart.")
			return nil
		},
	}
}
