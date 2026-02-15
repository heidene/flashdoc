# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**flashdoc** is a CLI tool that creates ephemeral Starlight documentation sites from any directory of markdown files. It's built with Go and follows **Behavior-Driven Development (BDD)** principles using Godog/Cucumber.

**Core concept**: Point flashdoc at a markdown folder, and it instantly creates a beautiful searchable documentation site in a temporary workspace, then cleans everything up on exit. Zero configuration, zero permanent files.

## BDD-First Development Philosophy

This project is **BDD-driven from the ground up**. This is not optional—it's the core development methodology:

### Feature Files Are The Source of Truth

- **All functionality** is defined in Gherkin `.feature` files in `features/` directory
- **14 feature files** organized into 4 development phases (foundation, markdown, starlight, server)
- **213 total scenarios**: 89 implemented and passing, 124 undefined (future work)
- Feature files must be read and understood before implementing any code
- Implementation must match the exact behavior described in scenarios

### Development Workflow

1. **Read the feature file** to understand expected behavior
2. **Check/write step definitions** in `features/steps/*.go` (16 step files)
3. **Implement code** in `internal/` packages to make tests pass
4. **Run tests** with `go test -v ./features/...` (NOT `make features`)
5. **Iterate** until all scenarios pass

### Running BDD Tests

```bash
# Run all BDD tests (current: 89 passing, 124 undefined)
go test -v ./features/...

# Run specific phase
go test -v ./features/... -run TestPhase1

# Run specific scenario
go test -v ./features/... -run TestPhase2/Support_various_markdown_extensions

# Count test status
go test -v ./features/... 2>&1 | grep "scenarios"
```

**IMPORTANT**:
- Use `go test -v ./features/...` to run tests (NOT `make features`)
- The Makefile `features` target may be outdated
- Step definitions are in `features/steps/`, organized by domain (cli_steps.go, scanner_steps.go, etc.)

### Phase-Based Organization

Features are organized into phases that must be completed sequentially:

- **Phase 1 (Foundation)**: CLI parsing, temp workspace, signal handling, cleanup - ✅ COMPLETE
- **Phase 2 (Markdown)**: File scanning, frontmatter injection, copying - ✅ COMPLETE
- **Phase 3 (Starlight)**: Template extraction, config generation, dependencies - ✅ COMPLETE
- **Phase 4 (Server & UX)**: Dev server, browser launch, terminal output - ✅ COMPLETE

Current status: **All implemented features passing (89/89)**. Remaining 124 scenarios are for advanced features (AI generation, file watching, cache optimization).

## Architecture

### High-Level Flow

```
CLI Input → Workspace Creation → Markdown Processing → Template Setup →
Dependency Install → Dev Server → Browser Launch → Cleanup on Exit
```

### Key Design Patterns

1. **Embedded Template**: Starlight template is embedded in binary via `go:embed` in `internal/template/starlight/`
2. **Temporary Workspace**: All builds happen in OS temp directory (`os.MkdirTemp`), cleaned up via `defer`
3. **Signal Handling**: SIGINT/SIGTERM captured to ensure clean shutdown
4. **Process Management**: Dev server runs as child process, terminated on exit
5. **Package Manager Detection**: Auto-detects pnpm → bun → npm in priority order

### Internal Package Structure

**Critical packages** (understand these first):

- `internal/workspace/`: Creates temp directory with Starlight structure
- `internal/cli/`: Cobra-based CLI with flags (--title, --port, --no-open, etc.)
- `internal/processor/`: Orchestrates markdown scanning and copying
- `internal/scanner/`: Walks directories, finds markdown files (.md, .markdown, .mdown, .mkd)
- `internal/frontmatter/`: Parses/injects YAML frontmatter for Starlight compatibility
- `internal/template/`: Extracts embedded Starlight template to workspace
- `internal/server/`: Manages Astro dev server lifecycle
- `internal/cleanup/`: Handles cleanup on exit (workspace + child processes)
- `internal/signal/`: Intercepts OS signals for graceful shutdown

**Supporting packages**:

- `internal/pkgmanager/`: Detects available package manager (pnpm/bun/npm)
- `internal/installer/`: Runs npm/pnpm/bun install with progress
- `internal/browser/`: Opens default browser to dev server URL

### Critical Implementation Details

#### Frontmatter Injection

The `frontmatter.Inject()` function:
- Parses existing YAML frontmatter (if any)
- Adds `title` field if missing (derived from filename)
- Handles special cases: README.md → "Home", index.md → parent dir name
- Preserves other frontmatter fields via inline YAML marshaling
- **Recent fix**: Upgraded Starlight from v0.29.3 to v0.37.3 to fix TypeError bug

#### Markdown File Processing

The scanner supports multiple extensions: `.md`, `.markdown`, `.mdown`, `.mkd`

Files are processed with directory structure preserved:
```
source/api/overview.md → workspace/src/content/docs/api/overview.md
source/guides/overview.md → workspace/src/content/docs/guides/overview.md
```

This allows duplicate filenames in different directories.

#### Test Infrastructure Patterns

**Step definitions** use a shared `TestContext` struct that tracks:
- `sourceDirectory`, `workspacePath`, `targetDirectory`
- `scannedFiles` slice for verification
- `output`, `errorOutput` buffers for CLI output checking
- `exitCode` for status verification
- `tempDirs` slice for cleanup tracking

**Common step patterns**:
- Directory creation with tree structures (uses directory stack for nesting)
- File content verification with frontmatter checking
- CLI execution with output capture
- Workspace validation and cleanup verification

**Tree structure parsing**: Steps that accept Gherkin DocStrings with tree format:
```gherkin
Given a directory "./docs" with files:
  """
  ./docs/
  ├── api/
  │   └── overview.md
  └── guides/
      └── overview.md
  """
```
Use a directory stack to track nesting depth (4 spaces = 1 level) and build correct file paths.

## Development Commands

### Essential Commands

```bash
# Build binary
go build -o bin/flashdoc cmd/flashdoc/main.go

# Run locally with arguments
./bin/flashdoc ./features/test-docs

# Run all BDD tests (89 passing, 124 undefined)
go test -v ./features/...

# Run specific test phase
go test -v ./features/... -run TestPhase1
go test -v ./features/... -run TestPhase2

# Run specific scenario
go test -v ./features/... -run "TestPhase2/Support_various_markdown_extensions"

# Count test results
go test -v ./features/... 2>&1 | grep "^    --- PASS:" | wc -l
go test -v ./features/... 2>&1 | grep "^    --- FAIL:" | wc -l

# Format code
go fmt ./...

# Tidy dependencies
go mod tidy
```

### Testing Workflow

When fixing failing tests:
1. Run the specific failing test to see exact error
2. Read the feature file to understand expected behavior
3. Check step implementation in `features/steps/`
4. Verify/fix production code in `internal/`
5. Rebuild binary if testing end-to-end: `go build -o bin/flashdoc cmd/flashdoc/main.go`
6. Re-run test to verify fix

### Makefile Alternatives

The Makefile exists but some targets may be outdated. Prefer direct Go commands:

```bash
# Instead of: make test
go test -v ./...

# Instead of: make features
go test -v ./features/...

# Instead of: make build
go build -o bin/flashdoc cmd/flashdoc/main.go
```

## Common Patterns & Gotchas

### Godog/Gherkin Patterns

**DocString vs Table**: Steps ending with `:` can accept either:
- **DocString** (triple quotes): Multi-line text, used for tree structures
- **Table** (pipe-delimited): Tabular data with headers

**IMPORTANT**: A single step pattern can only handle ONE format. If you need both, create separate step functions or convert feature files to use consistent format.

Example from recent fix:
```go
// This function handles DocString format only
func (ctx *TestContext) createDirectoryWithFiles(dirPath string, docString *godog.DocString) error {
    // Parse tree structure from docString.Content
}
```

### Tree Structure Parsing

When parsing Gherkin tree structures with `├──`, `└──`, `│` characters:
1. Remove ALL tree characters with `strings.ReplaceAll()`
2. Count leading spaces to determine nesting depth (4 spaces = 1 level)
3. Use a directory stack to track parent directories
4. Build paths by joining stack contents + current name

### Temporary Workspace Management

Tests create temp directories that must be tracked and cleaned up:
```go
ctx.TrackDir(tempDir)  // Registers for cleanup after test
```

The CLI extracts workspace path from output line: `📦 Workspace: /path/to/workspace`

### Package Versions

**Current embedded template versions** (in `internal/template/starlight/package.json`):
- `@astrojs/starlight`: `^0.37.3` (was 0.29.3, upgraded to fix TypeError)
- `astro`: `^5.0.0` (was 4.16.18)
- `sharp`: `^0.34.0` (was 0.33.5)

After changing template files, rebuild binary to embed changes.

## Testing Best Practices

1. **Always read the feature file first** before implementing or debugging
2. **Run specific tests** rather than full suite when debugging
3. **Check step definitions** in `features/steps/` to understand what's being tested
4. **Preserve existing patterns** when adding new step definitions
5. **Track temp directories** in tests with `ctx.TrackDir()`
6. **Use directory stack** for tree structure parsing (4 spaces per level)
7. **Extract paths from CLI output** rather than creating new workspaces in steps

## Known Issues & Recent Fixes

### Recent Fixes (Jan 2026)

1. **Orphaned temp directories test**: Fixed to extract workspace path from CLI output instead of creating new workspace
2. **Markdown extensions test**: Converted feature file from table format to DocString format for consistency
3. **Duplicate filenames test**: Implemented proper directory stack for tree structure parsing
4. **Starlight TypeError**: Upgraded from v0.29.3 to v0.37.3 to fix "Cannot read properties of undefined" error

### Current Test Status

- ✅ 89/89 implemented scenarios passing (100%)
- ⏸️ 124 undefined scenarios (future features: AI generation, file watching, caching)

## References

- **Feature files**: `features/phase{1-4}-*/*.feature` (14 files)
- **Step definitions**: `features/steps/*_steps.go` (16 files)
- **Main entry point**: `cmd/flashdoc/main.go`
- **Internal packages**: `internal/*/` (14 packages)
- **Embedded template**: `internal/template/starlight/`

For contribution guidelines, see `CONTRIBUTING.md`.
