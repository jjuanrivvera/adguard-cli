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
	// Password is stored in system keyring or encrypted file, never in config YAML.
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

// ResolvedInstance holds the full connection info including the password.
type ResolvedInstance struct {
	URL      string
	Username string
	Password string
}

// GetCurrentInstance returns the active instance config, checking env vars first.
func GetCurrentInstance() (*ResolvedInstance, error) {
	// Environment variables take precedence
	url := os.Getenv("ADGUARD_URL")
	user := os.Getenv("ADGUARD_USERNAME")
	pass := os.Getenv("ADGUARD_PASSWORD")

	if url != "" && user != "" && pass != "" {
		return &ResolvedInstance{URL: url, Username: user, Password: pass}, nil
	}

	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	instanceName := cfg.CurrentInstance
	inst, ok := cfg.Instances[instanceName]
	if !ok {
		return nil, fmt.Errorf("instance %q not found in config", instanceName)
	}

	// Get password from credential store
	store := NewCredentialStore()
	password, err := store.Get(instanceName)
	if err != nil {
		return nil, fmt.Errorf("retrieving credentials: %w", err)
	}

	return &ResolvedInstance{URL: inst.URL, Username: inst.Username, Password: password}, nil
}

// GetNamedInstance returns a specific instance by name from config + credential store.
func GetNamedInstance(name string) (*ResolvedInstance, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}

	inst, ok := cfg.Instances[name]
	if !ok {
		return nil, nil
	}

	store := NewCredentialStore()
	password, err := store.Get(name)
	if err != nil {
		return nil, fmt.Errorf("retrieving credentials: %w", err)
	}

	return &ResolvedInstance{URL: inst.URL, Username: inst.Username, Password: password}, nil
}

// SaveCredentials stores a password in the credential store.
func SaveCredentials(instance, password string) error {
	store := NewCredentialStore()
	return store.Set(instance, password)
}
