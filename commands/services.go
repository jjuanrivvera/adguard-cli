package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newServicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "Manage globally blocked services",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List all available services and their block status",
			RunE:  runServicesList,
		},
		&cobra.Command{
			Use:   "blocked",
			Short: "Show only blocked services",
			RunE:  runServicesBlocked,
		},
		&cobra.Command{
			Use:   "block [service1,service2,...]",
			Short: "Block one or more services",
			Args:  cobra.ExactArgs(1),
			RunE:  runServicesBlock,
		},
		&cobra.Command{
			Use:   "unblock [service1,service2,...]",
			Short: "Unblock one or more services",
			Args:  cobra.ExactArgs(1),
			RunE:  runServicesUnblock,
		},
	)

	return cmd
}

func runServicesList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	all, err := client.GetAllServices()
	if err != nil {
		return err
	}

	blocked, err := client.GetBlockedServices()
	if err != nil {
		return err
	}

	blockedSet := make(map[string]bool)
	for _, id := range blocked.IDs {
		blockedSet[id] = true
	}

	format := getFormat()

	type serviceRow struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Blocked bool   `json:"blocked"`
	}
	var data []serviceRow
	for _, s := range all {
		data = append(data, serviceRow{ID: s.ID, Name: s.Name, Blocked: blockedSet[s.ID]})
	}

	return output.Print(format, data,
		[]string{"ID", "Name", "Blocked"},
		func() [][]string {
			var rows [][]string
			for _, s := range data {
				status := ""
				if s.Blocked {
					status = "BLOCKED"
				}
				rows = append(rows, []string{s.ID, s.Name, status})
			}
			return rows
		},
	)
}

func runServicesBlocked(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	blocked, err := client.GetBlockedServices()
	if err != nil {
		return err
	}

	format := getFormat()

	return output.Print(format, blocked.IDs,
		[]string{"Blocked Service"},
		func() [][]string {
			var rows [][]string
			for _, id := range blocked.IDs {
				rows = append(rows, []string{id})
			}
			return rows
		},
	)
}

func runServicesBlock(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	toBlock := strings.Split(args[0], ",")

	current, err := client.GetBlockedServices()
	if err != nil {
		return err
	}

	existing := make(map[string]bool)
	for _, id := range current.IDs {
		existing[id] = true
	}
	for _, id := range toBlock {
		existing[id] = true
	}

	var merged []string
	for id := range existing {
		merged = append(merged, id)
	}

	current.IDs = merged
	if err := client.SetBlockedServices(*current); err != nil {
		return fmt.Errorf("updating blocked services: %w", err)
	}

	cmdutil.Infof("Blocked: %s\n", strings.Join(toBlock, ", "))
	return nil
}

func runServicesUnblock(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	toUnblock := strings.Split(args[0], ",")
	unblockSet := make(map[string]bool)
	for _, id := range toUnblock {
		unblockSet[id] = true
	}

	current, err := client.GetBlockedServices()
	if err != nil {
		return err
	}

	var filtered []string
	for _, id := range current.IDs {
		if !unblockSet[id] {
			filtered = append(filtered, id)
		}
	}

	current.IDs = filtered
	if err := client.SetBlockedServices(*current); err != nil {
		return fmt.Errorf("updating blocked services: %w", err)
	}

	cmdutil.Infof("Unblocked: %s\n", strings.Join(toUnblock, ", "))
	return nil
}
