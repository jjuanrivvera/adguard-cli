package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.NotNil(t, cfg.Instances)
	assert.Equal(t, "default", cfg.CurrentInstance)
	assert.Equal(t, "table", cfg.Output.Format)
	assert.Equal(t, "auto", cfg.Output.Color)
}

func TestSaveAndLoad(t *testing.T) {
	// Use temp dir as home
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfg := DefaultConfig()
	cfg.Instances["test"] = Instance{
		URL:      "http://192.168.0.1:8001",
		Username: "admin",
	}
	cfg.CurrentInstance = "test"

	err := Save(cfg)
	require.NoError(t, err)

	// Verify file exists with correct permissions
	path := filepath.Join(tmpDir, configDir, configFile)
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	// Load and verify
	loaded, err := Load()
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "test", loaded.CurrentInstance)
	assert.Equal(t, "http://192.168.0.1:8001", loaded.Instances["test"].URL)
	assert.Equal(t, "admin", loaded.Instances["test"].Username)
}

func TestLoad_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfg, err := Load()
	assert.NoError(t, err)
	assert.Nil(t, cfg)
}

func TestGetCurrentInstance_EnvVars(t *testing.T) {
	t.Setenv("ADGUARD_URL", "http://10.0.0.1:3000")
	t.Setenv("ADGUARD_USERNAME", "testuser")
	t.Setenv("ADGUARD_PASSWORD", "testpass")

	inst, err := GetCurrentInstance()
	require.NoError(t, err)
	require.NotNil(t, inst)
	assert.Equal(t, "http://10.0.0.1:3000", inst.URL)
	assert.Equal(t, "testuser", inst.Username)
	assert.Equal(t, "testpass", inst.Password)
}

func TestGetCurrentInstance_EnvVarsPartial(t *testing.T) {
	// Only URL set — should fall through to config
	t.Setenv("ADGUARD_URL", "http://10.0.0.1:3000")
	t.Setenv("ADGUARD_USERNAME", "")
	t.Setenv("ADGUARD_PASSWORD", "")

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// No config file → nil
	inst, err := GetCurrentInstance()
	assert.NoError(t, err)
	assert.Nil(t, inst)
}

func TestGetCurrentInstance_InstanceNotFound(t *testing.T) {
	t.Setenv("ADGUARD_URL", "")
	t.Setenv("ADGUARD_USERNAME", "")
	t.Setenv("ADGUARD_PASSWORD", "")

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfg := DefaultConfig()
	cfg.CurrentInstance = "nonexistent"
	err := Save(cfg)
	require.NoError(t, err)

	_, err = GetCurrentInstance()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent")
}

func TestConfigDir(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	dir, err := ConfigDir()
	require.NoError(t, err)
	assert.Equal(t, filepath.Join(tmpDir, ".adguard-cli"), dir)
}

func TestNoPasswordInYAML(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfg := DefaultConfig()
	cfg.Instances["secure"] = Instance{
		URL:      "http://localhost:8001",
		Username: "admin",
	}
	err := Save(cfg)
	require.NoError(t, err)

	// Read raw file and verify no password field
	path := filepath.Join(tmpDir, configDir, configFile)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.NotContains(t, string(data), "password")
}
