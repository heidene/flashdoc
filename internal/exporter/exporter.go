package exporter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Exporter handles exporting the built static site to a directory
type Exporter struct {
	distPath   string
	exportPath string
	output     io.Writer
}

// New creates a new Exporter
func New(distPath, exportPath string, output io.Writer) *Exporter {
	return &Exporter{
		distPath:   distPath,
		exportPath: exportPath,
		output:     output,
	}
}

// Export copies the built static site to the export directory
func (e *Exporter) Export() error {
	// Resolve export path to absolute
	absExportPath, err := filepath.Abs(e.exportPath)
	if err != nil {
		return fmt.Errorf("failed to resolve export path: %w", err)
	}

	// Check if export directory exists
	if _, err := os.Stat(absExportPath); err == nil {
		fmt.Fprintf(e.output, "⚠️  Warning: export directory already exists, overwriting...\n")
	}

	// Create export directory
	if err := os.MkdirAll(absExportPath, 0755); err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	fmt.Fprintf(e.output, "Copying files to %s...\n", e.exportPath)

	// Copy all files from dist to export directory
	fileCount := 0
	err = filepath.Walk(e.distPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from dist directory
		relPath, err := filepath.Rel(e.distPath, path)
		if err != nil {
			return err
		}

		// Skip the dist directory itself
		if relPath == "." {
			return nil
		}

		// Destination path
		destPath := filepath.Join(absExportPath, relPath)

		// If it's a directory, create it
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		if err := copyFile(path, destPath); err != nil {
			return fmt.Errorf("failed to copy %s: %w", relPath, err)
		}

		fileCount++
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy files: %w", err)
	}

	fmt.Fprintf(e.output, "Exported %d files\n", fileCount)
	fmt.Fprintf(e.output, "✅ Exported to %s\n", e.exportPath)

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}
