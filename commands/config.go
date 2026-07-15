package commands

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/config"
	clierrors "github.com/jjuanrivvera/adguard-cli/internal/errors"
	"github.com/jjuanrivvera/adguard-cli/internal/output"
)

// newConfigCmd manages the configured AdGuard Home instances (the CLI's "profiles").
// The instance data model always supported several instances; this command is the
// fleet-standard surface to list them, switch the active one, and remove one — without
// re-running `setup` or hand-editing config.yaml.
func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "View and manage configured instances (profiles)",
		Long: `Inspect the config file and manage AdGuard Home instances: list them, switch
the active instance, or remove one. Add an instance with ` + "`adguard-home setup`" + `.`,
	}
	cmd.AddCommand(newConfigViewCmd(), newConfigListCmd(), newConfigUseCmd(), newConfigRemoveCmd(), newConfigPathCmd())
	return cmd
}

func loadConfigOrEmpty() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	if cfg.Instances == nil {
		cfg.Instances = map[string]config.Instance{}
	}
	return cfg, nil
}

func sortedInstanceNames(cfg *config.Config) []string {
	names := make([]string, 0, len(cfg.Instances))
	for n := range cfg.Instances {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func unknownInstance(name string) error {
	return clierrors.WithHint(
		clierrors.New(clierrors.NotFound, fmt.Sprintf("instance %q not found", name)),
		"Run `adguard-home config list` to see configured instances.",
	)
}

func newConfigViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Show all configured instances and which one is active",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := loadConfigOrEmpty()
			if err != nil {
				return err
			}
			names := sortedInstanceNames(cfg)
			if len(names) == 0 {
				cmdutil.Infoln("No instances configured yet — run `adguard-home setup`.")
				return nil
			}
			type instanceView struct {
				Name     string `json:"name" yaml:"name"`
				URL      string `json:"url" yaml:"url"`
				Username string `json:"username" yaml:"username"`
				Active   bool   `json:"active" yaml:"active"`
			}
			views := make([]instanceView, 0, len(names))
			for _, n := range names {
				inst := cfg.Instances[n]
				views = append(views, instanceView{n, inst.URL, inst.Username, n == cfg.CurrentInstance})
			}
			return output.Print(getFormat(), views,
				[]string{"Name", "URL", "Username", "Active"},
				func() [][]string {
					rows := make([][]string, 0, len(views))
					for _, v := range views {
						active := ""
						if v.Active {
							active = "*"
						}
						rows = append(rows, []string{v.Name, v.URL, v.Username, active})
					}
					return rows
				},
			)
		},
	}
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured instance names (active marked with *)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := loadConfigOrEmpty()
			if err != nil {
				return err
			}
			names := sortedInstanceNames(cfg)
			if len(names) == 0 {
				cmdutil.Infoln("No instances configured yet — run `adguard-home setup`.")
				return nil
			}
			for _, n := range names {
				marker := "  "
				if n == cfg.CurrentInstance {
					marker = "* "
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s%s\n", marker, n)
			}
			return nil
		},
	}
}

func newConfigUseCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use <instance>",
		Short: "Set the active instance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfigOrEmpty()
			if err != nil {
				return err
			}
			name := args[0]
			if _, ok := cfg.Instances[name]; !ok {
				return unknownInstance(name)
			}
			cfg.CurrentInstance = name
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}
			cmdutil.Infof("Active instance is now %q\n", name)
			return nil
		},
	}
}

func newConfigRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <instance>",
		Aliases: []string{"rm"},
		Short:   "Remove an instance and its stored credentials",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfigOrEmpty()
			if err != nil {
				return err
			}
			name := args[0]
			if _, ok := cfg.Instances[name]; !ok {
				return unknownInstance(name)
			}
			delete(cfg.Instances, name)
			// If the active instance was removed, repoint it to any remaining one so the
			// config never dangles at a deleted instance.
			if cfg.CurrentInstance == name {
				cfg.CurrentInstance = ""
				if remaining := sortedInstanceNames(cfg); len(remaining) > 0 {
					cfg.CurrentInstance = remaining[0]
				}
			}
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("saving config: %w", err)
			}
			// Best-effort: the password may live only in the keyring, or not at all.
			_ = config.DeleteCredentials(name)
			cmdutil.Infof("Removed instance %q\n", name)
			if cfg.CurrentInstance != "" {
				cmdutil.Infof("Active instance is now %q\n", cfg.CurrentInstance)
			}
			return nil
		},
	}
}

func newConfigPathCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Print the config file path",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			p, err := config.Path()
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), p)
			return nil
		},
	}
}
