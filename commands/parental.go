package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newParentalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parental",
		Short: "Manage parental control",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "status",
			Short: "Show parental control status",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				s, err := client.GetParentalStatus()
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
						return [][]string{{"Parental Control", status}}
					},
				)
			},
		},
		&cobra.Command{
			Use:   "enable",
			Short: "Enable parental control",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.SetParental(true); err != nil {
					return err
				}
				cmdutil.Infoln("Parental control enabled.")
				return nil
			},
		},
		&cobra.Command{
			Use:   "disable",
			Short: "Disable parental control",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.SetParental(false); err != nil {
					return err
				}
				cmdutil.Infoln("Parental control disabled.")
				return nil
			},
		},
	)

	return cmd
}
