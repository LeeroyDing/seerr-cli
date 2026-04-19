---
name: seerr-cli
description: "Manage media requests, discover content, and handle admin tasks for Seerr/Overseerr using the CLI. WHEN: \"seerr search\", \"add imdb\", \"request media\", \"approve requests\", \"interactive mode\"."
license: MIT
metadata:
  author: Leeroy Ding
  version: "1.0.0"
---

# Seerr CLI Skill

This skill provides instructions for interacting with the `seerr-cli` tool to manage media requests and discovery on Seerr/Overseerr instances.

## When to Use

- Configuring connection to a Seerr/Overseerr instance
- Searching for movies or TV shows
- Adding media via IMDb links or IDs
- Requesting media for download
- Managing requests (list, cancel, approve, deny)
- Browsing trending or popular content

## Core Commands

### Configuration
Before use, configure the instance URL and API key.
```bash
go run main.go config --url https://seerr.example.com --api-key YOUR_API_KEY
```

### Discovery
- **Search**: `go run main.go search "Inception"`
- **Browse**: `go run main.go browse [trending|movies|tv]`
- **Info**: `go run main.go info <id> --type [movie|tv]`

### Requests
- **Request**: `go run main.go request <id> --type [movie|tv]`
- **List**: `go run main.go list`
- **Cancel**: `go run main.go cancel <id>`
- **Add (IMDb)**: `go run main.go add <imdb_url_or_id>`

### Admin Operations
Requires an admin API key.
- **List Pending**: `go run main.go admin list-pending`
- **Approve**: `go run main.go admin approve <request_id>`
- **Deny**: `go run main.go admin deny <request_id>`

### Interactive Mode
Launch a menu-driven interface for easier navigation.
```bash
go run main.go interactive
```

## Best Practices

- **Use Interactive Mode**: For browsing and discovery, `interactive` mode is often more efficient.
- **IMDb Links**: Use `seerr add` with a direct IMDb link for the fastest way to request a specific item.
- **Type Flag**: Always specify `--type movie` or `--type tv` when using `info` or `request` if prompted.
- **Admin Key**: Ensure you are using an admin API key for `admin` subcommands.
- **Command Help**: Use `seerr [command] --help` to see detailed instructions and examples for any command. The help messages are optimized for both humans and AI agents.
