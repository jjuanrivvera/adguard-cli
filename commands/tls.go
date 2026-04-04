package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newTLSCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tls",
		Short: "Manage TLS/HTTPS configuration",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "status",
			Short: "Show TLS configuration status",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				s, err := client.GetTLSStatus()
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
							{"Server Name", s.ServerName},
							{"HTTPS Port", fmt.Sprintf("%d", s.PortHTTPS)},
							{"DNS-over-TLS Port", fmt.Sprintf("%d", s.PortDNSOverTLS)},
							{"DNS-over-QUIC Port", fmt.Sprintf("%d", s.PortDNSOverQUIC)},
							{"Force HTTPS", fmt.Sprintf("%t", s.ForceHTTPS)},
							{"Valid Cert", fmt.Sprintf("%t", s.ValidCert)},
							{"Valid Key", fmt.Sprintf("%t", s.ValidKey)},
							{"Valid Pair", fmt.Sprintf("%t", s.ValidPair)},
							{"Cert Path", s.CertificatePath},
							{"Key Path", s.PrivateKeyPath},
						}
					},
				)
			},
		},
	)

	return cmd
}
