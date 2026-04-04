package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

func newStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show DNS query statistics",
		RunE:  runStats,
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Reset all statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}
			if err := client.ResetStats(); err != nil {
				return err
			}
			cmdutil.Infoln("Statistics reset.")
			return nil
		},
	})

	return cmd
}

func runStats(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	stats, err := client.GetStats()
	if err != nil {
		return err
	}

	format := getFormat()

	blocked := 0
	if stats.NumDNSQueries > 0 {
		blocked = stats.NumBlockedFiltering * 100 / stats.NumDNSQueries
	}

	return output.Print(format, stats,
		[]string{"Metric", "Value"},
		func() [][]string {
			return [][]string{
				{"Total Queries", fmt.Sprintf("%d", stats.NumDNSQueries)},
				{"Blocked (filtering)", fmt.Sprintf("%d (%d%%)", stats.NumBlockedFiltering, blocked)},
				{"Blocked (safebrowsing)", fmt.Sprintf("%d", stats.NumReplacedSafebrowsing)},
				{"Blocked (parental)", fmt.Sprintf("%d", stats.NumReplacedParental)},
				{"Avg Processing Time", fmt.Sprintf("%.1fms", stats.AvgProcessingTime*1000)},
			}
		},
	)
}
