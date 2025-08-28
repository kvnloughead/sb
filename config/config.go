package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the configuration for the sb CLI tool.
type Config struct {
	Repo  string            `yaml:"repo"`  // Absolute path to the git repository
	Slugs map[string]string `yaml:"slugs"` // Map of slug to branch name
}

// LoadConfig loads the configuration from the given path, or the default location (~/.config/sb.yaml).
func LoadConfig() (*Config, error) {
	configPath := os.Getenv("SB_CONFIG")
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(home, ".config", "sb.yaml")
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
