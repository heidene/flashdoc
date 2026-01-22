package pkgmanager

import (
	"fmt"
	"os/exec"
)

// PackageManager represents a Node.js package manager
type PackageManager string

const (
	Pnpm PackageManager = "pnpm"
	Bun  PackageManager = "bun"
	Npm  PackageManager = "npm"
)

// Detect finds the best available package manager
func Detect() (PackageManager, error) {
	// Check for pnpm first (fastest)
	if isAvailable("pnpm") {
		return Pnpm, nil
	}

	// Check for bun
	if isAvailable("bun") {
		return Bun, nil
	}

	// Check for npm (most common, usually installed with Node.js)
	if isAvailable("npm") {
		return Npm, nil
	}

	return "", fmt.Errorf("no package manager found (tried: pnpm, bun, npm)")
}

// isAvailable checks if a command is available in PATH
func isAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// InstallCommand returns the install command for the package manager
func (pm PackageManager) InstallCommand() []string {
	switch pm {
	case Pnpm:
		return []string{"pnpm", "install"}
	case Bun:
		return []string{"bun", "install"}
	case Npm:
		return []string{"npm", "install"}
	default:
		return []string{"npm", "install"}
	}
}

// BuildCommand returns the build command for the package manager
func (pm PackageManager) BuildCommand() []string {
	switch pm {
	case Pnpm:
		return []string{"pnpm", "run", "build"}
	case Bun:
		return []string{"bun", "run", "build"}
	case Npm:
		return []string{"npm", "run", "build"}
	default:
		return []string{"npm", "run", "build"}
	}
}

// String returns the string representation of the package manager
func (pm PackageManager) String() string {
	return string(pm)
}
