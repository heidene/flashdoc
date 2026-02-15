package steps

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/nicovandenhove/flashdoc/internal/pkgmanager"
	"github.com/nicovandenhove/flashdoc/internal/scanner"
	"github.com/nicovandenhove/flashdoc/internal/server"
	"github.com/nicovandenhove/flashdoc/internal/workspace"
)

// TestContext holds shared state between step definitions
type TestContext struct {
	// CLI execution results
	exitCode    int
	output      *bytes.Buffer
	errorOutput *bytes.Buffer
	cmd         *exec.Cmd

	// Parsed CLI values
	parsedPath  string
	parsedTitle string
	parsedPort  int
	noOpen      bool

	// Workspace tracking
	tempDir       string
	workspacePath string

	// Process tracking
	runningPID int
	childPIDs  []int

	// Flags
	cliAvailable          bool
	interrupted           bool
	tempDirNotWritable    bool
	expectDirectoryMissing bool

	// Test helpers
	createdDirs  []string
	createdFiles []string

	// Phase 2: Markdown Processing
	scannedFiles     []scanner.MarkdownFile
	processedContent string
	sourceDirectory  string
	targetDirectory  string
	copiedFiles      []string

	// Phase 3: Starlight Setup
	extractedFiles  []string
	configContent   string
	detectedPM      pkgmanager.PackageManager
	mockPMAvailable map[string]bool // For mocking availability
	installOutput   string

	// Phase 4: Server & Browser
	server         *server.Server
	serverURL      string
	serverPort     int
	serverReady    bool
	serverStarting bool
	browserOpened  bool
	browserCommand string // Captured for verification
	outputLines    []string

	// Additional state flags
	npmInstalling bool
	slowCleanup   bool

	// Export functionality
	exportPath       string
	buildTriggered   bool
	buildShouldFail  bool
	exportShouldFail bool
	forbiddenPath    string
	expectedTitle    string
}

// NewTestContext creates a new test context
func NewTestContext() *TestContext {
	return &TestContext{
		output:          new(bytes.Buffer),
		errorOutput:     new(bytes.Buffer),
		cliAvailable:    true,
		createdDirs:     make([]string, 0),
		createdFiles:    make([]string, 0),
		childPIDs:       make([]int, 0),
		scannedFiles:    make([]scanner.MarkdownFile, 0),
		copiedFiles:     make([]string, 0),
		extractedFiles:  make([]string, 0),
		mockPMAvailable: make(map[string]bool),
		outputLines:     make([]string, 0),
	}
}

// Reset clears the context for the next scenario
func (ctx *TestContext) Reset() {
	ctx.exitCode = 0
	ctx.output.Reset()
	ctx.errorOutput.Reset()
	ctx.cmd = nil
	ctx.parsedPath = ""
	ctx.parsedTitle = ""
	ctx.parsedPort = 0
	ctx.noOpen = false
	ctx.tempDir = ""
	ctx.workspacePath = ""
	ctx.runningPID = 0
	ctx.childPIDs = make([]int, 0)

	// Phase 2 fields
	ctx.scannedFiles = make([]scanner.MarkdownFile, 0)
	ctx.processedContent = ""
	ctx.sourceDirectory = ""
	ctx.targetDirectory = ""
	ctx.copiedFiles = make([]string, 0)

	// Phase 3 fields
	ctx.extractedFiles = make([]string, 0)
	ctx.configContent = ""
	ctx.detectedPM = ""
	ctx.mockPMAvailable = make(map[string]bool)
	ctx.installOutput = ""

	// Phase 4 fields
	if ctx.server != nil {
		// Server cleanup would happen here
		ctx.server = nil
	}
	ctx.serverURL = ""
	ctx.serverPort = 0
	ctx.serverReady = false
	ctx.browserOpened = false
	ctx.browserCommand = ""
	ctx.outputLines = make([]string, 0)

	// Clean up test files and directories
	for _, file := range ctx.createdFiles {
		os.Remove(file)
	}
	for _, dir := range ctx.createdDirs {
		os.RemoveAll(dir)
	}
	ctx.createdDirs = make([]string, 0)
	ctx.createdFiles = make([]string, 0)
}

// TrackDir adds a directory to be cleaned up after the test
func (ctx *TestContext) TrackDir(path string) {
	ctx.createdDirs = append(ctx.createdDirs, path)
}

// TrackFile adds a file to be cleaned up after the test
func (ctx *TestContext) TrackFile(path string) {
	ctx.createdFiles = append(ctx.createdFiles, path)
}

// CreateTestWorkspace creates a workspace with temporary run and shared directories
func (ctx *TestContext) CreateTestWorkspace() (*workspace.Workspace, error) {
	// Create temp run directory
	runDir, err := os.MkdirTemp("", "stardoc-test-run-")
	if err != nil {
		return nil, fmt.Errorf("failed to create test run directory: %w", err)
	}
	ctx.TrackDir(runDir)

	// Create temp shared directory
	sharedDir, err := os.MkdirTemp("", "stardoc-test-shared-")
	if err != nil {
		return nil, fmt.Errorf("failed to create test shared directory: %w", err)
	}
	ctx.TrackDir(sharedDir)

	// Create workspace
	ws, err := workspace.New(runDir, sharedDir)
	if err != nil {
		return nil, err
	}

	ctx.workspacePath = ws.Path
	return ws, nil
}
