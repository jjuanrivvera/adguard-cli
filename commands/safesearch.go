package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newSafeSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "safesearch",
		Short: "Manage safe search enforcement across search engines",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "status",
			Short: "Show safe search configuration",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				s, err := client.GetSafeSearchStatus()
				if err != nil {
					return err
				}
				return output.Print(getFormat(), s,
					[]string{"Engine", "Enforced"},
					func() [][]string {
						return [][]string{
							{"Global", fmt.Sprintf("%t", s.Enabled)},
							{"Google", fmt.Sprintf("%t", s.Google)},
							{"Bing", fmt.Sprintf("%t", s.Bing)},
							{"DuckDuckGo", fmt.Sprintf("%t", s.DuckDuckGo)},
							{"YouTube", fmt.Sprintf("%t", s.YouTube)},
							{"Yandex", fmt.Sprintf("%t", s.Yandex)},
							{"Ecosia", fmt.Sprintf("%t", s.Ecosia)},
							{"Pixabay", fmt.Sprintf("%t", s.Pixabay)},
						}
					},
				)
			},
		},
	)

	return cmd
}
