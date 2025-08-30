package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "sb.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return p
}

func TestLoadConfigFromEnv(t *testing.T) {
	cfgContent := "" +
		"repo: /tmp/repo\n" +
		"slugs:\n" +
		"  dev: development\n" +
		"  main: main\n"
	p := writeTempConfig(t, cfgContent)

	t.Setenv("SB_CONFIG", p)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.Repo != "/tmp/repo" {
		t.Fatalf("expected repo=/tmp/repo, got %q", cfg.Repo)
	}
	if cfg.Slugs["dev"] != "development" || cfg.Slugs["main"] != "main" {
		t.Fatalf("unexpected slugs: %#v", cfg.Slugs)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	t.Setenv("SB_CONFIG", filepath.Join(t.TempDir(), "nope.yaml"))
	if _, err := LoadConfig(); err == nil {
		t.Fatalf("expected error for missing config file")
	}
}

func TestLoadConfigMissingRepo(t *testing.T) {
	cfgContent := "slugs:\n  a: b\n"
	p := writeTempConfig(t, cfgContent)
	t.Setenv("SB_CONFIG", p)
	if _, err := LoadConfig(); err == nil {
		t.Fatalf("expected error for missing repo field")
	}
}
