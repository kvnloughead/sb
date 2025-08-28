package config

import (
	"os"
	"testing"
)

func TestDefaultConfigGenerated(t *testing.T) {
	if DefaultConfig.Repo == "" {
		t.Errorf("DefaultConfig.Repo should not be empty")
	}
	if len(DefaultConfig.Slugs) == 0 {
		t.Errorf("DefaultConfig.Slugs should not be empty")
	}
}

func TestLoadConfigUsesDefaultsIfMissing(t *testing.T) {
	// Temporarily unset SB_CONFIG and rename user config if present
	origEnv := os.Getenv("SB_CONFIG")
	os.Setenv("SB_CONFIG", "/tmp/nonexistent_sb.yaml")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should not error if config missing: %v", err)
	}
	if cfg.Repo != DefaultConfig.Repo {
		t.Errorf("Expected default repo, got %s", cfg.Repo)
	}
	if len(cfg.Slugs) != len(DefaultConfig.Slugs) {
		t.Errorf("Expected default slugs, got %v", cfg.Slugs)
	}

	os.Setenv("SB_CONFIG", origEnv)
}
