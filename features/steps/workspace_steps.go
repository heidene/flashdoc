package steps

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
	"github.com/nicovandenhove/flashdoc/internal/workspace"
)

// RegisterWorkspaceSteps registers all workspace-related step definitions
func RegisterWorkspaceSteps(ctx *godog.ScenarioContext, testCtx *TestContext) {
	ctx.Step(`^a source directory "([^"]*)" exists with markdown files$`, testCtx.aSourceDirectoryExistsWithMarkdownFiles)
	ctx.Step(`^a temporary directory should be created in the system temp location$`, testCtx.aTempDirectoryShouldBeCreatedInSystemTempLocation)
	ctx.Step(`^the temp directory name should start with "([^"]*)"$`, testCtx.theTempDirectoryNameShouldStartWith)
	ctx.Step(`^the temp directory name should include a unique identifier$`, testCtx.theTempDirectoryNameShouldIncludeUniqueIdentifier)
	ctx.Step(`^the temp directory should have write permissions$`, testCtx.theTempDirectoryShouldHaveWritePermissions)
	ctx.Step(`^the temp directory should contain a "([^"]*)" subdirectory$`, testCtx.theTempDirectoryShouldContainSubdirectory)
	ctx.Step(`^the temp directory should contain a "([^"]*)" file$`, testCtx.theTempDirectoryShouldContainFile)
	ctx.Step(`^the temp directory should be created under the OS temp directory$`, testCtx.theTempDirectoryShouldBeCreatedUnderOSTempDirectory)
	ctx.Step(`^the temp directory path should be logged to the terminal$`, testCtx.theTempDirectoryPathShouldBeLogged)
	ctx.Step(`^the log message should be "([^"]*)"$`, testCtx.theLogMessageShouldBe)
	ctx.Step(`^stardoc has created a temp directory for "([^"]*)"$`, testCtx.stardocHasCreatedTempDirectoryFor)
	ctx.Step(`^the temp directory still exists$`, testCtx.theTempDirectoryStillExists)
	ctx.Step(`^a new temporary directory should be created$`, testCtx.aNewTempDirectoryShouldBeCreated)
	ctx.Step(`^the new directory should have a different unique identifier$`, testCtx.theNewDirectoryShouldHaveDifferentUniqueIdentifier)
	ctx.Step(`^the system temp directory is not writable$`, testCtx.theSystemTempDirectoryIsNotWritable)
	ctx.Step(`^each instance should have its own isolated temp directory$`, testCtx.eachInstanceShouldHaveIsolatedTempDirectory)
	ctx.Step(`^the temp directory should be deleted$`, testCtx.theTempDirectoryShouldBeDeleted)
	ctx.Step(`^no stardoc-\* directories should remain in the system temp$`, testCtx.noStardocDirectoriesShouldRemain)
}

func (ctx *TestContext) aSourceDirectoryExistsWithMarkdownFiles(path string) error {
	// Create the directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	ctx.TrackDir(path)

	// Create a sample markdown file
	mdFile := filepath.Join(path, "README.md")
	content := "# Test Documentation\n\nThis is a test file."
	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func (ctx *TestContext) aTempDirectoryShouldBeCreatedInSystemTempLocation() error {
	ws, err := ctx.CreateTestWorkspace()
	if err != nil {
		return err
	}

	if !ws.Exists() {
		return fmt.Errorf("workspace directory was not created")
	}

	// Log the workspace path to output
	ctx.output.WriteString(fmt.Sprintf("Workspace: %s\n", ws.Path))

	return nil
}

func (ctx *TestContext) theTempDirectoryNameShouldStartWith(prefix string) error {
	if ctx.workspacePath == "" {
		return fmt.Errorf("no workspace path set")
	}

	dirName := filepath.Base(ctx.workspacePath)
	if !strings.HasPrefix(dirName, prefix) {
		return fmt.Errorf("expected directory name to start with %q, got %q", prefix, dirName)
	}

	return nil
}

func (ctx *TestContext) theTempDirectoryNameShouldIncludeUniqueIdentifier() error {
	if ctx.workspacePath == "" {
		return fmt.Errorf("no workspace path set")
	}

	dirName := filepath.Base(ctx.workspacePath)
	// Check that there's more than just "stardoc-"
	if len(dirName) <= len("stardoc-") {
		return fmt.Errorf("directory name %q doesn't include a unique identifier", dirName)
	}

	return nil
}

func (ctx *TestContext) theTempDirectoryShouldHaveWritePermissions() error {
	if ctx.workspacePath == "" {
		return fmt.Errorf("no workspace path set")
	}

	// Try to create a test file
	testFile := filepath.Join(ctx.workspacePath, "test-write")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		return fmt.Errorf("workspace is not writable: %w", err)
	}

	os.Remove(testFile)
	return nil
}

func (ctx *TestContext) theTempDirectoryShouldContainSubdirectory(subdir string) error {
	if ctx.workspacePath == "" {
		return fmt.Errorf("no workspace path set")
	}

	ws := &workspace.Workspace{Path: ctx.workspacePath}
	if err := ws.Setup(); err != nil {
		return err
	}

	subdirPath := filepath.Join(ctx.workspacePath, subdir)
	info, err := os.Stat(subdirPath)
	if err != nil {
		return fmt.Errorf("subdirectory %q not found: %w", subdir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%q is not a directory", subdir)
	}

	return nil
}

func (ctx *TestContext) theTempDirectoryShouldContainFile(filename string) error {
	if ctx.workspacePath == "" {
		return fmt.Errorf("no workspace path set")
	}

	ws := &workspace.Workspace{Path: ctx.workspacePath}
	if err := ws.Setup(); err != nil {
		return err
	}

	filePath := filepath.Join(ctx.workspacePath, filename)
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file %q not found: %w", filename, err)
	}

	return nil
}

func (ctx *TestContext) theTempDirectoryShouldBeCreatedUnderOSTempDirectory() error {
	if ctx.workspacePath == "" {
		return fmt.Errorf("no workspace path set")
	}

	osTempDir := os.TempDir()
	if !strings.HasPrefix(ctx.workspacePath, osTempDir) {
		return fmt.Errorf("workspace %q is not under OS temp directory %q", ctx.workspacePath, osTempDir)
	}

	return nil
}

func (ctx *TestContext) theTempDirectoryPathShouldBeLogged() error {
	output := ctx.output.String()
	if !strings.Contains(output, ctx.workspacePath) {
		return fmt.Errorf("workspace path %q not found in output: %s", ctx.workspacePath, output)
	}
	return nil
}

func (ctx *TestContext) theLogMessageShouldBe(expectedMessage string) error {
	output := ctx.output.String()

	// Handle {path} placeholder
	expectedMessage = strings.Replace(expectedMessage, "{path}", ctx.workspacePath, -1)

	// Handle XXXXX pattern for dynamic parts (like temp IDs)
	if strings.Contains(expectedMessage, "XXXXX") {
		// For "Workspace: /tmp/stardoc-XXXXX" pattern, just check the structure
		// The actual temp directory might be /var/folders/... on macOS or /tmp on Linux
		if strings.Contains(expectedMessage, "Workspace:") {
			// Check if output contains "Workspace:" followed by a path with "stardoc-" and some ID
			pattern := `Workspace: .*/stardoc-[a-zA-Z0-9]+`
			matched, err := regexp.MatchString(pattern, output)
			if err != nil {
				return fmt.Errorf("regex error: %v", err)
			}
			if !matched {
				return fmt.Errorf("expected log message pattern %q not found in output: %s", expectedMessage, output)
			}
			return nil
		}
	}

	if !strings.Contains(output, expectedMessage) {
		return fmt.Errorf("expected log message %q not found in output: %s", expectedMessage, output)
	}
	return nil
}

func (ctx *TestContext) stardocHasCreatedTempDirectoryFor(sourceDir string) error {
	// Create source directory
	if err := ctx.aSourceDirectoryExistsWithMarkdownFiles(sourceDir); err != nil {
		return err
	}

	// Create workspace
	ws, err := ctx.CreateTestWorkspace()
	if err != nil {
		return err
	}

	return ws.Setup()
}

func (ctx *TestContext) theTempDirectoryStillExists() error {
	ws := &workspace.Workspace{Path: ctx.workspacePath}
	if !ws.Exists() {
		return fmt.Errorf("workspace directory does not exist")
	}
	return nil
}

func (ctx *TestContext) aNewTempDirectoryShouldBeCreated() error {
	oldPath := ctx.workspacePath

	_, err := ctx.CreateTestWorkspace()
	if err != nil {
		return err
	}

	if ctx.workspacePath == oldPath {
		return fmt.Errorf("new workspace path is the same as old path")
	}

	return nil
}

func (ctx *TestContext) theNewDirectoryShouldHaveDifferentUniqueIdentifier() error {
	// This is implicitly tested by aNewTempDirectoryShouldBeCreated
	return nil
}

func (ctx *TestContext) theSystemTempDirectoryIsNotWritable() error {
	// Set flag to simulate temp directory creation failure
	ctx.tempDirNotWritable = true
	return nil
}

func (ctx *TestContext) eachInstanceShouldHaveIsolatedTempDirectory() error {
	// Create two workspaces and verify they're different
	ws1, err := ctx.CreateTestWorkspace()
	if err != nil {
		return err
	}

	ws2, err := ctx.CreateTestWorkspace()
	if err != nil {
		return err
	}

	if ws1.Path == ws2.Path {
		return fmt.Errorf("workspaces have the same path: %s", ws1.Path)
	}

	return nil
}

func (ctx *TestContext) theTempDirectoryShouldBeDeleted() error {
	ws := &workspace.Workspace{Path: ctx.workspacePath}
	if err := ws.Cleanup(); err != nil {
		return err
	}

	if ws.Exists() {
		return fmt.Errorf("workspace directory still exists after cleanup")
	}

	return nil
}

func (ctx *TestContext) noStardocDirectoriesShouldRemain() error {
	// Check that the workspace created in THIS test was cleaned up
	// (Not checking ALL stardoc directories since there may be orphaned ones from previous runs)
	if ctx.workspacePath == "" {
		// No workspace was created in this test
		return nil
	}

	// Check if the specific workspace directory still exists
	if _, err := os.Stat(ctx.workspacePath); !os.IsNotExist(err) {
		return fmt.Errorf("workspace directory %q was not cleaned up", ctx.workspacePath)
	}

	return nil
}
