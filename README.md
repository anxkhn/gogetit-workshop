# gogetit-workshop

An educational repository with seeded bugs for learning open source contributions.

## Project Description

gogetit is a concurrent CLI file downloader and website metadata scraper written in Go. This workshop repository is designed to help developers practice finding, fixing, and preventing bugs while learning the open source contribution workflow.

## Features

- **Download**: Concurrently download multiple files from URLs with progress tracking
- **Scrape**: Extract metadata from websites including title, description, and Open Graph tags
- **TUI Progress**: Terminal user interface with real-time progress visualization using Bubble Tea
- **Config File**: YAML-based configuration for customizing download settings and defaults

## Learning Objectives

By working through the issues in this repository, you will learn:

- Go concurrency patterns (goroutines, channels, sync primitives)
- CLI development with Cobra framework
- Terminal User Interface patterns with Bubble Tea
- Error handling and edge cases in concurrent code
- Testing concurrent code effectively
- Open source contribution workflow and best practices

## Quick Start

```bash
# Clone the repository
git clone https://github.com/anxkhn/gogetit-workshop.git
cd gogetit-workshop

# Install dependencies
go mod download

# Build the binary
go build -o gogetit ./cmd/gogetit

# Download files concurrrently
./gogetit download https://example.com/file1.zip https://example.com/file2.pdf

# Scrape website metadata
./gogetit scrape https://example.com

# View version
./gogetit version
```

## Installation

### Prerequisites

- Go 1.24 or higher
- golangci-lint (for linting)

### Dependencies

- github.com/spf13/cobra v1.10.2 - CLI framework
- github.com/charmbracelet/bubbletea v1.3.10 - TUI framework
- github.com/charmbracelet/bubbles v0.20.0 - UI components
- github.com/charmbracelet/lipgloss v0.13.0 - Styling
- gopkg.in/yaml.v3 - YAML configuration parsing

### Build from Source

```bash
# Clone the repository
git clone https://github.com/anxkhn/gogetit-workshop.git
cd gogetit-workshop

# Build
go build -o gogetit ./cmd/gogetit

# Run tests
go test ./...

# Install locally
go install ./cmd/gogetit
```

## Project Structure

```
gogetit-workshop/
├── cmd/
│   └── gogetit/
│       └── main.go              # Application entry point
├── internal/
│   ├── cmd/
│   │   ├── root.go              # Root command
│   │   ├── download.go          # Download command
│   │   ├── scrape.go            # Scrape command
│   │   └── version.go           # Version command
│   ├── downloader/
│   │   ├── downloader.go        # Download logic
│   │   ├── pool.go              # Worker pool
│   │   ├── worker.go            # Worker implementation
│   │   └── config.go            # Downloader config
│   ├── scraper/
│   │   ├── scraper.go           # Scraping logic
│   │   ├── metadata.go          # Metadata types
│   │   └── parser.go            # HTML parser
│   ├── progress/
│   │   ├── progress.go          # Progress tracking
│   │   └── model.go             # TUI model
│   ├── config/
│   │   └── config.go            # Configuration handling
│   └── version/
│       └── version.go           # Version info
├── pkg/
│   └── utils/
│       ├── http.go              # HTTP utilities
│       ├── file.go              # File utilities
│       └── url.go               # URL utilities
├── test/
│   └── integration/
│       └── download_test.go     # Integration tests
├── .github/
│   ├── workflows/
│   │   ├── commit-check.yml     # Commit message validation
│   │   └── pr-check.yml         # PR description validation
│   ├── ISSUE_TEMPLATE/
│   │   ├── bug_report.yml
│   │   ├── feature_request.yml
│   │   └── config.yml
│   └── PULL_REQUEST_TEMPLATE.md
├── go.mod
├── go.sum
├── Makefile
├── .golangci.yml
├── .gitignore
├── LICENSE
├── README.md
└── CONTRIBUTING.md
```

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/downloader/...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race ./...
```

## Issue Labels Guide

Issues in this repository are categorized by difficulty:

| Label | Description | Skills Required |
|-------|-------------|-----------------|
| `good-first-issue` | Great for newcomers (15-30 min) | Basic Go knowledge, documentation |
| `intermediate` | Moderate complexity (1-2 hours) | Concurrency, error handling, testing |
| `advanced` | Complex problems (3-6 hours) | Deep understanding of concurrent patterns, TUI |

Look for these labels to find issues matching your skill level.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on:

- Setting up your development environment
- Contribution workflow
- Code style requirements
- Pull request process

## Using github.dev

This repository supports github.dev for browser-based editing:

1. Press `.` (period) while viewing any file in GitHub
2. Or replace `github.com` with `github.dev` in the URL
3. Make your changes in the web-based VS Code editor
4. Commit and create a PR directly from the browser

This is useful for small fixes and documentation improvements without setting up a local development environment.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Quick Reference

### Download Command
```bash
./gogetit download <url> [urls...]
```

Flags:
- `-o, --output-dir` - Output directory (default: current)
- `-c, --concurrency` - Number of concurrent downloads (default: 3)
- `-t, --timeout` - Timeout in seconds (default: 30)
- `-r, --retry` - Number of retry attempts (default: 3)

### Scrape Command
```bash
./gogetit scrape <url> [urls...]
```

### Examples
```bash
# Download a single file
./gogetit download https://example.com/file.zip

# Download multiple files with custom concurrency
./gogetit download -c 5 https://example.com/file1.zip https://example.com/file2.zip

# Scrape metadata from a website
./gogetit scrape https://example.com
```
