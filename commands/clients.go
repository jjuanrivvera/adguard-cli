package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/api"
	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	clierrors "github.com/jjuanrivvera/adguard-cli/internal/errors"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newClientsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clients",
		Short: "Manage AdGuard Home clients",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all configured clients",
		RunE:  runClientsList,
	}

	findCmd := &cobra.Command{
		Use:   "find [ip]",
		Short: "Find which client an IP belongs to",
		Args:  cobra.ExactArgs(1),
		RunE:  runClientsFind,
	}

	addCmd := &cobra.Command{
		Use:   "add [name] [ip1,ip2,...]",
		Short: "Add a new client",
		Args:  cobra.ExactArgs(2),
		RunE:  runClientsAdd,
	}

	deleteCmd := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a client by name",
		Args:  cobra.ExactArgs(1),
		RunE:  runClientsDelete,
	}

	cmd.AddCommand(listCmd, findCmd, addCmd, deleteCmd)
	return cmd
}

func runClientsList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	resp, err := client.GetClients()
	if err != nil {
		return err
	}

	format := getFormat()

	return output.Print(format, resp.Clients,
		[]string{"Name", "IDs", "Blocked Services", "Filtering"},
		func() [][]string {
			var rows [][]string
			for _, c := range resp.Clients {
				services := "global"
				if !c.UseGlobalBlockedServices {
					if len(c.BlockedServices) == 0 {
						services = "none"
					} else {
						services = strings.Join(c.BlockedServices, ", ")
					}
				}
				filtering := "off"
				if c.FilteringEnabled {
					filtering = "on"
				}
				rows = append(rows, []string{
					c.Name,
					strings.Join(c.IDs, ", "),
					services,
					filtering,
				})
			}
			return rows
		},
	)
}

func runClientsFind(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	result, err := client.FindClient(args[0])
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return clierrors.ClientNotFound(args[0])
	}

	return output.PrintJSON(result)
}

func runClientsAdd(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	ids := strings.Split(args[1], ",")
	entry := api.ClientEntry{
		Name:                    args[0],
		IDs:                     ids,
		UseGlobalSettings:       true,
		UseGlobalBlockedServices: true,
		FilteringEnabled:        true,
		SafebrowsingEnabled:     true,
	}

	if err := client.AddClient(entry); err != nil {
		return fmt.Errorf("adding client: %w", err)
	}

	cmdutil.Infof("Client %q added with IDs: %s\n", args[0], strings.Join(ids, ", "))
	return nil
}

func runClientsDelete(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if err := client.DeleteClient(args[0]); err != nil {
		return fmt.Errorf("deleting client: %w", err)
	}

	cmdutil.Infof("Client %q deleted.\n", args[0])
	return nil
}
