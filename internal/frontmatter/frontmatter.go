package frontmatter

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Frontmatter represents parsed frontmatter data
type Frontmatter struct {
	Title       string                 `yaml:"title,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	Other       map[string]interface{} `yaml:",inline"`
}

// Parse extracts frontmatter from markdown content
func Parse(content string) (*Frontmatter, string, error) {
	// Check if content starts with frontmatter delimiters
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return nil, content, nil
	}

	// Find the closing delimiter
	lines := strings.Split(content, "\n")
	closingIndex := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			closingIndex = i
			break
		}
	}

	if closingIndex == -1 {
		// Malformed frontmatter - no closing delimiter
		return nil, content, nil
	}

	// Extract frontmatter content
	fmContent := strings.Join(lines[1:closingIndex], "\n")
	bodyContent := strings.Join(lines[closingIndex+1:], "\n")

	// Parse YAML
	var fm Frontmatter
	if err := yaml.Unmarshal([]byte(fmContent), &fm); err != nil {
		// Malformed YAML - return nil frontmatter
		return nil, content, nil
	}

	return &fm, bodyContent, nil
}

// Inject adds or updates frontmatter in markdown content
func Inject(content, filename, parentDir string) (string, error) {
	// Parse existing frontmatter
	fm, body, err := Parse(content)
	if err != nil {
		return "", err
	}

	// Create new frontmatter if none exists
	if fm == nil {
		fm = &Frontmatter{
			Other: make(map[string]interface{}),
		}
	}

	// Add title if missing
	if fm.Title == "" {
		fm.Title = GenerateTitle(filename, parentDir)
	}

	// Serialize frontmatter
	fmBytes, err := yaml.Marshal(fm)
	if err != nil {
		return "", fmt.Errorf("failed to marshal frontmatter: %w", err)
	}

	// Combine frontmatter and body
	result := fmt.Sprintf("---\n%s---\n%s", string(fmBytes), body)
	return result, nil
}

// GenerateTitle creates a title from a filename
func GenerateTitle(filename, parentDir string) string {
	// Remove extension
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Handle index/README files - use parent directory name
	if name == "index" || name == "README" {
		if parentDir == "" || parentDir == "." {
			return "Home"
		}
		name = filepath.Base(parentDir)
	}

	// Remove numbered prefixes (e.g., "01-", "002-")
	re := regexp.MustCompile(`^\d+-`)
	name = re.ReplaceAllString(name, "")

	// Replace separators with spaces
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Capitalize first letter of each word
	words := strings.Fields(name)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	title := strings.Join(words, " ")
	if title == "" {
		return "Untitled"
	}

	return title
}

// HasFrontmatter checks if content has frontmatter
func HasFrontmatter(content string) bool {
	fm, _, _ := Parse(content)
	return fm != nil
}
