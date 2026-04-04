package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show AdGuard Home server status",
		RunE:  runStatus,
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "enable",
			Short: "Enable DNS protection",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.SetProtection(true); err != nil {
					return err
				}
				cmdutil.Infoln("DNS protection enabled.")
				return nil
			},
		},
		&cobra.Command{
			Use:   "disable",
			Short: "Disable DNS protection",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.SetProtection(false); err != nil {
					return err
				}
				cmdutil.Infoln("DNS protection disabled.")
				return nil
			},
		},
	)

	return cmd
}

func runStatus(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	status, err := client.GetStatus()
	if err != nil {
		return err
	}

	format := getFormat()

	protection := "Disabled"
	if status.ProtectionEnabled {
		protection = "Enabled"
	}

	running := "Stopped"
	if status.Running {
		running = "Running"
	}

	return output.Print(format, status,
		[]string{"Field", "Value"},
		func() [][]string {
			return [][]string{
				{"Version", status.Version},
				{"Running", running},
				{"Protection", protection},
				{"DNS Port", fmt.Sprintf("%d", status.DNSPort)},
				{"HTTP Port", fmt.Sprintf("%d", status.HTTPPort)},
			}
		},
	)
}
