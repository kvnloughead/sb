package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration for the sb CLI tool.
type Config struct {
	Repo  string            `yaml:"repo"`  // Absolute path to the git repository
	Slugs map[string]string `yaml:"slugs"` // Map of slug to branch name
}

// LoadConfig loads the configuration from the given path, or the default
// location (~/.config/sb.yaml). If the config file can't be found or is
// invalid, the (sensible) default values from defaults.go are used.
func LoadConfig() (*Config, error) {
	configPath := os.Getenv("SB_CONFIG")
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			// If we can't get the home dir, just use defaults
			return &DefaultConfig, nil
		}
		configPath = filepath.Join(home, ".config", "sb.yaml")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		// If config file not found, use defaults
		return &DefaultConfig, nil
	}
	var cfg Config = DefaultConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		// If config file is invalid, use defaults
		return &DefaultConfig, nil
	}
	// Fill missing values from defaults
	if cfg.Repo == "" {
		cfg.Repo = DefaultConfig.Repo
	}
	if len(cfg.Slugs) == 0 {
		cfg.Slugs = DefaultConfig.Slugs
	}
	return &cfg, nil
}
