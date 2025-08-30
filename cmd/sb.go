package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sb/config"
	"strings"
)

// main is the entry point for the sb CLI tool.
// It loads configuration, switches to the repo directory, and optionally checks out a branch.

func main() {
	args := os.Args[1:]

	// Handle install command
	if len(args) > 0 && args[0] == "install" {
		runInstaller()
		return
	}

	// Handle completions command
	if len(args) > 0 && args[0] == "completions" {
		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		// Print all aliases for shell completion
		for alias := range cfg.Aliases {
			fmt.Println(alias)
		}
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\nRun 'sb install' to set up configuration.\n", err)
		os.Exit(1)
	}

	var branch string
	if len(args) > 0 {
		alias := args[0]
		var ok bool
		branch, ok = cfg.Aliases[alias]
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown alias: %s\n", alias)
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

// runInstaller prompts the user for configuration and sets up sb
func runInstaller() {
	fmt.Println("Welcome to sb installer!")
	fmt.Println()

	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	// Check if we're already installed (trying to reinstall)
	if strings.Contains(execPath, ".local/bin") || strings.Contains(execPath, "/usr/local/bin") || strings.Contains(execPath, "/usr/bin") {
		fmt.Println("It looks like sb is already installed.")
		fmt.Printf("Reinstall sb? [y/N]: ")
		response := readInput()
		if !(strings.ToLower(response) == "y" || strings.ToLower(response) == "yes") {
			fmt.Println("Installation cancelled.")
			return
		}
	}

	// Get current executable path
	execPath, err = os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting executable path: %v\n", err)
		os.Exit(1)
	}

	// Prompt for install directory
	fmt.Printf("Where would you like to install sb? [~/.local/bin]: ")
	installDir := readInput()
	if installDir == "" {
		home, _ := os.UserHomeDir()
		installDir = filepath.Join(home, ".local", "bin")
	} else if strings.HasPrefix(installDir, "~/") {
		home, _ := os.UserHomeDir()
		installDir = filepath.Join(home, installDir[2:])
	}

	// Prompt for repo path
	fmt.Printf("What is your repository path? [~/repo]: ")
	repoPath := readInput()
	if repoPath == "" {
		repoPath = "~/repo"
	}

	// Prompt for editor
	fmt.Printf("What is your preferred editor? [code]: ")
	editor := readInput()
	if editor == "" {
		editor = "code"
	}

	// Create config directory and file
	// If SB_CONFIG is set, honor it; otherwise default to ~/.config/sb.yaml
	configOutPath := os.Getenv("SB_CONFIG")
	if strings.TrimSpace(configOutPath) == "" {
		home2, _ := os.UserHomeDir()
		configOutPath = filepath.Join(home2, ".config", "sb.yaml")
	}
	configDir := filepath.Dir(configOutPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		os.Exit(1)
	}
	// Create config file with example (use aliases)
	configContent := fmt.Sprintf(`repo: %s
aliases:
	# Add your branch aliases here
	# Examples:
	# main: main
	# dev: development
	# feature: feature-branch
`, repoPath)
	if err := os.WriteFile(configOutPath, []byte(configContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Config created at: %s\n", configOutPath)

	// Install binary
	err = os.MkdirAll(installDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating install directory: %v\n", err)
		os.Exit(1)
	}

	destPath := filepath.Join(installDir, "sb")
	err = copyFile(execPath, destPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error installing binary: %v\n", err)
		os.Exit(1)
	}

	err = os.Chmod(destPath, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting binary permissions: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✓ Binary installed to: %s\n", destPath)
	fmt.Printf("\nOpen config file in %s? [Y/n]: ", editor)

	response := readInput()
	if response == "" || strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		cmd := exec.Command(editor, configOutPath)
		cmd.Run()
	} else {
		fmt.Printf("\nNext steps:\n")
		fmt.Printf("1. Edit %s to add your branch mappings\n", configOutPath)
		fmt.Printf("2. Add the shell function to your .bashrc or .zshrc for seamless integration\n")
		fmt.Printf("3. Make sure %s is in your PATH\n", installDir)
	}
}

// readInput reads a line from stdin
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}
