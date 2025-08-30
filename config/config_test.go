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

func TestLoadConfigFromEnv_Aliases(t *testing.T) {
	cfgContent := "" +
		"repo: /tmp/repo\n" +
		"aliases:\n" +
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
	if cfg.Aliases["dev"] != "development" || cfg.Aliases["main"] != "main" {
		t.Fatalf("unexpected aliases: %#v", cfg.Aliases)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	t.Setenv("SB_CONFIG", filepath.Join(t.TempDir(), "nope.yaml"))
	if _, err := LoadConfig(); err == nil {
		t.Fatalf("expected error for missing config file")
	}
}

func TestLoadConfigMissingRepo(t *testing.T) {
	cfgContent := "aliases:\n  a: b\n"
	p := writeTempConfig(t, cfgContent)
	t.Setenv("SB_CONFIG", p)
	if _, err := LoadConfig(); err == nil {
		t.Fatalf("expected error for missing repo field")
	}
}

func TestLoadConfig_ContainsAliases(t *testing.T) {
	cfgContent := "repo: /tmp/repo\naliases:\n  dev: development\n"
	p := writeTempConfig(t, cfgContent)
	t.Setenv("SB_CONFIG", p)
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig returned error: %v", err)
	}
	if cfg.Aliases["dev"] != "development" {
		t.Fatalf("expected alias mapping, got: %#v", cfg.Aliases)
	}
}
