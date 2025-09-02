# sb CLI Tool

A simple Go command-line tool for easy switching between git branches. It lets you set aliases for branches and switch to them with `sb <alias>`. Its use case is a bit niche. But if you have a repo with a lot of long-lasting branches in it that you need to periodically refer to, this tool can help you avoid the pain of trying to remember the names of the branches. There is an optional script to enable tab completions.

Not tested on Windows yet.

## Features

- `sb` — switches to the repo directory.
- `sb [alias]` — switches to the repo directory and checks out the corresponding branch.
- Configurable mapping of alias to branch and repo path.
- Config file location is user-defined, defaulting to `~/.config/sb.yaml`. To place the file somewhere else, set the `SB_CONFIG` environmental variable.

## Usage

It's very simple.

```sh
sb           # switches to the repo directory
sb dev       # switches to the repo directory and checks out the branch mapped to alias 'dev'
```

## Installation

It is recommended to install from the latest release, but instructions for installing from source are at the bottom of this page.

- Download the [latest release](https://github.com/kvnloughead/sb/releases/) 
- Run: `sb-* install`

The installer will:

- Create a starter config at `~/.config/sb.yaml` (or your chosen location)
- Copy the binary to an install directory (default: `~/.local/bin`)
- Optionally open the config in your editor

When finished, make sure the chosen install directory is on your PATH.

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

## License
MIT

## Manual installation from source

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

## Local development: block pushes if tests fail

This repo includes a local git pre-push hook that runs tests and prevents a push if they fail. Enable it once per clone:

```sh
make hooks
```

After that, every `git push` will run `go test ./...` first and abort on failures.

