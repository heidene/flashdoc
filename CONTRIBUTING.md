# Contributing to Stardoc

Thank you for your interest in contributing to Stardoc! This document provides guidelines for contributing to the project.

## Project Status

**Current Phase**: Feature validation and BDD test setup
**Implementation Status**: Not yet started - all feature files are complete and ready for implementation

## Development Philosophy

This project follows **Behavior-Driven Development (BDD)** principles:

1. **Features Define Behavior**: All functionality is specified in Gherkin `.feature` files
2. **Tests Before Code**: Write tests that match feature scenarios before implementing
3. **Phase-by-Phase**: Implement incrementally, completing one phase before moving to the next
4. **Validation First**: Features are validated and approved before implementation begins

## Project Structure

```
stardoc/
├── features/                    # BDD feature files (specifications)
│   ├── phase1-foundation/      # Core CLI, workspace, cleanup
│   ├── phase2-markdown/        # File processing and frontmatter
│   ├── phase3-starlight/       # Template and dependency management
│   └── phase4-server/          # Dev server and browser integration
├── cmd/stardoc/                # CLI entry point
├── internal/                   # Internal packages
│   ├── cli/                    # Command-line parsing and validation
│   ├── server/                 # Dev server lifecycle management
│   ├── markdown/               # Markdown file processing
│   ├── template/               # Embedded template extraction
│   └── pkgmgr/                 # Package manager detection
├── templates/starlight/        # Embedded Starlight template (to be created)
├── README.md                   # Project overview
├── FEATURE_VALIDATION.md       # Feature review and validation checklist
├── Makefile                    # Development commands
└── go.mod                      # Go module definition
```

## Getting Started

### Prerequisites

- Go 1.21 or later
- Node.js 18 or later (for testing Starlight integration)
- One of: pnpm, bun, or npm

### Initial Setup

```bash
# Clone the repository
cd ~/Code/Tools/stardoc

# Download Go dependencies
make mod

# Install development tools
make deps
```

### Development Commands

```bash
# Build the binary
make build

# Run tests (once implemented)
make test

# Run BDD feature tests
make features

# Format code
make fmt

# Run linters
make lint

# Install locally for testing
make install
```

## Implementation Roadmap

### Phase 1: Foundation (Current Target)
Focus: CLI parsing, temp workspace, signal handling, cleanup

**Features to implement**:
1. `cli-parsing.feature` - Cobra-based CLI with flags
2. `temp-workspace.feature` - OS temp directory management
3. `signal-handling.feature` - Graceful shutdown with signals
4. `cleanup.feature` - Resource cleanup on exit

**Implementation steps**:
1. Set up `internal/cli` package with Cobra
2. Implement workspace creation in dedicated package
3. Add signal handling with context cancellation
4. Wire up cleanup with defer and signal handlers
5. Write godog tests for each feature
6. Ensure all scenarios pass

### Phase 2: Markdown Processing
Focus: File discovery, frontmatter injection, copying

### Phase 3: Starlight Integration
Focus: Template extraction, config generation, dependencies

### Phase 4: Server & UX
Focus: Dev server management, browser launch, terminal output

## How to Contribute

### 1. Review Feature Files

Start by reading the feature files in the `features/` directory. These define the expected behavior.

Each `.feature` file contains:
- **Feature**: High-level description
- **Background**: Common setup for scenarios
- **Scenarios**: Specific test cases with Given/When/Then steps

Example:
```gherkin
Scenario: Valid directory path provided
  When I run "stardoc ./my-docs"
  Then the CLI should parse the path "./my-docs"
  And the CLI should validate that the path exists
  And the CLI should proceed with site generation
```

### 2. Set Up BDD Tests

We use [godog](https://github.com/cucumber/godog) for BDD testing.

```bash
# Install godog
go install github.com/cucumber/godog/cmd/godog@latest

# Create test file structure (example for Phase 1)
mkdir -p features/phase1-foundation
# Feature files already exist!

# Create step definitions
mkdir -p test/steps
touch test/steps/cli_steps.go
```

### 3. Write Step Definitions

Step definitions map Gherkin steps to Go code:

```go
// test/steps/cli_steps.go
package steps

import (
    "github.com/cucumber/godog"
)

func InitializeScenario(ctx *godog.ScenarioContext) {
    ctx.Step(`^I run "([^"]*)"$`, iRun)
    ctx.Step(`^the CLI should parse the path "([^"]*)"$`, cliShouldParsePath)
    // ... more step definitions
}

func iRun(command string) error {
    // Implementation
    return nil
}

func cliShouldParsePath(path string) error {
    // Assertion
    return nil
}
```

### 4. Implement Features

With tests in place, implement the actual functionality:

```go
// internal/cli/parser.go
package cli

import "github.com/spf13/cobra"

func NewRootCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "stardoc <directory>",
        Short: "Ephemeral Starlight documentation sites",
        Args:  cobra.ExactArgs(1),
        RunE:  run,
    }

    // Add flags
    cmd.Flags().String("title", "", "Custom site title")
    cmd.Flags().Int("port", 4321, "Dev server port")
    cmd.Flags().Bool("no-open", false, "Don't open browser")

    return cmd
}

func run(cmd *cobra.Command, args []string) error {
    // Implementation
    return nil
}
```

### 5. Run Tests

```bash
# Run all feature tests
make features

# Run specific feature
godog features/phase1-foundation/cli-parsing.feature

# Run with tags
godog --tags=@phase1 features/
```

### 6. Iterate

- Write test → Implement → Verify → Refine
- Ensure all scenarios pass before moving on
- Keep commits focused on single features

## Code Style

### Go Conventions

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Document exported functions and types
- Keep functions small and focused
- Prefer composition over inheritance

### Error Handling

```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create workspace: %w", err)
}

// Use custom error types for specific cases
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}
```

### Package Organization

```
internal/
├── cli/           # CLI parsing and command setup
├── workspace/     # Temp directory management
├── scanner/       # Markdown file discovery
├── processor/     # Frontmatter injection
├── copier/        # File copying logic
├── template/      # Template extraction
├── config/        # Config generation
├── pkgmgr/        # Package manager detection
├── installer/     # Dependency installation
├── server/        # Dev server management
└── browser/       # Browser opening
```

## Testing Guidelines

### Unit Tests
- Test individual functions and methods
- Mock external dependencies
- Use table-driven tests where appropriate

```go
func TestGenerateTitle(t *testing.T) {
    tests := []struct {
        name     string
        filename string
        want     string
    }{
        {"simple", "guide.md", "Guide"},
        {"hyphenated", "getting-started.md", "Getting Started"},
        {"numbered", "01-intro.md", "Intro"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := GenerateTitle(tt.filename)
            if got != tt.want {
                t.Errorf("got %q, want %q", got, tt.want)
            }
        })
    }
}
```

### Integration Tests (BDD)
- Use godog to implement feature scenarios
- Test component interactions
- Verify end-to-end behavior

### Manual Testing
- Build and test the actual binary
- Verify cross-platform behavior
- Test with real markdown directories

## Pull Request Process

1. **Fork and Branch**: Create a feature branch from `main`
2. **Implement**: Write tests first, then implementation
3. **Test**: Ensure all tests pass (`make test`)
4. **Format**: Run `make fmt` and `make lint`
5. **Commit**: Use clear, descriptive commit messages
6. **PR**: Submit PR with reference to feature file(s)

### Commit Message Format

```
type(scope): subject

body (optional)

Refs: #issue
```

Examples:
```
feat(cli): implement directory path validation

Implements scenarios from cli-parsing.feature:
- Valid directory path provided
- Directory does not exist
- Path is a file not a directory

Refs: #1

test(workspace): add temp directory cleanup tests

Covers cleanup.feature scenarios for graceful shutdown

docs(readme): update installation instructions
```

## Questions?

- **Feature Questions**: Review the feature files and FEATURE_VALIDATION.md
- **Implementation Questions**: Open an issue for discussion
- **Bug Reports**: Use issue template (to be created)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Happy Coding!** 🚀

Remember: Features first, tests second, implementation third. This approach ensures we build exactly what's specified and nothing more.
