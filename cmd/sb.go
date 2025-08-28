package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sb/config"
)

// main is the entry point for the sb CLI tool.
// It loads configuration, switches to the repo directory, and optionally checks out a branch.
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		// No slug provided, just switch to repo directory
		changeDirAndBranch(cfg.Repo, "")
		return
	}

	slug := args[0]
	branch, ok := cfg.Slugs[slug]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown branch slug: %s\n", slug)
		os.Exit(1)
	}
	changeDirAndBranch(cfg.Repo, branch)
}

// changeDir switches to the given directory and opens a new shell.
// dir: absolute or relative path to the directory.
// changeDirAndBranch switches to the given directory and optionally checks out a branch, then opens a new shell.
// dir: absolute or relative path to the directory.
// branch: name of the branch to checkout (empty for none).
func changeDirAndBranch(dir string, branch string) {
	// Expand ~ to home directory if present
	if len(dir) > 1 && dir[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err == nil {
			dir = filepath.Join(home, dir[2:])
		}
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid repo path: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Switching to repo directory: %s\n", absDir)
	var cmdStr string
	if branch != "" {
		fmt.Printf("Checking out branch: %s\n", branch)
		cmdStr = fmt.Sprintf("cd '%s' && git checkout '%s'; bash", absDir, branch)
	} else {
		cmdStr = fmt.Sprintf("cd '%s'; bash", absDir)
	}
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
