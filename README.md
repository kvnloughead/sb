# sb CLI Tool

A simple Go command-line tool to switch to a git repository directory and optionally check out a branch using a slug.

## Features

- `sb` — switches to the repo directory.
- `sb [branchSlug]` — switches to the repo directory and checks out the corresponding branch.
- Configurable mapping of slug to branch and repo path.
- Config file location is user-defined, defaulting to `~/.config/sb.yaml`. To place the file somewhere else, set the `SB_CONFIG` environmental variable.

## Installation

### Option A: Installer (recommended)

- If you're in this repo, you can run it directly: `./sb install`
- If the binary is already on your PATH, run: `sb install`

The installer will:
- Create a starter config at `~/.config/sb.yaml` (or your chosen location)
- Copy the binary to an install directory (default: `~/.local/bin`)
- Optionally open the config in your editor

When finished, ensure the chosen install directory is on your PATH.

### Option B: Manual installation from source

Prerequisites: Go toolchain installed.

Using the Makefile:

```sh
make build      # builds local binary ./sb
make install    # installs to ~/.local/bin
make build-all  # cross-compile binaries into ./bin
```

Or directly via Go:

```sh
go build -o sb ./cmd/sb.go
install -d ~/.local/bin
install sb ~/.local/bin/sb
```

Then either run `sb install` to generate a starter config, or create the config file manually (see below).

## Usage

```sh
sb           # switches to the repo directory
sb dev       # switches to the repo directory and checks out the branch mapped to 'dev'
```

## Configuration

Create a config file (default: `~/.config/sb.yaml`) with the following structure, or run `sb install` to generate one:

Example:

```yaml
repo: /path/to/your/repo
slugs:
  dev: development
  main: main
  feat: feature-branch
```

- `repo`: Absolute path to your git repository.
- `slugs`: Map of slug to branch name.

You can specify a custom config file location with the `SB_CONFIG` environment variable.

## Shell integration and completions

To preserve shell history and enable tab-completions for branch slugs, use the scripts in [SCRIPTS.md](SCRIPTS.md).

Included:
- A wrapper function that preserves history when changing directories
- Bash completions that list your configured slugs

Add them to your shell config (e.g., `~/.bashrc`, `~/.bash_profile`, or `~/.zshrc`) and reload your shell.

## License
MIT

