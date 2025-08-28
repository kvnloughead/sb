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
	var branch string
	if len(args) > 0 {
		slug := args[0]
		var ok bool
		branch, ok = cfg.Slugs[slug]
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown branch slug: %s\n", slug)
			os.Exit(1)
		}
	}

	// Expand ~ to home directory if present
	dir := cfg.Repo
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

	// If running in a shell function, print info for wrapper
	if os.Getenv("SB_SHELL_WRAPPER") == "1" {
		fmt.Printf("DIR:%s\n", absDir)
		if branch != "" {
			fmt.Printf("BRANCH:%s\n", branch)
		}
		return
	}

	// Otherwise, launch a subshell in the repo directory
	var cmdStr string
	if branch != "" {
		fmt.Printf("Switching to repo directory: %s\n", absDir)
		fmt.Printf("Checking out branch: %s\n", branch)
		cmdStr = fmt.Sprintf("cd '%s' && git checkout '%s'; bash", absDir, branch)
	} else {
		fmt.Printf("Switching to repo directory: %s\n", absDir)
		cmdStr = fmt.Sprintf("cd '%s'; bash", absDir)
	}
	cmd := exec.Command("bash", "-c", cmdStr)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
