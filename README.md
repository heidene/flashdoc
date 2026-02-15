# flashdoc

[![Development Status](https://img.shields.io/badge/status-alpha-orange?style=flat-square)](https://github.com/heidene/flashdoc/releases)
[![License](https://img.shields.io/badge/license-Beerware-yellow?style=flat-square)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/heidene/flashdoc?style=flat-square)](https://github.com/heidene/flashdoc/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/heidene/flashdoc?style=flat-square)](go.mod)
[![CI](https://img.shields.io/github/actions/workflow/status/heidene/flashdoc/ci.yml?branch=main&style=flat-square&label=tests)](https://github.com/heidene/flashdoc/actions)

> Ephemeral Starlight documentation sites from any folder of markdown files.

flashdoc is a CLI tool that transforms any directory of markdown files into a beautiful, searchable documentation site powered by [Astro Starlight](https://starlight.astro.build/). Think "man pages++" - instant documentation without the hassle.

## Features

- **Zero Configuration**: Just point to a folder and go
- **Automatic Setup**: Creates temporary workspace, installs dependencies, starts server
- **Smart Processing**: Auto-generates frontmatter from filenames
- **Clean UX**: Beautiful terminal output with real-time progress
- **Auto Cleanup**: Removes all temporary files on exit
- **Package Manager Detection**: Automatically uses pnpm, bun, or npm

## Quick Start

```bash
# View documentation from a folder
flashdoc ./docs

# With custom title
flashdoc ./api-docs --title "API Reference"

# Custom port
flashdoc ./guides --port 8080

# Don't open browser automatically
flashdoc ./docs --no-open
```

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/heidene/flashdoc/main/install.sh | sh
```

This will automatically download the latest release for your OS and architecture.

### Alternative Methods

**Using Go:**
```bash
go install github.com/heidene/flashdoc/cmd/flashdoc@latest
```

**Manual Download:**
Download a binary from the [releases page](https://github.com/heidene/flashdoc/releases).

## How It Works

1. **Scan**: Discovers all markdown files in your source directory
2. **Process**: Adds/fixes frontmatter for Starlight compatibility
3. **Setup**: Creates temporary workspace with embedded Starlight template
4. **Install**: Installs dependencies using your preferred package manager
5. **Serve**: Starts Astro dev server and opens in your browser
6. **Clean**: Removes everything on exit (Ctrl+C)

## Project Structure

The directory structure is organized by BDD phases:

```
flashdoc/
├── features/                    # Gherkin feature files (BDD specs)
│   ├── phase1-foundation/      # CLI, workspace, signals, cleanup
│   ├── phase2-markdown/        # Scanning, frontmatter, copying
│   ├── phase3-starlight/       # Template, config, dependencies
│   └── phase4-server/          # Dev server, browser, terminal UX
├── cmd/
│   └── flashdoc/
│       └── main.go             # CLI entry point
├── internal/
│   ├── cli/                    # Command-line parsing
│   ├── server/                 # Dev server management
│   ├── markdown/               # File processing
│   ├── template/               # Embedded template
│   └── pkgmgr/                 # Package manager detection
├── templates/
│   └── starlight/              # Embedded Starlight template
├── go.mod
├── Makefile
└── README.md
```

## Development Approach

This project follows **Behavior-Driven Development (BDD)** principles:

1. **Features First**: All functionality is defined in Gherkin `.feature` files
2. **Validation**: Feature files are reviewed before implementation
3. **Test-Driven**: Implementation is guided by feature scenarios
4. **Incremental**: Built phase by phase (Foundation → Markdown → Starlight → Server)

### Development Phases

| Phase | Focus | Features |
|-------|-------|----------|
| **Phase 1** | Foundation | CLI parsing, temp workspace, signal handling, cleanup |
| **Phase 2** | Markdown | Folder scanning, frontmatter injection, file copying |
| **Phase 3** | Starlight | Template extraction, config generation, dependencies |
| **Phase 4** | Server & UX | Dev server, browser open, terminal output |
| **Phase 5** | Advanced | AI generation, source watching, cache optimization |

## Building

```bash
# Install dependencies
go mod download

# Run tests
make test

# Build binary
make build

# Install locally
make install
```

## Usage Examples

### Basic Usage

```bash
# View markdown docs from current directory
flashdoc .

# View docs from specific folder
flashdoc ./documentation

# View nested structure
flashdoc ~/projects/api-docs
```

### With Flags

```bash
# Custom title
flashdoc ./docs --title "My Project Documentation"

# Custom port
flashdoc ./docs --port 3000

# Quiet mode (minimal output)
flashdoc ./docs --quiet

# Verbose mode (debug info)
flashdoc ./docs --verbose

# Don't auto-open browser
flashdoc ./docs --no-open

# Force specific package manager
flashdoc ./docs --package-manager pnpm
```

## CLI Reference

```
Usage: flashdoc <directory> [flags]

Flags:
  --title string             Custom site title (default: directory name)
  --port int                 Dev server port (default: 4321)
  --no-open                  Don't open browser automatically
  --quiet                    Minimal output
  --verbose                  Verbose output with debug info
  --package-manager string   Force package manager (pnpm, bun, npm)
  --silent                   Suppress package manager output
  --timestamps               Include timestamps in log output
  --help                     Show help
  --version                  Show version
```

## Requirements

- Go 1.21+ (for building)
- Node.js 18+ (for running Astro)
- One of: pnpm, bun, or npm

## Architecture Decisions

### Why Go?
- Single binary distribution
- Excellent process management
- `go:embed` for template bundling
- Fast startup time

### Why Starlight?
- Best-in-class documentation UX
- Auto-generates navigation from file structure
- Built-in search
- Mobile-friendly
- Minimal configuration required

### Why Temporary Workspace?
- No permanent files to manage
- Clean slate every run
- No git pollution
- Parallel instances don't conflict

## Feature Status

**Current Phase**: Feature validation and review

**Completed**:
- ✅ Project structure created
- ✅ Phase 1 features written (Foundation)
- ✅ Phase 2 features written (Markdown Processing)
- ✅ Phase 3 features written (Starlight Integration)
- ✅ Phase 4 features written (Server & UX)

**Next Steps**:
1. Review and validate all feature files
2. Set up BDD test framework (godog)
3. Implement Phase 1 (Foundation)
4. Implement Phase 2 (Markdown)
5. Implement Phase 3 (Starlight)
6. Implement Phase 4 (Server)

## Contributing

This project is in early development. Feature files define the expected behavior - implementation PRs should align with these specifications.

### Development Workflow

1. Review feature files in `features/` directory
2. Write tests using godog that match the scenarios
3. Implement code to make tests pass
4. Ensure all scenarios in the phase pass before moving to next phase

## License

[Beerware License](LICENSE) 🍺

Free to use, modify, and distribute. If we meet someday and you think this tool is worth it, you can buy me a beer!

## Inspiration

Inspired by tools like:
- `man` pages (instant reference)
- `python -m http.server` (zero-config serving)
- Starlight (beautiful documentation)

Combines the best of all worlds: instant setup, beautiful output, automatic cleanup.

---

**Status**: 🚧 In Development - Feature validation phase

**Version**: 0.1.0-alpha
