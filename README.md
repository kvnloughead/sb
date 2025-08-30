# sb CLI Tool

A simple Go command-line tool to switch to a git repository directory and optionally check out a branch using a slug.

## Features

- `sb` — switches to the repo directory.
- `sb [branchSlug]` — switches to the repo directory and checks out the corresponding branch.
- Configurable mapping of slug to branch and repo path.
- Config file location is user-defined, defaulting to `~/.config/sb.yaml`. To place the file somewhere else, set the `SB_CONFIG` environmental variable.

## Installation

- Download the latest release for your OS
- Move it to `~/.local/bin` or another directory in `$PATH`.
- Optionally, run `sb install` to set up a starter config and confirm PATH setup
- Optionally, add the shell function below to prevent shell history disruption when changing directories

## Usage

```sh
sb           # switches to the repo directory
sb dev       # switches to the repo directory and checks out the branch mapped to 'dev'
```

## Configuration

Create a config file (default: `~/.config/sb.yaml`) with the following structure or run `sb install` to generate one:

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

## Building from source

```sh
make build      # build local binary ./sb
make install    # install to ~/.local/bin
make build-all  # cross-compile binaries into ./bin
```

Ensure the install directory is on your PATH.

## Preserving shell history

The default installation spawns a subshell which doesn't share the previous process's history. To preserve the history, add this script to `.bashrc`, `.bash_profile`, or an analogous file for your shell. Then source the file, with `source /path/to/file`. I've only tested this function on Bash, but it should work on zsh as well.

```bash
# sb shell function to preserve history after switching directory
sb() {
  # Bypass wrapper for commands that must print directly
  case "$1" in
    install|-h|--help|completions)
      command sb "$@"
      return
      ;;
  esac
  local out dir branch
  out=$(SB_SHELL_WRAPPER=1 command sb "$@")
  dir=$(echo "$out" | grep '^DIR:' | cut -d: -f2- | xargs)
  branch=$(echo "$out" | grep '^BRANCH:' | cut -d: -f2- | xargs)
  [ -n "$dir" ] || return
  cd "$dir" || return
  if [ -n "$branch" ]; then
    git checkout "$branch"
  fi
}
```

## License
MIT

