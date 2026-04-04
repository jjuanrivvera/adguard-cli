package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra/doc"

	"github.com/jjuanrivvera/adguard-cli/commands"
)

func main() {
	outputDir := "./docs/commands"

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	rootCmd := commands.NewRootCommand("dev", "none", "unknown")

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, ".md")
		return base + ".md"
	}

	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		name = strings.TrimSuffix(name, ".md")
		title := strings.ReplaceAll(name, "_", " ")
		return fmt.Sprintf("---\ntitle: %s\n---\n\n", title)
	}

	if err := doc.GenMarkdownTreeCustom(rootCmd, outputDir, filePrepender, linkHandler); err != nil {
		log.Fatalf("Failed to generate docs: %v", err)
	}

	fmt.Printf("Documentation generated in %s\n", outputDir)
}
