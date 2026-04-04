package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newSafeBrowsingCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "safebrowsing",
		Short: "Manage safe browsing protection",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "status",
			Short: "Show safe browsing status",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				s, err := client.GetSafeBrowsingStatus()
				if err != nil {
					return err
				}
				return output.Print(getFormat(), s,
					[]string{"Feature", "Status"},
					func() [][]string {
						status := "disabled"
						if s.Enabled {
							status = "enabled"
						}
						return [][]string{{"Safe Browsing", status}}
					},
				)
			},
		},
		&cobra.Command{
			Use:   "enable",
			Short: "Enable safe browsing",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.SetSafeBrowsing(true); err != nil {
					return err
				}
				cmdutil.Infoln("Safe browsing enabled.")
				return nil
			},
		},
		&cobra.Command{
			Use:   "disable",
			Short: "Disable safe browsing",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.SetSafeBrowsing(false); err != nil {
					return err
				}
				cmdutil.Infoln("Safe browsing disabled.")
				return nil
			},
		},
	)

	return cmd
}
