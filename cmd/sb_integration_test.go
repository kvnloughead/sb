package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// runSB runs `go run ./cmd/sb.go` with args and returns stdout, stderr, and error
func runSB(t *testing.T, env []string, args ...string) (string, string, error) {
	t.Helper()
	// Determine the cmd directory of this package
	_, thisFile, _, _ := runtime.Caller(0)
	cmdDir := filepath.Dir(thisFile)
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Dir = cmdDir
	// Preserve base env and add overrides
	cmd.Env = append(os.Environ(), env...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// runSBWithInput is like runSB but feeds input to stdin
func runSBWithInput(t *testing.T, input string, env []string, args ...string) (string, string, error) {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	cmdDir := filepath.Dir(thisFile)
	cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
	cmd.Dir = cmdDir
	cmd.Env = append(os.Environ(), env...)
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func TestCompletionsOutputsAliases(t *testing.T) {
	// Ensure repo dir exists for validation
	_ = os.MkdirAll("/tmp/repo", 0o755)
	cfg := "repo: /tmp/repo\naliases:\n  p1: branch1\n  p2: branch2\n"
	// Supply config via SB_CONFIG env
	tmp := t.TempDir()
	cfgPath := tmp + "/sb.yaml"
	if err := os.WriteFile(cfgPath, []byte(cfg), 0644); err != nil {
		t.Fatalf("write cfg: %v", err)
	}
	out, stderr, err := runSB(t, []string{"SB_CONFIG=" + cfgPath}, "completions")
	if err != nil {
		t.Fatalf("completions err: %v, stderr=%s", err, stderr)
	}
	if !strings.Contains(out, "p1") || !strings.Contains(out, "p2") {
		t.Fatalf("expected aliases in output, got: %q", out)
	}
}

func TestShellWrapperProtocol(t *testing.T) {
	// Ensure repo dir exists for validation
	_ = os.MkdirAll("/tmp/repo", 0o755)
	cfg := "repo: /tmp/repo\naliases:\n  d: dev\n"
	tmp := t.TempDir()
	cfgPath := tmp + "/sb.yaml"
	if err := os.WriteFile(cfgPath, []byte(cfg), 0644); err != nil {
		t.Fatalf("write cfg: %v", err)
	}
	// No alias
	out, stderr, err := runSB(t, []string{"SB_CONFIG=" + cfgPath, "SB_SHELL_WRAPPER=1"})
	if err != nil {
		t.Fatalf("wrapper err: %v, stderr=%s", err, stderr)
	}
	if !strings.Contains(out, "DIR:/tmp/repo") {
		t.Fatalf("expected DIR line, got: %q", out)
	}
	// With alias
	out, stderr, err = runSB(t, []string{"SB_CONFIG=" + cfgPath, "SB_SHELL_WRAPPER=1"}, "d")
	if err != nil {
		t.Fatalf("wrapper with alias err: %v, stderr=%s", err, stderr)
	}
	if !strings.Contains(out, "BRANCH:dev") {
		t.Fatalf("expected BRANCH line, got: %q", out)
	}
}

func TestInstallerCreatesFiles(t *testing.T) {
	// Use explicit paths under a temp dir without changing HOME
	tmpRoot := t.TempDir()
	installDir := filepath.Join(tmpRoot, "bin")
	cfgPath := filepath.Join(tmpRoot, "sb.yaml")

	// Provide inputs for prompts:
	// 1) install dir -> installDir
	// 2) repo path   -> /tmp/repo
	// 3) editor      -> code
	// 4) open config -> n
	input := strings.Join([]string{
		installDir,
		"/tmp/repo",
		"code",
		"n",
	}, "\n") + "\n"

	out, stderr, err := runSBWithInput(t, input, []string{"SB_CONFIG=" + cfgPath}, "install")
	if err != nil {
		t.Fatalf("installer err: %v, stderr=%s, stdout=%s", err, stderr, out)
	}

	// Verify binary installed
	if fi, err := os.Stat(filepath.Join(installDir, "sb")); err != nil || fi.IsDir() {
		t.Fatalf("expected installed binary at %s", filepath.Join(installDir, "sb"))
	}

	// Verify config created
	if fi, err := os.Stat(cfgPath); err != nil || fi.IsDir() {
		t.Fatalf("expected config file at %s", cfgPath)
	}
}
