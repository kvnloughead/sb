# sb CLI Tool

A simple Go command-line tool to switch to a git repository directory and optionally check out a branch using a slug.

## Features
- `sb` — switches to the repo directory.
- `sb [branchSlug]` — switches to the repo directory and checks out the corresponding branch.
- Configurable mapping of slug to branch and repo path.
- Config file location is user-defined, defaulting to `~/.config/sb.yaml`.

## Usage

```sh
sb           # switches to the repo directory
sb dev       # switches to the repo directory and checks out the branch mapped to 'dev'
```

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
go build -o sb ./cmd/sb.go && mv sb ~/.local/bin
```

Make sure that the directory you move it to is in your PATH.

## License
MIT

