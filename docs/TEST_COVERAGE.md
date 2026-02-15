# flashdoc Test Coverage Report

**Total**: 213 scenarios (89 implemented ✅, 124 undefined ⏸️)

This document provides a comprehensive overview of all BDD scenarios defined in Gherkin feature files, showing which scenarios have been implemented and tested.

---

## Phase 1: Foundation (38 scenarios - 33 implemented, 5 undefined)

### CLI Parsing
**Status**: 12/12 implemented ✅

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Valid directory path provided | ✅ | Yes |
| Absolute path provided | ✅ | Yes |
| No arguments provided | ✅ | Yes |
| Too many arguments provided | ✅ | Yes |
| Directory does not exist | ✅ | Yes |
| Path is a file not a directory | ✅ | Yes |
| Help flag | ✅ | Yes |
| Version flag | ✅ | Yes |
| Custom title flag | ✅ | Yes |
| Port flag | ✅ | Yes |
| Invalid port flag | ✅ | Yes |
| No-open flag | ✅ | Yes |

### Temporary Workspace Management
**Status**: 7/9 implemented (2 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Create temporary directory on startup | ✅ | Yes |
| Temp directory contains Starlight structure | ✅ | Yes |
| Temp directory location is deterministic | ✅ | Yes |
| Reuse existing temp directory if still valid | ✅ | Yes |
| Handle temp directory creation failure | ⏸️ | No |
| Temp directory isolation | ⏸️ | No |
| Temp directory cleanup on successful exit | ✅ | Yes |
| Temp directory cleanup on error during setup | ✅ | Yes |
| Handle temp directory creation failure | ✅ | Yes |

### Signal Handling
**Status**: 3/8 implemented (5 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Handle SIGINT (Ctrl+C) | ✅ | Yes |
| Handle SIGTERM | ✅ | Yes |
| Multiple interrupt signals | ✅ | Yes |
| Interrupt during dependency installation | ⏸️ | No |
| Interrupt during server startup | ⏸️ | No |
| Interrupt with child processes running | ⏸️ | No |
| Cleanup timeout handling | ⏸️ | No |
| Signal handling does not interfere with normal operation | ⏸️ | No |

### Cleanup
**Status**: 10/10 implemented ✅

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Normal cleanup on exit | ✅ | Yes |
| Cleanup on early error | ✅ | Yes |
| Cleanup child processes | ✅ | Yes |
| Cleanup with multiple running processes | ✅ | Yes |
| Cleanup on panic | ✅ | Yes |
| Partial cleanup on force exit | ✅ | Yes |
| Cleanup verification logging | ✅ | Yes |
| Handle cleanup errors gracefully | ✅ | Yes |
| No orphaned temp directories from previous runs | ✅ | Yes |
| Cleanup node_modules if present | ✅ | Yes |

---

## Phase 2: Markdown Processing (44 scenarios - 24 implemented, 20 undefined)

### Frontmatter Injection
**Status**: 14/18 implemented (4 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Add frontmatter to file without any | ✅ | Yes |
| Preserve existing frontmatter | ✅ | Yes |
| Add missing title to existing frontmatter | ✅ | Yes |
| Generate title from filename | ✅ | Yes |
| Generate title from kebab-case filename | ✅ | Yes |
| Generate title from snake_case filename | ✅ | Yes |
| Handle numbered prefixes in filenames | ✅ | Yes |
| Handle special characters in filenames | ✅ | Yes |
| Handle index files | ✅ | Yes |
| Handle README files | ✅ | Yes |
| Handle README in root | ✅ | Yes |
| Preserve existing title even with misleading filename | ✅ | Yes |
| Handle malformed frontmatter | ⏸️ | No |
| Handle empty frontmatter | ✅ | Yes |
| Preserve other frontmatter fields | ✅ | Yes |
| Handle unicode in titles | ⏸️ | No |
| Strip file extension from generated titles | ✅ | Yes |

### File Copying
**Status**: 7/19 implemented (12 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Copy files to Starlight content directory | ✅ | Yes |
| Preserve directory structure | ✅ | Yes |
| Rename README.md to index.md | ✅ | Yes |
| Rename nested README files | ✅ | Yes |
| Create necessary parent directories | ✅ | Yes |
| Handle file copy errors | ⏸️ | No |
| Process files before copying | ✅ | Yes |
| Preserve file timestamps | ⏸️ | No |
| Handle special characters in filenames | ✅ | Yes |
| Copy files in correct order | ⏸️ | No |
| Skip non-markdown files | ⏸️ | No |
| Handle duplicate filenames in different directories | ⏸️ | No |
| Report copy statistics | ⏸️ | No |
| Handle unicode filenames | ⏸️ | No |
| Atomic copy prevents partial files | ⏸️ | No |
| Interrupted copy cleanup | ⏸️ | No |

### Folder Scanning
**Status**: 6/11 implemented (5 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Scan directory with markdown files | ✅ | Yes |
| Scan nested directories | ✅ | Yes |
| Handle empty directory | ✅ | Yes |
| Ignore non-markdown files | ✅ | Yes |
| Support various markdown extensions | ✅ | Yes |
| Ignore hidden files and directories | ✅ | Yes |
| Ignore common exclude patterns | ⏸️ | No |
| Preserve file order for navigation | ⏸️ | No |
| Handle symbolic links | ⏸️ | No |
| Handle scan errors gracefully | ⏸️ | No |
| Report scan summary | ⏸️ | No |

---

## Phase 3: Starlight Setup (69 scenarios - 31 implemented, 38 undefined)

### Template Extraction
**Status**: 10/12 implemented (2 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Extract template to temp workspace | ✅ | Yes |
| Template contains Starlight dependencies | ✅ | Yes |
| Template package.json has correct structure | ✅ | Yes |
| Template astro.config.mjs is valid | ✅ | Yes |
| Template directory structure matches Starlight conventions | ✅ | Yes |
| Template uses latest stable Starlight version | ✅ | Yes |
| Template extraction is idempotent | ✅ | Yes |
| Handle extraction errors | ⏸️ | No |
| Template contains minimal boilerplate | ✅ | Yes |
| Template tsconfig.json extends Starlight defaults | ✅ | Yes |
| Template is embedded using go:embed | ✅ | Yes |
| Template supports custom Starlight config | ⏸️ | No |

### Config Generation
**Status**: 10/16 implemented (6 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Generate config with default title | ✅ | Yes |
| Generate config with custom title flag | ✅ | Yes |
| Replace template placeholder | ✅ | Yes |
| Generate title from directory name with hyphens | ✅ | Yes |
| Generate title from directory name with underscores | ✅ | Yes |
| Handle special characters in auto-generated title | ✅ | Yes |
| Generate config with sidebar autogeneration | ✅ | Yes |
| Generate config without social links | ⏸️ | No |
| Generate config with default locale handling | ⏸️ | No |
| Config is valid JavaScript | ✅ | Yes |
| Handle title with quotes | ⏸️ | No |
| Handle very long titles | ✅ | Yes |
| Generate config preserves template structure | ✅ | Yes |
| Config generation error handling | ⏸️ | No |
| Config includes documentation URL | ⏸️ | No |
| Config supports future customization | ⏸️ | No |

### Package Manager Detection
**Status**: 5/14 implemented (9 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Detect pnpm when available | ✅ | Yes |
| Detect bun when pnpm is not available | ✅ | Yes |
| Fall back to npm | ✅ | Yes |
| Detection priority order | ✅ | Yes |
| No package manager available | ✅ | Yes |
| Package manager detection via version check | ⏸️ | No |
| Respect package manager override flag | ⏸️ | No |
| Invalid package manager override | ⏸️ | No |
| Package manager detection caching | ⏸️ | No |
| Package manager in PATH but not executable | ⏸️ | No |
| Detect package manager version compatibility | ⏸️ | No |
| Warn about old package manager versions | ⏸️ | No |
| Package manager detection on different platforms | ⏸️ | No |
| Log detection process | ⏸️ | No |

### Dependency Installation
**Status**: 6/20 implemented (14 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Install dependencies with pnpm | ✅ | Yes |
| Install dependencies with bun | ✅ | Yes |
| Install dependencies with npm | ✅ | Yes |
| Display installation progress | ✅ | Yes |
| Successful installation | ✅ | Yes |
| Installation failure | ✅ | Yes |
| Network error during installation | ⏸️ | No |
| Corrupted package.json | ⏸️ | No |
| Install with frozen lockfile (pnpm) | ⏸️ | No |
| Install with frozen lockfile (npm) | ⏸️ | No |
| Install with frozen lockfile (bun) | ⏸️ | No |
| No lockfile present | ⏸️ | No |
| Installation timeout | ⏸️ | No |
| Handle SIGINT during installation | ⏸️ | No |
| Retry on transient failure | ⏸️ | No |
| Silent install option | ⏸️ | No |
| Install creates node_modules | ⏸️ | No |
| Install duration logging | ⏸️ | No |
| Parallel installations are isolated | ⏸️ | No |
| Installation with npm audit warnings | ⏸️ | No |

---

## Phase 4: Server & UX (62 scenarios - 1 implemented, 61 undefined)

### Dev Server
**Status**: 0/19 implemented (19 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Start dev server with default port | ⏸️ | No |
| Dev server starts successfully | ⏸️ | No |
| Stream server output | ⏸️ | No |
| Server startup with custom port | ⏸️ | No |
| Port already in use | ⏸️ | No |
| Server startup failure | ⏸️ | No |
| Monitor server process | ⏸️ | No |
| Server crashes during runtime | ⏸️ | No |
| Server responds to requests | ⏸️ | No |
| Hot reload works | ⏸️ | No |
| Server shutdown on exit | ⏸️ | No |
| Force kill server if graceful shutdown fails | ⏸️ | No |
| Server uses correct package manager | ⏸️ | No |
| Server starts with environment variables | ⏸️ | No |
| Parse server URL from output | ⏸️ | No |
| Handle non-standard Astro output | ⏸️ | No |
| Server keeps running until interrupted | ⏸️ | No |
| Multiple dev servers in parallel | ⏸️ | No |
| Dev server uses production-like settings | ⏸️ | No |

### Browser Open
**Status**: 0/16 implemented (16 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Open browser on macOS | ⏸️ | No |
| Open browser on Linux | ⏸️ | No |
| Open browser on Windows | ⏸️ | No |
| Wait for server before opening browser | ⏸️ | No |
| Browser opens after server ready log | ⏸️ | No |
| Do not open browser with --no-open flag | ⏸️ | No |
| Handle browser open failure | ⏸️ | No |
| Browser is already running | ⏸️ | No |
| Open with custom port | ⏸️ | No |
| Open URL with base path | ⏸️ | No |
| Log browser open action | ⏸️ | No |
| Network interface binding | ⏸️ | No |
| Browser open timeout | ⏸️ | No |
| Concurrent browser opens | ⏸️ | No |
| Browser preference detection | ⏸️ | No |
| Headless environment detection | ⏸️ | No |
| Open browser exactly once | ⏸️ | No |

### Terminal Output
**Status**: 1/22 implemented (21 undefined)

| Scenario | Status | Implemented |
|----------|--------|-------------|
| Clean startup output | ⏸️ | No |
| Use emoji consistently | ⏸️ | No |
| Progress indicators for long operations | ⏸️ | No |
| Streaming output from subprocesses | ⏸️ | No |
| Error output is distinct | ⏸️ | No |
| Warning output is distinct | ⏸️ | No |
| Verbose mode | ⏸️ | No |
| Quiet mode | ⏸️ | No |
| Color support detection | ⏸️ | No |
| No color support | ⏸️ | No |
| NO_COLOR environment variable | ⏸️ | No |
| Interactive vs non-interactive detection | ⏸️ | No |
| Logging format options | ⏸️ | No |
| Timestamp option | ⏸️ | No |
| File processing summary | ⏸️ | No |
| Cleanup messages | ⏸️ | No |
| Error context and suggestions | ⏸️ | No |
| Help text formatting | ⏸️ | No |
| Version output | ⏸️ | No |
| Output buffering | ⏸️ | No |
| Multiline error messages | ⏸️ | No |
| Server output passthrough | ⏸️ | No |
| Handle terminal width | ✅ | Yes |
| Exit message clarity | ⏸️ | No |

---

## Summary by Phase

| Phase | Feature | Total | Implemented | Undefined | % Complete |
|-------|---------|-------|-------------|-----------|------------|
| **Phase 1** | **Foundation** | **38** | **33** | **5** | **87%** |
| | CLI Parsing | 12 | 12 | 0 | 100% |
| | Temp Workspace | 9 | 7 | 2 | 78% |
| | Signal Handling | 8 | 3 | 5 | 38% |
| | Cleanup | 10 | 10 | 0 | 100% |
| **Phase 2** | **Markdown** | **44** | **24** | **20** | **55%** |
| | Frontmatter Injection | 18 | 14 | 4 | 78% |
| | File Copying | 19 | 7 | 12 | 37% |
| | Folder Scanning | 11 | 6 | 5 | 55% |
| **Phase 3** | **Starlight** | **69** | **31** | **38** | **45%** |
| | Template Extraction | 12 | 10 | 2 | 83% |
| | Config Generation | 16 | 10 | 6 | 63% |
| | Package Manager Detection | 14 | 5 | 9 | 36% |
| | Dependency Installation | 20 | 6 | 14 | 30% |
| **Phase 4** | **Server & UX** | **62** | **1** | **61** | **2%** |
| | Dev Server | 19 | 0 | 19 | 0% |
| | Browser Open | 16 | 0 | 16 | 0% |
| | Terminal Output | 22 | 1 | 21 | 5% |
| **Total** | | **213** | **89** | **124** | **42%** |

---

## Notes

- ✅ **Implemented**: Scenario has passing tests
- ⏸️ **Undefined**: Scenario defined in feature file but not yet implemented

### How to Run Tests

```bash
# Run all tests
go test -v ./features/...

# Run specific phase
go test -v ./features/... -run TestPhase1
go test -v ./features/... -run TestPhase2
go test -v ./features/... -run TestPhase3
go test -v ./features/... -run TestPhase4

# Run specific scenario (example)
go test -v ./features/... -run "TestPhase1/Valid_directory_path_provided"
```

### Test File Locations

- **Feature files**: `features/phase{1-4}-*/*.feature`
- **Step definitions**: `features/steps/*_steps.go`
- **Test runner**: `features/features_test.go`

---

*Last updated: 2026-01-21*
