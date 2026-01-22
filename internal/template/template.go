package template

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed starlight/* starlight/src/*
var starlightTemplate embed.FS

// Extract extracts the embedded Starlight template to the workspace
func Extract(workspacePath string) error {
	entries, err := starlightTemplate.ReadDir("starlight")
	if err != nil {
		return fmt.Errorf("failed to extract template: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		content, err := starlightTemplate.ReadFile(filepath.Join("starlight", filename))
		if err != nil {
			return fmt.Errorf("failed to extract template: %w", err)
		}

		targetPath := filepath.Join(workspacePath, filename)
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to extract template: %w", err)
		}
	}

	// Create required Starlight directory structure
	requiredDirs := []string{
		filepath.Join(workspacePath, "public"),
		filepath.Join(workspacePath, "src"),
		filepath.Join(workspacePath, "src", "content"),
		filepath.Join(workspacePath, "src", "content", "docs"),
	}

	for _, dir := range requiredDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to extract template: %w", err)
		}
	}

	return nil
}

// GenerateConfig replaces the {{SITE_TITLE}} placeholder in astro.config.mjs
func GenerateConfig(workspacePath, title string) error {
	configPath := filepath.Join(workspacePath, "astro.config.mjs")

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	// Replace the placeholder
	newContent := strings.ReplaceAll(string(content), "{{SITE_TITLE}}", title)

	if err := os.WriteFile(configPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to generate config: %w", err)
	}

	return nil
}

// GenerateTitle creates a title from a directory name
func GenerateTitle(dirPath string) string {
	// Get the base directory name
	baseName := filepath.Base(dirPath)

	// Handle absolute paths
	if baseName == "/" || baseName == "." {
		return "Documentation"
	}

	// Replace separators with spaces
	title := strings.ReplaceAll(baseName, "-", " ")
	title = strings.ReplaceAll(title, "_", " ")

	// Capitalize first letter of each word
	words := strings.Fields(title)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	result := strings.Join(words, " ")
	if result == "" {
		return "Documentation"
	}

	return result
}

// ExtractToShared extracts the embedded template to the shared directory
func ExtractToShared(sharedDir string) error {
	entries, err := starlightTemplate.ReadDir("starlight")
	if err != nil {
		return fmt.Errorf("failed to extract template to shared: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		content, err := starlightTemplate.ReadFile(filepath.Join("starlight", filename))
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", filename, err)
		}

		targetPath := filepath.Join(sharedDir, filename)
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write template file %s: %w", filename, err)
		}
	}

	return nil
}

// GetEmbeddedPackageHash returns a hash of the embedded package.json for cache invalidation
func GetEmbeddedPackageHash() (string, error) {
	content, err := starlightTemplate.ReadFile("starlight/package.json")
	if err != nil {
		return "", fmt.Errorf("failed to read embedded package.json: %w", err)
	}

	hash := sha256.Sum256(content)
	return hex.EncodeToString(hash[:]), nil
}

// ExtractConfigOnly extracts only the config files (not package.json) to the workspace
func ExtractConfigOnly(workspacePath string) error {
	// Files to extract (excluding package.json which is symlinked)
	configFiles := []string{
		"astro.config.mjs",
		"tsconfig.json",
		"src/content.config.ts",
	}

	for _, filename := range configFiles {
		content, err := starlightTemplate.ReadFile(filepath.Join("starlight", filename))
		if err != nil {
			// If file doesn't exist in template, skip it
			if os.IsNotExist(err) {
				continue
			}
			return fmt.Errorf("failed to read template file %s: %w", filename, err)
		}

		targetPath := filepath.Join(workspacePath, filename)

		// Create parent directory if needed
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", filename, err)
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write config file %s: %w", filename, err)
		}
	}

	return nil
}
