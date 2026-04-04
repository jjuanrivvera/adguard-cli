package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/api"
	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newDHCPCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dhcp",
		Short: "Manage DHCP server",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "status",
			Short: "Show DHCP server status",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				s, err := client.GetDHCPStatus()
				if err != nil {
					return err
				}
				return output.Print(getFormat(), s,
					[]string{"Field", "Value"},
					func() [][]string {
						enabled := "disabled"
						if s.Enabled {
							enabled = "enabled"
						}
						return [][]string{
							{"Enabled", enabled},
							{"Interface", s.InterfaceName},
							{"IPv4 Gateway", s.V4.GatewayIP},
							{"IPv4 Range", fmt.Sprintf("%s - %s", s.V4.RangeStart, s.V4.RangeEnd)},
							{"IPv4 Subnet", s.V4.SubnetMask},
							{"Active Leases", fmt.Sprintf("%d", len(s.Leases))},
							{"Static Leases", fmt.Sprintf("%d", len(s.StaticLeases))},
						}
					},
				)
			},
		},
		&cobra.Command{
			Use:   "leases",
			Short: "List active and static DHCP leases",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				s, err := client.GetDHCPStatus()
				if err != nil {
					return err
				}
				all := append(s.StaticLeases, s.Leases...)
				return output.Print(getFormat(), all,
					[]string{"MAC", "IP", "Hostname"},
					func() [][]string {
						var rows [][]string
						for _, l := range all {
							rows = append(rows, []string{l.MAC, l.IP, l.Hostname})
						}
						return rows
					},
				)
			},
		},
		&cobra.Command{
			Use:   "add-lease [mac] [ip] [hostname]",
			Short: "Add a static DHCP lease",
			Args:  cobra.ExactArgs(3),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.AddStaticLease(api.DHCPStaticLease{MAC: args[0], IP: args[1], Hostname: args[2]}); err != nil {
					return err
				}
				cmdutil.Infof("Static lease added: %s -> %s (%s)\n", args[0], args[1], args[2])
				return nil
			},
		},
		&cobra.Command{
			Use:   "remove-lease [mac] [ip] [hostname]",
			Short: "Remove a static DHCP lease",
			Args:  cobra.ExactArgs(3),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.RemoveStaticLease(api.DHCPStaticLease{MAC: args[0], IP: args[1], Hostname: args[2]}); err != nil {
					return err
				}
				cmdutil.Infoln("Static lease removed.")
				return nil
			},
		},
		&cobra.Command{
			Use:   "interfaces",
			Short: "List available network interfaces for DHCP",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				ifaces, err := client.GetDHCPInterfaces()
				if err != nil {
					return err
				}
				return output.PrintJSON(ifaces)
			},
		},
		&cobra.Command{
			Use:   "reset",
			Short: "Reset DHCP configuration",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.ResetDHCP(); err != nil {
					return err
				}
				cmdutil.Infoln("DHCP configuration reset.")
				return nil
			},
		},
	)

	return cmd
}
