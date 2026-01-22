package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MarkdownFile represents a discovered markdown file
type MarkdownFile struct {
	Path     string // Path relative to source directory
	FullPath string // Absolute path
}

// Scanner discovers markdown files in a directory
type Scanner struct {
	sourceDir string
	files     []MarkdownFile
}

// New creates a new scanner for the given source directory
func New(sourceDir string) *Scanner {
	return &Scanner{
		sourceDir: sourceDir,
		files:     make([]MarkdownFile, 0),
	}
}

// Scan discovers all markdown files in the source directory
func (s *Scanner) Scan() ([]MarkdownFile, error) {
	err := filepath.WalkDir(s.sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Log warning but continue scanning
			fmt.Fprintf(os.Stderr, "Warning: cannot access %s: %v\n", path, err)
			return nil
		}

		// Skip hidden files and directories
		if strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		// Skip common exclude patterns
		if d.IsDir() && s.shouldSkipDir(d.Name()) {
			return fs.SkipDir
		}

		// Check if file is a markdown file
		if !d.IsDir() && s.isMarkdownFile(d.Name()) {
			relPath, err := filepath.Rel(s.sourceDir, path)
			if err != nil {
				return err
			}

			s.files = append(s.files, MarkdownFile{
				Path:     relPath,
				FullPath: path,
			})
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	// Sort files alphabetically for consistent ordering
	sort.Slice(s.files, func(i, j int) bool {
		return s.files[i].Path < s.files[j].Path
	})

	return s.files, nil
}

// GetFiles returns the discovered files (after Scan has been called)
func (s *Scanner) GetFiles() []MarkdownFile {
	return s.files
}

// Count returns the number of discovered files
func (s *Scanner) Count() int {
	return len(s.files)
}

// isMarkdownFile checks if a filename has a markdown extension
func (s *Scanner) isMarkdownFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".md" || ext == ".markdown" || ext == ".mdown" || ext == ".mkd"
}

// shouldSkipDir checks if a directory should be skipped
func (s *Scanner) shouldSkipDir(dirname string) bool {
	skipDirs := []string{
		"node_modules",
		"dist",
		"build",
		".obsidian",
		".vscode",
		".idea",
		"vendor",
	}

	for _, skip := range skipDirs {
		if dirname == skip {
			return true
		}
	}

	return false
}
