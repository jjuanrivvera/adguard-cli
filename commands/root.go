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
	Profile      string
}

var flags GlobalFlags

func NewRootCommand(version, commit, date string) *cobra.Command {
	root := &cobra.Command{
		Use:           "adguard-home",
		Short:         "The missing CLI for AdGuard Home",
		Long:          "A command-line interface for managing AdGuard Home DNS filtering.\nSupports clients, blocked services, DNS rewrites, query logs, filters, and more.",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().StringVarP(&flags.OutputFormat, "output", "o", "table", "Output format: table, json, yaml")
	root.PersistentFlags().StringVar(&flags.Instance, "instance", "", "AdGuard Home instance name from config (default: current_instance)")
	root.PersistentFlags().StringVar(&flags.Profile, "profile", "", "alias for --instance")
	_ = root.PersistentFlags().MarkHidden("profile")

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
		newServerCmd(),
		newSelfUpdateCmd(version),
		newDoctorCmd(),
		newSetupCmd(),
		newConfigCmd(),
		newMCPCmd(),
	)

	// Must run after the tree is complete: ophis reads cmd.Annotations when the mcp
	// subcommand runs, and a host in read-only mode drops a server with no read-only tool.
	applyMCPAnnotations(root)

	return root
}

func getFormat() output.Format {
	return output.ParseFormat(flags.OutputFormat)
}

// selectedInstance returns the instance requested for this invocation via --instance,
// falling back to its --profile alias (kept for vocabulary parity with the rest of the CLI
// fleet). Empty means "use the configured current_instance".
func selectedInstance() string {
	if flags.Instance != "" {
		return flags.Instance
	}
	return flags.Profile
}

func getClient() (*api.Client, error) {
	inst, err := config.GetCurrentInstance()
	if err != nil {
		return nil, clierrors.Wrap(clierrors.ConfigError, "loading config", err)
	}
	if inst == nil {
		return nil, clierrors.ConfigNotFound()
	}

	// Override instance if --instance (or its --profile alias) was provided
	if name := selectedInstance(); name != "" {
		named, err := config.GetNamedInstance(name)
		if err != nil {
			return nil, clierrors.Wrap(clierrors.ConfigError, "loading instance", err)
		}
		if named == nil {
			return nil, clierrors.WithHint(
				clierrors.New(clierrors.NotFound, "instance '"+name+"' not found"),
				"Run `adguard-home config list` to see configured instances.",
			)
		}
		inst = named
	}

	return api.NewClient(inst.URL, inst.Username, inst.Password), nil
}
