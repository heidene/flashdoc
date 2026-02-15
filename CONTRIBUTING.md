# Contributing to flashdoc

Thank you for your interest in contributing to flashdoc! This document provides guidelines for contributing to the project.

## Project Status

**Current Phase**: Active development and maintenance
**Implementation Status**: ✅ Core features complete (89/89 BDD scenarios passing)
**Next Steps**: Bug fixes, new features, documentation improvements, and community feedback

## Development Philosophy

This project follows **Behavior-Driven Development (BDD)** principles:

1. **Features Define Behavior**: All functionality is specified in Gherkin `.feature` files in `features/`
2. **Tests Drive Development**: BDD tests validate behavior before and after changes
3. **Phase-Based Organization**: Features are organized into 4 phases (foundation, markdown, starlight, server)
4. **Quality First**: All changes must pass existing tests and add tests for new features

## Project Structure

```
flashdoc/
├── features/                    # BDD feature files and test steps
│   ├── phase1-foundation/       # ✅ CLI, workspace, signals, cleanup
│   ├── phase2-markdown/         # ✅ File processing, frontmatter
│   ├── phase3-starlight/        # ✅ Template, config, dependencies
│   ├── phase4-server/           # ✅ Dev server, browser, export
│   ├── steps/                   # Step definitions (16 files)
│   └── test-docs/               # Test fixtures (markdown files)
├── cmd/flashdoc/                 # CLI entry point
├── internal/                    # Internal packages
│   ├── cli/                     # Command-line parsing (Cobra)
│   ├── workspace/               # Temp workspace management
│   ├── scanner/                 # Markdown file discovery
│   ├── frontmatter/             # YAML frontmatter injection
│   ├── processor/               # File processing orchestration
│   ├── template/                # Embedded Starlight template
│   │   └── starlight/           # Astro + Starlight template
│   ├── config/                  # Starlight config generation
│   ├── pkgmanager/              # Package manager detection
│   ├── installer/               # Dependency installation
│   ├── server/                  # Dev server lifecycle
│   ├── exporter/                # Static site export
│   ├── browser/                 # Browser opening
│   ├── cleanup/                 # Resource cleanup
│   └── signal/                  # Signal handling
├── docs/                        # Project documentation
├── .github/workflows/           # CI/CD pipelines
│   ├── ci.yml                   # Tests and linting
│   └── release.yml              # Release automation
├── .goreleaser.yaml             # Release configuration
├── LICENSE                      # Beerware license
├── THIRD-PARTY-LICENSES.md      # Dependency licenses
├── CLAUDE.md                    # Claude Code instructions
├── README.md                    # User documentation
├── CONTRIBUTING.md              # This file
└── Makefile                     # Development commands
```

## Getting Started

### Prerequisites

- **Go 1.22+** (project uses 1.25.5)
- **Node.js 18+** (for Starlight template)
- **Package manager**: pnpm, bun, or npm
- **Git** for version control

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/nicovandenhove/flashdoc.git
cd flashdoc

# Download Go dependencies
go mod download

# Install development tools (optional)
go install github.com/cucumber/godog/cmd/godog@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Development Commands

```bash
# Build the binary
make build
# or: go build -o bin/flashdoc cmd/flashdoc/main.go

# Run all tests (unit + BDD)
go test -v ./...

# Run BDD feature tests only
go test -v ./features/...

# Run specific phase tests
go test -v ./features/... -run TestPhase1
go test -v ./features/... -run TestPhase2

# Run specific scenario
go test -v ./features/... -run "TestPhase4/Export_to_default_directory"

# Format code
go fmt ./...

# Run linters (requires golangci-lint)
make lint
# or: golangci-lint run ./...

# Install locally for testing
make install
# or: go install ./cmd/flashdoc
```

## How to Contribute

### 1. Understanding the Codebase

Start by reading:
- **README.md**: User-facing documentation and usage examples
- **CLAUDE.md**: Comprehensive development guide (BDD workflow, architecture, patterns)
- **Feature files** in `features/phase*/`: Behavior specifications in Gherkin

### 2. Types of Contributions

**Bug Fixes:**
1. Identify the failing behavior
2. Find or write a BDD scenario that reproduces the bug
3. Fix the bug in the relevant `internal/` package
4. Verify all tests pass

**New Features:**
1. Write a `.feature` file describing the behavior
2. Add step definitions in `features/steps/`
3. Implement the feature in `internal/` packages
4. Ensure tests pass

**Documentation:**
- Improve README.md, CONTRIBUTING.md, or comments
- Add examples or tutorials
- Update CLAUDE.md for development patterns

**Refactoring:**
- Ensure all existing tests pass before and after
- Maintain backward compatibility
- Document architectural changes

### 3. Development Workflow

```bash
# 1. Create a feature branch
git checkout -b feature/your-feature-name

# 2. Make your changes
# - Edit code in internal/ packages
# - Add/update tests in features/
# - Update documentation

# 3. Run tests frequently
go test -v ./features/...

# 4. Format and lint
go fmt ./...
golangci-lint run ./...

# 5. Commit with conventional commits format
git add .
git commit -m "feat: add support for custom CSS"

# 6. Push and create PR
git push origin feature/your-feature-name
```

### 4. BDD Testing Pattern

flashdoc uses **godog** for BDD testing. Here's how to add tests:

```go
// features/steps/my_feature_steps.go
package steps

import "github.com/cucumber/godog"

func (ctx *TestContext) iDoSomething() error {
    // Implementation
    return nil
}

func (ctx *TestContext) iShouldSeeResult(expected string) error {
    // Assertion
    if ctx.output != expected {
        return fmt.Errorf("expected %q, got %q", expected, ctx.output)
    }
    return nil
}

// Register steps in features_test.go InitializeScenario
```

See existing step files in `features/steps/` for patterns.

### 5. Code Style

**Go Conventions:**
- Follow standard Go formatting (`gofmt`)
- Use meaningful names (avoid abbreviations)
- Keep functions small and focused
- Document exported types and functions
- Prefer composition over inheritance

**Error Handling:**
```go
// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create workspace: %w", err)
}

// Return early on errors
if err := validate(input); err != nil {
    return err
}
```

**Package Organization:**
- Each package has a single, clear responsibility
- Avoid circular dependencies
- Keep internal packages private

### 6. Testing Guidelines

**BDD Tests (Required):**
- All user-facing behavior must have BDD scenarios
- Test happy paths and error cases
- Use realistic test data in `features/test-docs/`

**Unit Tests (Encouraged):**
- Test complex logic in isolation
- Use table-driven tests for multiple cases
- Mock external dependencies (filesystem, network)

**Manual Testing:**
```bash
# Build and test locally
make build
./bin/flashdoc ./features/test-docs

# Test with real documentation
./bin/flashdoc ~/my-real-docs --title "My Docs"

# Test export functionality
./bin/flashdoc ./features/test-docs --export ./output
```

## Pull Request Process

### Before Submitting

- [ ] All tests pass: `go test -v ./features/...`
- [ ] Code is formatted: `go fmt ./...`
- [ ] Linting passes: `golangci-lint run ./...` (if installed)
- [ ] Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/)
- [ ] Documentation updated (if needed)

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding or updating tests
- `chore`: Changes to build process or auxiliary tools
- `ci`: Changes to CI/CD configuration

**Examples:**
```
feat(export): add --export flag for static site generation

Implements static site export functionality that runs astro build
and copies the output to a specified directory.

Closes #42

fix(frontmatter): handle empty YAML blocks correctly

Previously empty frontmatter blocks caused parser errors.
Now they're treated as empty documents.

docs(readme): add installation instructions for Homebrew

test(scanner): add test for nested directory structures
```

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] All existing tests pass
- [ ] New tests added (if applicable)
- [ ] Manually tested

## Checklist
- [ ] Code follows project style
- [ ] Self-reviewed the code
- [ ] Commented complex logic
- [ ] Updated documentation
- [ ] No new warnings

## Related Issues
Closes #issue_number
```

## Release Process

Releases are automated via GoReleaser and GitHub Actions:

1. **Version Tag**: Create and push a version tag
   ```bash
   git tag -a v0.3.0 -m "Release v0.3.0"
   git push origin v0.3.0
   ```

2. **Automated Build**: GitHub Actions automatically:
   - Runs all tests
   - Builds binaries for multiple platforms
   - Creates GitHub Release
   - Updates Homebrew tap
   - Generates changelog

3. **Distribution**: Users can install via:
   - `go install github.com/nicovandenhove/flashdoc/cmd/flashdoc@latest`
   - `brew install nicovandenhove/tap/flashdoc`
   - Download binaries from GitHub Releases

**Note**: Only maintainers can create releases.

## License and Compliance

### Project License

flashdoc is licensed under the [Beerware License](LICENSE). By contributing, you agree that your contributions will be licensed under the same license.

**What this means:**
- Anyone can use, modify, and distribute this software
- The only requirement: retain the license notice
- Optional: Buy the author a beer if you meet and like it! 🍺

### Third-Party Dependencies

When adding new dependencies:

1. **Check the license**: Must be permissive (MIT, Apache 2.0, BSD)
2. **Update THIRD-PARTY-LICENSES.md**: Add the dependency with license info
3. **Update .goreleaser.yaml**: Include license file in releases (if needed)

See [THIRD-PARTY-LICENSES.md](THIRD-PARTY-LICENSES.md) for current dependencies.

## Getting Help

- **BDD/Testing Questions**: See [CLAUDE.md](CLAUDE.md) sections on BDD patterns
- **Architecture Questions**: Review package documentation in `internal/`
- **Feature Requests**: Open an issue with the `enhancement` label
- **Bug Reports**: Open an issue with reproduction steps
- **General Discussion**: Use GitHub Discussions (if enabled)

## Community Guidelines

- Be respectful and inclusive
- Provide constructive feedback
- Help others learn and grow
- Focus on the behavior, not the person
- Celebrate contributions of all sizes

## Recognition

Contributors will be recognized in:
- GitHub contributor list
- Release notes (for significant contributions)
- README.md (for major features)

---

**Thank you for contributing to flashdoc!** 🚀

Every contribution—whether it's code, documentation, bug reports, or feedback—helps make flashdoc better for everyone.

**Questions?** Open an issue or reach out to [@heidene](https://github.com/heidene).

---

*Last updated: 2026-02-15*
