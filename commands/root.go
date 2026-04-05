package commands

import (
	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/api"
	"github.com/jjuanrivvera/adguard-cli/internal/config"
	clierrors "github.com/jjuanrivvera/adguard-cli/internal/errors"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

type GlobalFlags struct {
	OutputFormat string
	Instance     string
}

var flags GlobalFlags

func NewRootCommand(version, commit, date string) *cobra.Command {
	root := &cobra.Command{
		Use:   "adguard-home",
		Short: "The missing CLI for AdGuard Home",
		Long:  "A command-line interface for managing AdGuard Home DNS filtering.\nSupports clients, blocked services, DNS rewrites, query logs, filters, and more.",
		Version: version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVarP(&flags.OutputFormat, "output", "o", "table", "Output format: table, json, yaml")
	root.PersistentFlags().StringVar(&flags.Instance, "instance", "", "AdGuard Home instance name from config (default: current_instance)")

	root.AddCommand(
		newStatusCmd(),
		newStatsCmd(),
		newClientsCmd(),
		newServicesCmd(),
		newRewritesCmd(),
		newLogCmd(),
		newFiltersCmd(),
		newDHCPCmd(),
		newTLSCmd(),
		newDNSCmd(),
		newSafeBrowsingCmd(),
		newParentalCmd(),
		newSafeSearchCmd(),
		newAccessCmd(),
		newVersionCmd(),
		newUpdateCmd(),
		newDoctorCmd(),
		newSetupCmd(),
		newMCPCmd(),
	)

	return root
}

func getFormat() output.Format {
	return output.ParseFormat(flags.OutputFormat)
}

func getClient() (*api.Client, error) {
	inst, err := config.GetCurrentInstance()
	if err != nil {
		return nil, clierrors.Wrap(clierrors.ConfigError, "loading config", err)
	}
	if inst == nil {
		return nil, clierrors.ConfigNotFound()
	}

	// Override instance if --instance flag was provided
	if flags.Instance != "" {
		named, err := config.GetNamedInstance(flags.Instance)
		if err != nil {
			return nil, clierrors.Wrap(clierrors.ConfigError, "loading instance", err)
		}
		if named == nil {
			return nil, clierrors.WithHint(
				clierrors.New(clierrors.NotFound, "instance '"+flags.Instance+"' not found"),
				"Check available instances in ~/.adguard-cli/config.yaml",
			)
		}
		inst = named
	}

	return api.NewClient(inst.URL, inst.Username, inst.Password), nil
}
