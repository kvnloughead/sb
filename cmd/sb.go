package main

import (
	"bufio"
	"flag"
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
		runInstaller(args[1:])
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

	// Validate that the repo directory exists before proceeding
	if fi, statErr := os.Stat(absDir); statErr != nil || !fi.IsDir() {
		if os.IsNotExist(statErr) {
			fmt.Fprintf(os.Stderr, "Repository directory does not exist: %s\n", absDir)
		} else if statErr != nil {
			fmt.Fprintf(os.Stderr, "Error accessing repository directory %s: %v\n", absDir, statErr)
		} else {
			fmt.Fprintf(os.Stderr, "Repository path is not a directory: %s\n", absDir)
		}
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
func runInstaller(argv []string) {
	fmt.Println("Welcome to sb installer!")
	fmt.Println()

	// Flags for non-interactive / scripted installs
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	var (
		flagDir    string
		flagRepo   string
		flagEditor string
		flagYes    bool
		flagNoOpen bool
	)
	fs.StringVar(&flagDir, "dir", "", "install directory (default ~/.local/bin)")
	fs.StringVar(&flagRepo, "repo", "", "repository path (default ~/repo)")
	fs.StringVar(&flagEditor, "editor", "", "editor command to open config (default 'code')")
	fs.BoolVar(&flagYes, "yes", false, "accept defaults; do not prompt")
	fs.BoolVar(&flagYes, "y", false, "shorthand for --yes")
	fs.BoolVar(&flagNoOpen, "no-open", false, "do not prompt to open the config file")
	_ = fs.Parse(argv)

	// If launched via a shell wrapper and we're about to prompt, exit with guidance to avoid hangs
	if os.Getenv("SB_SHELL_WRAPPER") == "1" && !flagYes && flagDir == "" && flagRepo == "" && flagEditor == "" {
		fmt.Fprintln(os.Stderr, "Detected shell wrapper. Run 'command sb install' or pass --yes/--dir/--repo/--editor flags to avoid prompts.")
		os.Exit(2)
	}

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

	// Re-resolve executable path (absolute + symlinks)
	if abs, err := filepath.Abs(execPath); err == nil {
		execPath = abs
	}
	if link, err := filepath.EvalSymlinks(execPath); err == nil && link != "" {
		execPath = link
	}

	// Install directory
	var installDir string
	if flagDir != "" {
		installDir = flagDir
	} else if flagYes {
		home, _ := os.UserHomeDir()
		installDir = filepath.Join(home, ".local", "bin")
	} else {
		fmt.Printf("Where would you like to install sb? [~/.local/bin]: ")
		installDir = readInput()
	}
	if installDir == "" {
		home, _ := os.UserHomeDir()
		installDir = filepath.Join(home, ".local", "bin")
	} else if strings.HasPrefix(installDir, "~/") {
		home, _ := os.UserHomeDir()
		installDir = filepath.Join(home, installDir[2:])
	}

	// Repo path
	var repoPath string
	if flagRepo != "" {
		repoPath = flagRepo
	} else if flagYes {
		repoPath = "~/repo"
	} else {
		fmt.Printf("What is your repository path? [~/repo]: ")
		repoPath = readInput()
	}
	if repoPath == "" {
		repoPath = "~/repo"
	}

	// Editor
	var editor string
	if flagEditor != "" {
		editor = flagEditor
	} else if flagYes {
		editor = "code"
	} else {
		fmt.Printf("What is your preferred editor? [code]: ")
		editor = readInput()
	}
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
	// Create config file with example (use aliases) if it doesn't already exist
	if _, err := os.Stat(configOutPath); err == nil {
		fmt.Printf("✓ Config already exists: %s (leaving it unchanged)\n", configOutPath)
	} else {
		configContent := fmt.Sprintf(`repo: %s
aliases:
  # Add your branch aliases here. Use spaces for indenting.
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
	}

	// Install binary
	err = os.MkdirAll(installDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating install directory: %v\n", err)
		os.Exit(1)
	}

	destPath := filepath.Join(installDir, "sb")
	// Resolve destination for robust comparison
	destResolved := destPath
	if abs, err := filepath.Abs(destResolved); err == nil {
		destResolved = abs
	}
	if link, err := filepath.EvalSymlinks(destResolved); err == nil && link != "" {
		destResolved = link
	}

	same := false
	if si, err1 := os.Stat(execPath); err1 == nil {
		if di, err2 := os.Stat(destResolved); err2 == nil {
			same = os.SameFile(si, di)
		}
	}

	if same {
		// Avoid copying a file onto itself (would truncate). Just ensure perms.
		if err := os.Chmod(destResolved, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting binary permissions: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := copyFileAtomic(execPath, destResolved); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing binary: %v\n", err)
			os.Exit(1)
		}
		if err := os.Chmod(destResolved, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error setting binary permissions: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\n✓ Binary installed to: %s\n", destResolved)
	if !flagNoOpen && !flagYes {
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

// copyFileAtomic copies src to dst via a temporary file and atomic rename.
func copyFileAtomic(src, dst string) error {
	dir := filepath.Dir(dst)
	tmp, err := os.CreateTemp(dir, ".sb-tmp-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	// Ensure cleanup on failure
	defer func() {
		tmp.Close()
		os.Remove(tmpPath)
	}()

	// Copy contents
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	if _, err := tmp.ReadFrom(srcF); err != nil {
		return err
	}
	if err := tmp.Chmod(0755); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}

	// Atomic replace
	return os.Rename(tmpPath, dst)
}
