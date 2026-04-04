package main

import (
	"os"

	"github.com/jjuanrivvera/adguard-cli/commands"
	"github.com/jjuanrivvera/adguard-cli/internal/cmdutil"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	root := commands.NewRootCommand(version, commit, date)

	if err := root.Execute(); err != nil {
		cmdutil.HandleError(err)
		os.Exit(1)
	}
}
