package commands

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newAccessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access",
		Short: "Manage access control lists (allowed/disallowed clients and blocked domains)",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "Show access control configuration",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				a, err := client.GetAccessList()
				if err != nil {
					return err
				}
				return output.Print(getFormat(), a,
					[]string{"Category", "Entries"},
					func() [][]string {
						return [][]string{
							{"Allowed Clients", strings.Join(a.AllowedClients, ", ")},
							{"Disallowed Clients", strings.Join(a.DisallowedClients, ", ")},
							{"Blocked Hosts", strings.Join(a.BlockedHosts, ", ")},
						}
					},
				)
			},
		},
	)

	return cmd
}
