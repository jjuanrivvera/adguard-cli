package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/update"
)

// newSelfUpdateCmd builds the fleet-standard `update` command: it self-updates THIS CLI
// binary (adguard-home) from its GitHub releases. Server-side updates live under
// `server upgrade` / `server check-update`. The running version is threaded in from
// NewRootCommand (adguard-cli injects it via ldflags into main, not an importable package).
func newSelfUpdateCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the adguard-home CLI to the latest release",
		Long: `Check GitHub for a newer release of this CLI and, if one exists, download it,
verify it against the release checksums, and replace the running binary in place.

To update the AdGuard Home server instead, use ` + "`adguard-home server upgrade`" + `.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
			defer cancel()
			res := update.NewUpdater(version).CheckAndUpdate(ctx)
			if res.Error != nil {
				return res.Error
			}
			out := cmd.OutOrStdout()
			if res.Updated {
				fmt.Fprintf(out, "Updated %s → %s. Restart to use the new version.\n", res.FromVersion, res.ToVersion)
			} else {
				fmt.Fprintln(out, "Already on the latest version.")
			}
			return nil
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check for a newer CLI release without installing it",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
			defer cancel()
			rel, err := update.NewUpdater(version).GetLatestRelease(ctx)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Current: %s\nLatest:  %s\n", version, rel.TagName)
			if update.IsNewer(rel.TagName, version) {
				fmt.Fprintln(out, "An update is available. Run `adguard-home update` to install it.")
			} else {
				fmt.Fprintln(out, "You are on the latest version.")
			}
			return nil
		},
	})
	return cmd
}
