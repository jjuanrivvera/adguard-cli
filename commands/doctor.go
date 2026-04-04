package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jjuanrivvera/adguard-cli/internal/api"
	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/config"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run diagnostic checks against your AdGuard Home instance",
		RunE:  runDoctor,
	}
}

func runDoctor(cmd *cobra.Command, args []string) error {
	passed := 0

	// Check 1: Config
	cmdutil.Infof("[1/4] Checking configuration... ")
	inst, err := config.GetCurrentInstance()
	if err != nil || inst == nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "FAIL")
		cmdutil.Infoln("  No configuration found. Run 'adguard-home setup' first.")
		return nil
	}
	cmdutil.Infof("OK (%s)\n", inst.URL)
	passed++

	// Check 2: Connectivity
	cmdutil.Infof("[2/4] Checking connectivity... ")
	client := api.NewClient(inst.URL, inst.Username, inst.Password)
	if err := client.Ping(); err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "FAIL")
		cmdutil.Infof("  Cannot reach %s: %v\n", inst.URL, err)
		cmdutil.Infoln("  Hint: Is AdGuard Home running? Is the URL correct? Is Tailscale connected?")
		return nil
	}
	cmdutil.Infoln("OK")
	passed++

	// Check 3: Authentication
	cmdutil.Infof("[3/4] Checking authentication... ")
	status, err := client.GetStatus()
	if err != nil {
		fmt.Fprintln(cmd.ErrOrStderr(), "FAIL")
		cmdutil.Infof("  %v\n", err)
		cmdutil.Infoln("  Hint: Check credentials with 'adguard-home setup'")
		return nil
	}
	cmdutil.Infof("OK (v%s)\n", status.Version)
	passed++

	// Check 4: Protection
	cmdutil.Infof("[4/4] Checking protection status... ")
	if status.ProtectionEnabled {
		cmdutil.Infoln("OK (enabled)")
	} else {
		cmdutil.Infoln("WARNING (disabled)")
		cmdutil.Infoln("  DNS protection is off. Enable with 'adguard-home status enable'")
	}
	passed++

	cmdutil.Infof("\n%d/4 checks passed.\n", passed)
	return nil
}
