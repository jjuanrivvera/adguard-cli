package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newDNSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dns",
		Short: "Manage DNS server configuration",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "config",
			Short: "Show DNS server configuration",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				d, err := client.GetDNSConfig()
				if err != nil {
					return err
				}
				return output.Print(getFormat(), d,
					[]string{"Setting", "Value"},
					func() [][]string {
						return [][]string{
							{"Upstream DNS", strings.Join(d.UpstreamDNS, ", ")},
							{"Bootstrap DNS", strings.Join(d.BootstrapDNS, ", ")},
							{"Fallback DNS", strings.Join(d.FallbackDNS, ", ")},
							{"Upstream Mode", d.UpstreamMode},
							{"Blocking Mode", d.BlockingMode},
							{"Rate Limit", fmt.Sprintf("%d", d.RateLimit)},
							{"Cache Size", fmt.Sprintf("%d bytes", d.CacheSize)},
							{"Cache TTL", fmt.Sprintf("%d-%d sec", d.CacheMinTTL, d.CacheMaxTTL)},
							{"Cache Optimistic", fmt.Sprintf("%t", d.CacheOptimistic)},
							{"DNSSEC", fmt.Sprintf("%t", d.DNSSECEnabled)},
							{"EDNS CS", fmt.Sprintf("%t", d.EDNSCSEnabled)},
							{"Disable IPv6", fmt.Sprintf("%t", d.DisableIPv6)},
						}
					},
				)
			},
		},
		&cobra.Command{
			Use:   "cache-clear",
			Short: "Clear the DNS cache",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.ClearCache(); err != nil {
					return err
				}
				cmdutil.Infoln("DNS cache cleared.")
				return nil
			},
		},
		&cobra.Command{
			Use:   "check [hostname]",
			Short: "Check if a hostname is blocked by filtering rules",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				result, err := client.CheckHost(args[0])
				if err != nil {
					return err
				}
				return output.PrintJSON(result)
			},
		},
	)

	return cmd
}
