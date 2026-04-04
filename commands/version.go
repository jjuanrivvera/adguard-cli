package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
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

	return cmd
}

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update AdGuard Home to the latest version",
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
