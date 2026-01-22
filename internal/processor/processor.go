package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicovandenhove/stardoc/internal/frontmatter"
	"github.com/nicovandenhove/stardoc/internal/scanner"
)

// Processor handles markdown file processing and copying
type Processor struct {
	sourceDir   string
	targetDir   string
	filescopied int
}

// New creates a new processor
func New(sourceDir, targetDir string) *Processor {
	return &Processor{
		sourceDir: sourceDir,
		targetDir: targetDir,
	}
}

// Process scans and copies all markdown files with frontmatter injection
func (p *Processor) Process() error {
	// Scan for markdown files
	s := scanner.New(p.sourceDir)
	files, err := s.Scan()
	if err != nil {
		return fmt.Errorf("failed to scan source directory: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no markdown files found in %s", p.sourceDir)
	}

	fmt.Printf("Found %d markdown files\n", len(files))
	fmt.Printf("Processing %d files...\n", len(files))

	// Process each file
	for _, file := range files {
		if err := p.processFile(file); err != nil {
			return fmt.Errorf("failed to copy %s: %w", file.Path, err)
		}
		p.filescopied++
	}

	fmt.Printf("Copied %d files successfully\n", p.filescopied)

	return nil
}

// processFile processes a single markdown file
func (p *Processor) processFile(file scanner.MarkdownFile) error {
	// Read source file
	content, err := os.ReadFile(file.FullPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Determine target filename (rename README.md -> index.md)
	targetFilename := filepath.Base(file.Path)
	if strings.ToUpper(targetFilename) == "README.MD" {
		targetFilename = "index.md"
	}

	// Build target path
	targetDir := filepath.Join(p.targetDir, filepath.Dir(file.Path))
	targetPath := filepath.Join(targetDir, targetFilename)

	// Get parent directory for title generation
	parentDir := filepath.Dir(file.Path)
	if parentDir == "." {
		parentDir = ""
	}

	// Inject frontmatter
	processed, err := frontmatter.Inject(string(content), filepath.Base(file.Path), parentDir)
	if err != nil {
		return fmt.Errorf("failed to inject frontmatter: %w", err)
	}

	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Write processed file
	if err := os.WriteFile(targetPath, []byte(processed), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// GetCopiedCount returns the number of files copied
func (p *Processor) GetCopiedCount() int {
	return p.filescopied
}
