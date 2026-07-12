package commands

import (
	"bufio"
	"fmt"
	"io"
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
	passwordStr, err := readSecretRaw(os.Stdin)
	fmt.Fprintln(os.Stderr) // newline after masked input
	if err != nil {
		// Fallback to plain input if terminal is not available (e.g., piped input)
		cmdutil.Infof("Password (visible): ")
		pw, _ := reader.ReadString('\n')
		passwordStr = strings.TrimSpace(pw)
	}
	password := sanitizeSecret(passwordStr)

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

// sanitizeSecret strips terminal bracketed-paste markers (ESC[200~ … ESC[201~) and trims
// surrounding whitespace. With bracketed paste enabled, a raw read (unlike the shell's line
// editor) receives those wrappers around pasted text; left in they corrupt a pasted secret so
// it fails auth. Stripping them fixes the common "typing works, pasting fails".
func sanitizeSecret(s string) string {
	s = strings.ReplaceAll(s, "\x1b[200~", "")
	s = strings.ReplaceAll(s, "\x1b[201~", "")
	return strings.TrimSpace(s)
}

// readSecretRaw puts the terminal in raw mode (no echo, no line-length limit) and reads one line.
// term.ReadPassword instead reads in CANONICAL mode, capped at MAX_CANON (1024 bytes on macOS):
// pasting a longer secret fills the line buffer and the terminal BLOCKS — the "prompt hangs until
// Ctrl-C" bug. Raw mode has no such limit. On a pipe MakeRaw fails, so the caller falls back.
func readSecretRaw(f *os.File) (string, error) {
	fd := int(f.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() { _ = term.Restore(fd, oldState) }()
	return scanSecretLine(f)
}

// scanSecretLine reads bytes until CR/LF with no line-length limit. Ctrl-C cancels; Backspace/DEL
// edits. Split out so the byte handling is testable without a real terminal.
func scanSecretLine(r io.Reader) (string, error) {
	var buf []byte
	chunk := make([]byte, 256)
	for {
		n, readErr := r.Read(chunk)
		for i := 0; i < n; i++ {
			switch c := chunk[i]; c {
			case '\r', '\n':
				return string(buf), nil
			case 3: // Ctrl-C
				return "", fmt.Errorf("cancelled")
			case 127, 8: // DEL / Backspace
				if len(buf) > 0 {
					buf = buf[:len(buf)-1]
				}
			default:
				buf = append(buf, c)
			}
		}
		if readErr != nil {
			if len(buf) == 0 {
				return "", readErr
			}
			return string(buf), nil
		}
	}
}
