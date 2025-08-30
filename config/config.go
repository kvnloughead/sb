package config

import (
	"fmt"
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
// location (~/.config/sb.yaml). Returns an error if no config file is found.
func LoadConfig() (*Config, error) {
	configPath := os.Getenv("SB_CONFIG")
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("unable to get home directory: %w", err)
		}
		configPath = filepath.Join(home, ".config", "sb.yaml")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("config file not found at %s (run 'sb install' to create it): %w", configPath, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if cfg.Repo == "" {
		return nil, fmt.Errorf("repo field is required in config file")
	}

	return &cfg, nil
}
