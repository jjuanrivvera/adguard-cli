package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newLogCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "log",
		Short: "View DNS query log",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}

			ql, err := client.GetQueryLog(limit)
			if err != nil {
				return err
			}

			format := getFormat()

			return output.Print(format, ql.Data,
				[]string{"Time", "Client", "Domain", "Type", "Status", "Reason"},
				func() [][]string {
					var rows [][]string
					for _, entry := range ql.Data {
						domain := ""
						qtype := ""
						if q, ok := entry.Question["name"]; ok {
							if s, ok := q.(string); ok {
								domain = s
							}
						}
						if q, ok := entry.Question["type"]; ok {
							if s, ok := q.(string); ok {
								qtype = s
							}
						}
						rows = append(rows, []string{
							entry.Time,
							entry.Client,
							domain,
							qtype,
							entry.Status,
							entry.Reason,
						})
					}
					return rows
				},
			)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 25, "Number of log entries to show")

	return cmd
}
