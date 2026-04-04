package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/api"
	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newRewritesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rewrites",
		Short: "Manage DNS rewrites",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List all DNS rewrites",
			RunE:  runRewritesList,
		},
		&cobra.Command{
			Use:   "add [domain] [answer]",
			Short: "Add a DNS rewrite (e.g., adguard-home rewrites add example.local 192.168.0.10)",
			Args:  cobra.ExactArgs(2),
			RunE:  runRewritesAdd,
		},
		&cobra.Command{
			Use:   "delete [domain] [answer]",
			Short: "Delete a DNS rewrite",
			Args:  cobra.ExactArgs(2),
			RunE:  runRewritesDelete,
		},
	)

	return cmd
}

func runRewritesList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	rewrites, err := client.GetRewrites()
	if err != nil {
		return err
	}

	format := getFormat()

	return output.Print(format, rewrites,
		[]string{"Domain", "Answer"},
		func() [][]string {
			var rows [][]string
			for _, r := range rewrites {
				rows = append(rows, []string{r.Domain, r.Answer})
			}
			return rows
		},
	)
}

func runRewritesAdd(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if err := client.AddRewrite(api.RewriteEntry{Domain: args[0], Answer: args[1]}); err != nil {
		return err
	}

	cmdutil.Infof("Rewrite added: %s -> %s\n", args[0], args[1])
	return nil
}

func runRewritesDelete(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	if err := client.DeleteRewrite(api.RewriteEntry{Domain: args[0], Answer: args[1]}); err != nil {
		return err
	}

	cmdutil.Infof("Rewrite deleted: %s -> %s\n", args[0], args[1])
	return nil
}
