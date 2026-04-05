package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPCommand_Registered(t *testing.T) {
	root := NewRootCommand("test", "none", "unknown")

	// Verify mcp subcommand exists
	var mcpFound bool
	for _, cmd := range root.Commands() {
		if cmd.Use == "mcp" {
			mcpFound = true
			break
		}
	}
	assert.True(t, mcpFound, "mcp command should be registered on root")
}

func TestMCPCommand_HasSubcommands(t *testing.T) {
	root := NewRootCommand("test", "none", "unknown")

	var mcpCmd *cobra.Command
	for _, cmd := range root.Commands() {
		if cmd.Use == "mcp" {
			mcpCmd = cmd
			break
		}
	}
	require.NotNil(t, mcpCmd, "mcp command not found")

	expectedSubs := []string{"start", "stream", "tools", "claude", "vscode", "cursor"}
	for _, name := range expectedSubs {
		found := false
		for _, sub := range mcpCmd.Commands() {
			if sub.Use == name || sub.Name() == name {
				found = true
				break
			}
		}
		assert.True(t, found, "mcp should have subcommand %q", name)
	}
}

func TestMCPTools_ExportAndValidate(t *testing.T) {
	root := NewRootCommand("test", "none", "unknown")

	// Find and execute mcp tools command
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { _ = os.Chdir(origDir) }()

	root.SetArgs([]string{"mcp", "tools"})
	err := root.Execute()
	require.NoError(t, err)

	// Read exported tools
	data, err := os.ReadFile(filepath.Join(tmpDir, "mcp-tools.json"))
	require.NoError(t, err)

	var tools []map[string]any
	err = json.Unmarshal(data, &tools)
	require.NoError(t, err)

	// Should have tools exported
	assert.Greater(t, len(tools), 0, "should export at least one tool")

	// All tools should have the adguard prefix
	for _, tool := range tools {
		name, ok := tool["name"].(string)
		require.True(t, ok, "tool should have a name")
		assert.Contains(t, name, "adguard", "tool name should have adguard prefix: %s", name)
	}
}

func TestMCPTools_DestructiveCommandsExcluded(t *testing.T) {
	root := NewRootCommand("test", "none", "unknown")

	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { _ = os.Chdir(origDir) }()

	root.SetArgs([]string{"mcp", "tools"})
	err := root.Execute()
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(tmpDir, "mcp-tools.json"))
	require.NoError(t, err)

	var tools []map[string]any
	err = json.Unmarshal(data, &tools)
	require.NoError(t, err)

	// Verify destructive commands are excluded
	for _, tool := range tools {
		name := tool["name"].(string)
		assert.NotContains(t, name, "reset", "destructive command 'reset' should be excluded: %s", name)
		assert.NotContains(t, name, "update", "destructive command 'update' should be excluded: %s", name)
	}
}

func TestMCPTools_SensitiveFlagsExcluded(t *testing.T) {
	root := NewRootCommand("test", "none", "unknown")

	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() { _ = os.Chdir(origDir) }()

	root.SetArgs([]string{"mcp", "tools"})
	err := root.Execute()
	require.NoError(t, err)

	data, err := os.ReadFile(filepath.Join(tmpDir, "mcp-tools.json"))
	require.NoError(t, err)

	var tools []map[string]any
	err = json.Unmarshal(data, &tools)
	require.NoError(t, err)

	// Check that --instance flag is not in any tool's schema
	for _, tool := range tools {
		schema, ok := tool["inputSchema"].(map[string]any)
		if !ok {
			continue
		}
		props, ok := schema["properties"].(map[string]any)
		if !ok {
			continue
		}
		flags, ok := props["flags"].(map[string]any)
		if !ok {
			continue
		}
		flagProps, ok := flags["properties"].(map[string]any)
		if !ok {
			continue
		}
		_, hasInstance := flagProps["instance"]
		assert.False(t, hasInstance, "tool %s should not expose --instance flag", tool["name"])
	}
}
