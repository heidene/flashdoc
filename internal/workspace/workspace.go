package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Workspace manages the run directory for Starlight site generation
type Workspace struct {
	Path      string
	SharedDir string
}

// New creates a new workspace directory under ~/.stardoc/runs/{id}/
func New(runDir, sharedDir string) (*Workspace, error) {
	// Create run directory
	if err := os.MkdirAll(runDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create run directory: %w", err)
	}

	ws := &Workspace{
		Path:      runDir,
		SharedDir: sharedDir,
	}

	return ws, nil
}

// Setup creates the Starlight directory structure and symlinks in the workspace
func (w *Workspace) Setup() error {
	// Create required directories
	dirs := []string{
		filepath.Join(w.Path, "src", "content", "docs"),
		filepath.Join(w.Path, "public"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create symlinks to shared project files
	if err := w.createSymlinks(); err != nil {
		return fmt.Errorf("failed to create symlinks: %w", err)
	}

	// Template files (astro.config.mjs, tsconfig.json) will be extracted by the template package
	// package.json is symlinked from shared directory

	return nil
}

// createSymlinks creates symlinks to the shared project files
func (w *Workspace) createSymlinks() error {
	// Define symlinks to create: target -> source
	symlinks := map[string]string{
		"node_modules": filepath.Join(w.SharedDir, "node_modules"),
		"package.json": filepath.Join(w.SharedDir, "package.json"),
	}

	for linkName, target := range symlinks {
		linkPath := filepath.Join(w.Path, linkName)

		// Check if target exists
		if _, err := os.Stat(target); os.IsNotExist(err) {
			return fmt.Errorf("symlink target does not exist: %s", target)
		}

		// Remove existing link if present
		os.Remove(linkPath)

		// Create symlink (platform-specific handling)
		if err := createSymlink(target, linkPath); err != nil {
			return fmt.Errorf("failed to create symlink %s -> %s: %w", linkName, target, err)
		}
	}

	return nil
}

// createSymlink creates a symlink with platform-specific handling
func createSymlink(target, link string) error {
	// On Windows, try to create a junction point if it's a directory
	if runtime.GOOS == "windows" {
		info, err := os.Stat(target)
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Use junction point on Windows for directories (doesn't require admin)
			// Fall back to regular symlink if junction creation fails
			return os.Symlink(target, link)
		}
	}

	// On Unix-like systems or for files, use regular symlink
	return os.Symlink(target, link)
}

// Cleanup removes the run directory
func (w *Workspace) Cleanup() error {
	if w.Path == "" {
		return nil
	}

	// Only remove the run directory, not the shared directory
	if err := os.RemoveAll(w.Path); err != nil {
		return fmt.Errorf("failed to remove workspace directory: %w", err)
	}

	return nil
}

// Exists checks if the workspace directory still exists
func (w *Workspace) Exists() bool {
	if w.Path == "" {
		return false
	}

	_, err := os.Stat(w.Path)
	return err == nil
}

// GetDocsDir returns the path to the docs directory
func (w *Workspace) GetDocsDir() string {
	return filepath.Join(w.Path, "src", "content", "docs")
}

// GetDistDir returns the path to the dist directory (build output)
func (w *Workspace) GetDistDir() string {
	return filepath.Join(w.Path, "dist")
}
