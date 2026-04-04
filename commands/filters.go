package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newFiltersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filters",
		Short: "Manage DNS filter lists",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List all filter lists",
			RunE:  runFiltersList,
		},
		&cobra.Command{
			Use:   "add [name] [url]",
			Short: "Add a new filter list",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.AddFilter(args[0], args[1], true); err != nil {
					return fmt.Errorf("adding filter: %w", err)
				}
				cmdutil.Infof("Filter %q added.\n", args[0])
				return nil
			},
		},
		&cobra.Command{
			Use:   "remove [url]",
			Short: "Remove a filter list by URL",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.RemoveFilter(args[0]); err != nil {
					return fmt.Errorf("removing filter: %w", err)
				}
				cmdutil.Infoln("Filter removed.")
				return nil
			},
		},
		&cobra.Command{
			Use:   "refresh",
			Short: "Refresh all filter lists",
			RunE: func(cmd *cobra.Command, args []string) error {
				client, err := getClient()
				if err != nil {
					return err
				}
				if err := client.RefreshFilters(); err != nil {
					return err
				}
				cmdutil.Infoln("Filters refreshed.")
				return nil
			},
		},
	)

	return cmd
}

func runFiltersList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	fs, err := client.GetFiltering()
	if err != nil {
		return err
	}

	format := getFormat()

	return output.Print(format, fs.Filters,
		[]string{"Name", "Rules", "Enabled", "Last Updated", "URL"},
		func() [][]string {
			var rows [][]string
			for _, f := range fs.Filters {
				enabled := "off"
				if f.Enabled {
					enabled = "on"
				}
				rows = append(rows, []string{
					f.Name,
					fmt.Sprintf("%d", f.RulesCount),
					enabled,
					f.LastUpdated,
					f.URL,
				})
			}
			return rows
		},
	)
}
