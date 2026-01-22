package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const (
	// StardocDir is the directory name under user's home directory
	StardocDir = ".stardoc"
	// SharedDir is the subdirectory for shared project files
	SharedDir = "shared"
	// RunsDir is the subdirectory for individual run workspaces
	RunsDir = "runs"
	// VersionFile stores the hash for cache invalidation
	VersionFile = ".stardoc-version"
	// LockFile prevents concurrent installs
	LockFile = ".lock"
)

// Manager handles shared project directory operations
type Manager struct {
	homeDir string
}

// NewManager creates a new shared project manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	return &Manager{
		homeDir: homeDir,
	}, nil
}

// GetStardocDir returns the path to ~/.stardoc/
func (m *Manager) GetStardocDir() string {
	return filepath.Join(m.homeDir, StardocDir)
}

// GetSharedDir returns the path to ~/.stardoc/shared/
func (m *Manager) GetSharedDir() string {
	return filepath.Join(m.GetStardocDir(), SharedDir)
}

// GetRunsDir returns the path to ~/.stardoc/runs/
func (m *Manager) GetRunsDir() string {
	return filepath.Join(m.GetStardocDir(), RunsDir)
}

// GetVersionFilePath returns the path to ~/.stardoc/shared/.stardoc-version
func (m *Manager) GetVersionFilePath() string {
	return filepath.Join(m.GetSharedDir(), VersionFile)
}

// GetLockFilePath returns the path to ~/.stardoc/shared/.lock
func (m *Manager) GetLockFilePath() string {
	return filepath.Join(m.GetSharedDir(), LockFile)
}

// GenerateRunID creates a unique identifier for a new run
func (m *Manager) GenerateRunID() string {
	return uuid.New().String()
}

// GetRunDir returns the path to a specific run directory
func (m *Manager) GetRunDir(runID string) string {
	return filepath.Join(m.GetRunsDir(), runID)
}

// EnsureDirectories creates the necessary directory structure
func (m *Manager) EnsureDirectories() error {
	dirs := []string{
		m.GetStardocDir(),
		m.GetSharedDir(),
		m.GetRunsDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// IsSharedProjectCurrent checks if the shared project exists and matches the current version
func (m *Manager) IsSharedProjectCurrent(expectedHash string) (bool, error) {
	// Check if node_modules exists
	nodeModules := filepath.Join(m.GetSharedDir(), "node_modules")
	if _, err := os.Stat(nodeModules); os.IsNotExist(err) {
		return false, nil
	}

	// Check version file
	versionFile := m.GetVersionFilePath()
	data, err := os.ReadFile(versionFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to read version file: %w", err)
	}

	currentHash := string(data)
	return currentHash == expectedHash, nil
}

// SaveVersion saves the current version hash to disk
func (m *Manager) SaveVersion(hash string) error {
	versionFile := m.GetVersionFilePath()
	if err := os.WriteFile(versionFile, []byte(hash), 0644); err != nil {
		return fmt.Errorf("failed to write version file: %w", err)
	}
	return nil
}

// ComputeHash computes a SHA-256 hash from the given content
func ComputeHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// CleanupOldRuns removes run directories older than the specified duration
func (m *Manager) CleanupOldRuns(maxAge time.Duration) error {
	runsDir := m.GetRunsDir()

	entries, err := os.ReadDir(runsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read runs directory: %w", err)
	}

	now := time.Now()
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		runPath := filepath.Join(runsDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		age := now.Sub(info.ModTime())
		if age > maxAge {
			os.RemoveAll(runPath)
		}
	}

	return nil
}

// AcquireLock creates a lock file to prevent concurrent installs
func (m *Manager) AcquireLock() error {
	lockFile := m.GetLockFilePath()

	// Try to create lock file exclusively
	f, err := os.OpenFile(lockFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("another installation is in progress")
		}
		return fmt.Errorf("failed to create lock file: %w", err)
	}
	f.Close()

	return nil
}

// ReleaseLock removes the lock file
func (m *Manager) ReleaseLock() error {
	lockFile := m.GetLockFilePath()
	if err := os.Remove(lockFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove lock file: %w", err)
	}
	return nil
}
