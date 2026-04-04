package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/jjuanrivvera/adguard-cli/internal/api"
	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
	"github.com/jjuanrivvera/adguard-cli/internal/config"
)

func newSetupCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "setup",
		Short: "Configure your AdGuard Home instance (interactive wizard)",
		RunE:  runSetup,
	}
}

func runSetup(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	cmdutil.Infoln("AdGuard Home CLI Setup")
	cmdutil.Infoln("======================")
	cmdutil.Infoln("")

	// URL
	cmdutil.Infof("AdGuard Home URL (e.g., http://192.168.0.105:8001): ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)
	if url == "" {
		cmdutil.Infoln("URL is required.")
		return nil
	}

	// Username
	cmdutil.Infof("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	// Password (masked input)
	cmdutil.Infof("Password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr) // newline after masked input
	if err != nil {
		// Fallback to plain input if terminal is not available (e.g., piped input)
		cmdutil.Infof("Password (visible): ")
		pw, _ := reader.ReadString('\n')
		passwordBytes = []byte(strings.TrimSpace(pw))
	}
	password := string(passwordBytes)

	// Test connection
	cmdutil.Infof("\nTesting connection to %s... ", url)
	client := api.NewClient(url, username, password)
	status, err := client.GetStatus()
	if err != nil {
		cmdutil.Infof("FAIL\n  %v\n", err)
		cmdutil.Infoln("Check the URL and credentials and try again.")
		return nil
	}
	cmdutil.Infof("OK (AdGuard Home v%s)\n", status.Version)

	// Instance name
	cmdutil.Infof("Instance name (default: default): ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = "default"
	}

	// Save
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	cfg.Instances[name] = config.Instance{
		URL:      url,
		Username: username,
	}
	cfg.CurrentInstance = name

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	// Store password in system keyring or encrypted file
	if err := config.SaveCredentials(name, password); err != nil {
		return fmt.Errorf("saving credentials: %w", err)
	}

	dir, _ := config.ConfigDir()
	cmdutil.Infof("\nConfiguration saved to %s/config.yaml\n", dir)
	cmdutil.Infoln("Password stored in system keyring (or encrypted file as fallback).")
	cmdutil.Infoln("Run 'adguard-home doctor' to verify everything works.")

	return nil
}
