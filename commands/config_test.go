package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jjuanrivvera/adguard-cli/internal/config"
)

// seedTwoInstances writes a config with two instances ("home" active, "office") into an
// isolated HOME. Credentials are intentionally not stored — the config command surface
// (list/use/remove/view) operates on config.yaml, not the keyring.
func seedTwoInstances(t *testing.T) {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.Instances["home"] = config.Instance{URL: "http://192.168.0.5:8001", Username: "admin"}
	cfg.Instances["office"] = config.Instance{URL: "http://10.0.0.2:3000", Username: "ops"}
	cfg.CurrentInstance = "home"
	require.NoError(t, config.Save(cfg))
}

func runRootCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	root := NewRootCommand("test", "none", "unknown")
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(args)
	return capture(&buf, root.Execute())
}

func capture(buf *bytes.Buffer, err error) (string, error) { return buf.String(), err }

func TestConfigList_MarksActive(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	seedTwoInstances(t)
	out, err := runRootCmd(t, "config", "list")
	require.NoError(t, err)
	assert.Contains(t, out, "* home", "active instance should be marked with *")
	assert.Contains(t, out, "  office", "inactive instance should be unmarked")
}

func TestConfigUse_SwitchesActive(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	seedTwoInstances(t)
	_, err := runRootCmd(t, "config", "use", "office")
	require.NoError(t, err)

	cfg, err := config.Load()
	require.NoError(t, err)
	assert.Equal(t, "office", cfg.CurrentInstance)
}

func TestConfigUse_UnknownInstanceErrors(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	seedTwoInstances(t)
	_, err := runRootCmd(t, "config", "use", "nope")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestConfigRemove_RepointsActive(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	seedTwoInstances(t)
	// Remove the ACTIVE instance; current_instance must repoint to the remaining one.
	_, err := runRootCmd(t, "config", "remove", "home")
	require.NoError(t, err)

	cfg, err := config.Load()
	require.NoError(t, err)
	_, exists := cfg.Instances["home"]
	assert.False(t, exists, "removed instance should be gone")
	assert.Equal(t, "office", cfg.CurrentInstance, "active should repoint to the remaining instance")
}

func TestProfileFlagIsHiddenAliasForInstance(t *testing.T) {
	root := NewRootCommand("test", "none", "unknown")
	f := root.PersistentFlags().Lookup("profile")
	require.NotNil(t, f, "--profile alias should be registered")
	assert.True(t, f.Hidden, "--profile should be hidden (alias for --instance)")
}
