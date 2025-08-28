# sb CLI Tool

A simple Go command-line tool to switch to a git repository directory and optionally check out a branch using a slug.

## Features
- `sb` — switches to the repo directory.
- `sb [branchSlug]` — switches to the repo directory and checks out the corresponding branch.
- Configurable mapping of slug to branch and repo path.
- Config file location is user-defined, defaulting to `~/.config/sb.yaml`. To place the file somewhere else, set the `SB_CONFIG` environmental variable.

## Usage

```sh
sb           # switches to the repo directory
sb dev       # switches to the repo directory and checks out the branch mapped to 'dev'
```

## Configuration is not necessary

A default configuration is used that assumes the canonicals repo is in `~/tripleten/se-canonicals_en`, and provides aliases for the commonly referenced branches. These default values can be found in [defaults.yaml](defaults.yaml).

## Configuration

Create a config file (default: `~/.config/sb.yaml`) with the following structure:

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

## Installation

```sh
# Recommended: installs binary to ~/.local/bin
make install

# Manual
go build -o sb ./cmd/sb.go && mv sb /directory/in/path
```

Make sure that the directory you install to is in your path.

## Preserving shell history

The default installation spawns a subshell which doesn't share the previous process's history. To preserve the history, add this script to `.bashrc`, `.bash_profile`, or an analogous file for your shell. Then source the file, with `source /path/to/file`. I've only tested this function on Bash, but it should work on zsh as well.

```bash
# sb shell function to preserve history after switching directory
sb() {
  local out
  out=$(SB_SHELL_WRAPPER=1 command sb "$@")
  local dir branch
  dir=$(echo "$out" | grep '^DIR:' | cut -d: -f2- | xargs)
  branch=$(echo "$out" | grep '^BRANCH:' | cut -d: -f2- | xargs)
  cd "$dir" || return
  if [ -n "$branch" ]; then
    git checkout "$branch"
  fi
}
```

## License
MIT

