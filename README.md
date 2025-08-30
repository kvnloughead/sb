# sb CLI Tool

A simple Go command-line tool to switch to a git repository directory and optionally check out a branch using an alias.

## Features

- `sb` — switches to the repo directory.
- `sb [alias]` — switches to the repo directory and checks out the corresponding branch.
- Configurable mapping of alias to branch and repo path.
- Config file location is user-defined, defaulting to `~/.config/sb.yaml`. To place the file somewhere else, set the `SB_CONFIG` environmental variable.

## Installation

### Option A: Installer (recommended)

- If you're in this repo, run `make build` followed by `./sb install`
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
sb dev       # switches to the repo directory and checks out the branch mapped to alias 'dev'
```

## Configuration

Create a config file (default: `~/.config/sb.yaml`) with the following structure, or run `sb install` to generate one:

Example (aliases):

```yaml
repo: /path/to/your/repo
aliases:
  dev: development
  main: main
  feat: feature-branch
```

- `repo`: Absolute path to your git repository.
- `aliases`: Map of alias to branch name.

You can specify a custom config file location with the `SB_CONFIG` environment variable.

## Shell integration and completions

To preserve shell history and enable tab-completions for branch aliases, use the scripts in [SCRIPTS.md](SCRIPTS.md).

Included:
- A wrapper function that preserves history when changing directories
- Bash completions that list your configured aliases

Add them to your shell config (e.g., `~/.bashrc`, `~/.bash_profile`, or `~/.zshrc`) and reload your shell.

## Local development: block pushes if tests fail

This repo includes a local git pre-push hook that runs tests and prevents a push if they fail. Enable it once per clone:

```sh
make hooks
```

After that, every `git push` will run `go test ./...` first and abort on failures.

## Troubleshooting installation/run

- If `sb` appears to do nothing but `./sb install` works, your shell function wrapper may be swallowing errors. Bypass the wrapper with:
  - `command sb install`
  - or temporarily rename/remove the wrapper function in your shell config.
- Prefer installing from a release binary: download for your OS, then run `./sb-* install` (or `sb install` if it’s already on your PATH). The installer will copy the binary to `~/.local/bin` by default.

## License
MIT

