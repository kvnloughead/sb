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
		changeDir(cfg.Repo)
		return
	}

	slug := args[0]
	branch, ok := cfg.Slugs[slug]
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown branch slug: %s\n", slug)
		os.Exit(1)
	}
	changeDir(cfg.Repo)
	checkoutBranch(branch)
}

// changeDir switches to the given directory and opens a new shell.
// dir: absolute or relative path to the directory.
func changeDir(dir string) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid repo path: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Switching to repo directory: %s\n", absDir)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("cd '%s'; exec $SHELL", absDir))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// checkoutBranch checks out the given branch in the current repo.
// branch: name of the branch to checkout.
func checkoutBranch(branch string) {
	fmt.Printf("Checking out branch: %s\n", branch)
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
