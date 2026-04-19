# seerr-cli

A powerful, interactive command-line interface for Seerr and Overseerr. Manage your media requests, discover new content, and handle administrative tasks directly from your terminal.

## Features

- 🚀 **Interactive Mode**: A menu-driven interface for seamless browsing (`seerr interactive`).
- 🔍 **Media Discovery**: Search for movies and TV shows, or browse trending and popular content.
- 📬 **Request Management**: Create, list, and cancel requests.
- ℹ️ **Deep Insights**: View detailed metadata, ratings, and summaries for any media item.
- 🛡️ **Admin Suite**: Approve or deny pending requests (requires admin API key).
- ⚙️ **Configurable**: Easy setup for multiple instances and secure API key management.

## Installation

```bash
# Clone the repository
git clone https://github.com/LeeroyDing/seerr-cli.git

# Build the binary
cd seerr-cli
go build -o seerr main.go

# (Optional) Move to your PATH
mv seerr /usr/local/bin/
```

## Quick Start

### 1. Configure your instance
Set your Seerr instance URL and API Key:
```bash
seerr config --url https://seerr.example.com --api-key YOUR_API_KEY
```

### 2. Launch Interactive Mode
The easiest way to use the tool:
```bash
seerr interactive
```

### 3. Manual Commands
You can also use traditional CLI commands:
```bash
# Search for a movie
seerr search "Inception"

# View details
seerr info 27205 --type movie

# Request an item
seerr request 27205 --type movie

# List your requests
seerr list
```

## Roadmap Progress

| Version | Milestone | Status |
|:---|:---|:---:|
| v0.1.0 | Setup & Connectivity | ✅ |
| v0.2.0 | Media Discovery | ✅ |
| v0.3.0 | Media Requesting | ✅ |
| v0.4.0 | Request Management | ✅ |
| v0.5.0 | User Awareness | ✅ |
| v0.6.0 | Deep Dive | ✅ |
| v0.7.0 | Discovery Flows | ✅ |
| v0.8.0 | Quality Control | ✅ |
| v0.9.0 | Admin Control | ✅ |
| v0.10.0| Interactive Experience | ✅ |

## Development

Requires Go 1.21+.

```bash
go run main.go --help
```

## License

MIT
