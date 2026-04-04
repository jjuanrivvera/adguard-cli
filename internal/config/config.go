package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configDir  = ".adguard-cli"
	configFile = "config.yaml"
)

type Instance struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type OutputConfig struct {
	Format string `yaml:"format"`
	Color  string `yaml:"color"`
}

type Config struct {
	Instances       map[string]Instance `yaml:"instances"`
	CurrentInstance string              `yaml:"current_instance"`
	Output          OutputConfig        `yaml:"output"`
}

func DefaultConfig() *Config {
	return &Config{
		Instances:       make(map[string]Instance),
		CurrentInstance: "default",
		Output: OutputConfig{
			Format: "table",
			Color:  "auto",
		},
	}
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir, configFile), nil
}

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configDir), nil
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return cfg, nil
}

func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// GetCurrentInstance returns the active instance config, checking env vars first.
func GetCurrentInstance() (*Instance, error) {
	// Environment variables take precedence
	url := os.Getenv("ADGUARD_URL")
	user := os.Getenv("ADGUARD_USERNAME")
	pass := os.Getenv("ADGUARD_PASSWORD")

	if url != "" && user != "" && pass != "" {
		return &Instance{URL: url, Username: user, Password: pass}, nil
	}

	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	inst, ok := cfg.Instances[cfg.CurrentInstance]
	if !ok {
		return nil, fmt.Errorf("instance %q not found in config", cfg.CurrentInstance)
	}

	return &inst, nil
}
