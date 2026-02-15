package steps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MockCommandExecutor provides mock command execution for testing
type MockCommandExecutor struct {
	AvailableCommands map[string]bool
	CommandOutputs    map[string]string
	ExecutedCommands  []string
}

// NewMockCommandExecutor creates a new mock command executor
func NewMockCommandExecutor() *MockCommandExecutor {
	return &MockCommandExecutor{
		AvailableCommands: make(map[string]bool),
		CommandOutputs:    make(map[string]string),
		ExecutedCommands:  make([]string, 0),
	}
}

// IsAvailable checks if a command is available
func (m *MockCommandExecutor) IsAvailable(command string) bool {
	return m.AvailableCommands[command]
}

// Execute simulates command execution and returns output
func (m *MockCommandExecutor) Execute(command string, args ...string) (string, error) {
	fullCommand := command + " " + strings.Join(args, " ")
	m.ExecutedCommands = append(m.ExecutedCommands, fullCommand)

	if output, ok := m.CommandOutputs[fullCommand]; ok {
		return output, nil
	}

	return "", nil
}

// CreateDirectoryStructure creates a directory structure from a map
// map keys are paths (relative to baseDir), values are file contents
// Empty string values create directories
func (ctx *TestContext) CreateDirectoryStructure(baseDir string, structure map[string]string) error {
	for path, content := range structure {
		fullPath := filepath.Join(baseDir, path)

		// If content is empty and path ends with /, it's a directory
		if content == "" && strings.HasSuffix(path, "/") {
			if err := os.MkdirAll(fullPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
			}
			ctx.TrackDir(fullPath)
		} else {
			// It's a file - create parent directory first
			dir := filepath.Dir(fullPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create parent directory %s: %w", dir, err)
			}

			// Create the file
			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to create file %s: %w", fullPath, err)
			}
			ctx.TrackFile(fullPath)
		}
	}
	return nil
}

// CreateStandardDocStructure creates a standard documentation structure for testing
func (ctx *TestContext) CreateStandardDocStructure() error {
	baseDir, err := os.MkdirTemp("", "stardoc-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	ctx.TrackDir(baseDir)
	ctx.sourceDirectory = baseDir

	structure := map[string]string{
		"README.md": "# Test Documentation\n\nThis is a test.",
		"docs/":     "",
		"docs/guide.md": `# Guide

## Introduction
This is a guide.

## Usage
Use it like this.
`,
		"docs/api.md": `# API Reference

## Methods
- method1()
- method2()
`,
		"docs/advanced/": "",
		"docs/advanced/performance.md": `# Performance

Tips for better performance.
`,
	}

	return ctx.CreateDirectoryStructure(baseDir, structure)
}

// CreateMarkdownFile creates a markdown file with given content
func (ctx *TestContext) CreateMarkdownFile(dir, filename, content string) (string, error) {
	if dir == "" {
		var err error
		dir, err = os.MkdirTemp("", "stardoc-test-*")
		if err != nil {
			return "", fmt.Errorf("failed to create temp directory: %w", err)
		}
		ctx.TrackDir(dir)
	}

	fullPath := filepath.Join(dir, filename)
	parentDir := filepath.Dir(fullPath)

	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	ctx.TrackFile(fullPath)
	return fullPath, nil
}

// SetupStarlightWorkspace creates a mock Starlight workspace for testing
func (ctx *TestContext) SetupStarlightWorkspace() error {
	baseDir, err := os.MkdirTemp("", "stardoc-starlight-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	ctx.TrackDir(baseDir)
	ctx.tempDir = baseDir

	structure := map[string]string{
		"package.json": `{
  "name": "test-starlight-site",
  "type": "module",
  "version": "0.0.1",
  "scripts": {
    "dev": "astro dev",
    "build": "astro build"
  },
  "dependencies": {
    "@astrojs/starlight": "^0.20.0",
    "astro": "^4.4.0"
  }
}`,
		"astro.config.mjs": `import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
  integrations: [
    starlight({
      title: '{{SITE_TITLE}}',
      social: {},
    }),
  ],
});`,
		"tsconfig.json": `{
  "extends": "astro/tsconfigs/strict"
}`,
		"src/":              "",
		"src/content/":      "",
		"src/content/docs/": "",
		"public/":           "",
	}

	return ctx.CreateDirectoryStructure(baseDir, structure)
}

// ParseDocStructure parses a Gherkin doc string into a directory structure
// Format:
//
//	docs/
//	  README.md
//	  guide/
//	    intro.md
func ParseDocStructure(docString string) map[string]string {
	structure := make(map[string]string)
	lines := strings.Split(docString, "\n")

	var currentPath []string
	currentIndent := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Calculate indentation level
		indent := len(line) - len(strings.TrimLeft(line, " "))
		name := strings.TrimSpace(line)

		// Adjust path based on indentation
		if indent > currentIndent {
			// Going deeper
		} else if indent < currentIndent {
			// Going back up - pop from path
			diff := (currentIndent - indent) / 2
			if diff >= len(currentPath) {
				currentPath = []string{}
			} else {
				currentPath = currentPath[:len(currentPath)-diff]
			}
		} else if len(currentPath) > 0 {
			// Same level - replace last element
			currentPath = currentPath[:len(currentPath)-1]
		}

		// Add current name to path
		currentPath = append(currentPath, name)
		fullPath := strings.Join(currentPath, "/")

		// If it ends with /, it's a directory
		if strings.HasSuffix(name, "/") {
			structure[fullPath] = ""
		} else {
			// It's a file - add with default content
			structure[fullPath] = fmt.Sprintf("# %s\n\nContent for %s", name, name)
		}

		currentIndent = indent
	}

	return structure
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// ReadFileContent reads and returns file content
func ReadFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// CountFilesInDir counts files in a directory (non-recursive)
func CountFilesInDir(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			count++
		}
	}
	return count, nil
}

// FindMarkdownFiles recursively finds all markdown files in a directory
func FindMarkdownFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".markdown")) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
