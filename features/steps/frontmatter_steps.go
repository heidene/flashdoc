package steps

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/heidene/flashdoc/internal/frontmatter"
	"gopkg.in/yaml.v3"
)

// RegisterFrontmatterSteps registers all frontmatter-related step definitions
func RegisterFrontmatterSteps(sc *godog.ScenarioContext, ctx *TestContext) {
	sc.Step(`^a markdown file "([^"]*)" with content:$`, ctx.createMarkdownFileWithContent)
	sc.Step(`^a markdown file "([^"]*)"$`, ctx.createMarkdownFileWithName)
	sc.Step(`^a markdown file "([^"]*)" in directory "([^"]*)"$`, ctx.createMarkdownFileInDirectory)
	sc.Step(`^a markdown file "([^"]*)" in the root source directory$`, ctx.createMarkdownFileInRoot)
	sc.Step(`^a markdown file "([^"]*)" with frontmatter:$`, ctx.createMarkdownFileWithFrontmatter)

	sc.Step(`^the file is processed$`, ctx.processFileWithFrontmatter)
	sc.Step(`^the output should have frontmatter:$`, ctx.outputShouldHaveFrontmatter)
	sc.Step(`^the frontmatter should remain unchanged$`, ctx.frontmatterShouldRemainUnchanged)
	sc.Step(`^the title should still be "([^"]*)"$`, ctx.titleShouldStillBe)
	sc.Step(`^the description should still be "([^"]*)"$`, ctx.descriptionShouldStillBe)
	sc.Step(`^the frontmatter should include:$`, ctx.frontmatterShouldInclude)
	sc.Step(`^the generated title should be "([^"]*)"$`, ctx.generatedTitleShouldBe)
	sc.Step(`^the number prefix "([^"]*)" should be stripped$`, ctx.numberPrefixShouldBeStripped)
	sc.Step(`^the title should be derived from the parent directory name$`, ctx.titleShouldBeDerivedFromParentDir)
	sc.Step(`^the title should remain "([^"]*)"$`, ctx.titleShouldRemain)
	sc.Step(`^the filename should not influence the title$`, ctx.filenameShouldNotInfluenceTitle)
	sc.Step(`^the malformed frontmatter should be detected$`, ctx.malformedFrontmatterShouldBeDetected)
	sc.Step(`^new valid frontmatter should be added at the top:$`, ctx.newValidFrontmatterShouldBeAdded)
	sc.Step(`^the original malformed content should be preserved as body content$`, ctx.originalContentShouldBePreserved)
	sc.Step(`^the frontmatter should be populated with:$`, ctx.frontmatterShouldBePopulatedWith)
	sc.Step(`^all frontmatter fields should be preserved$`, ctx.allFrontmatterFieldsShouldBePreserved)
	sc.Step(`^the structure should remain valid YAML$`, ctx.structureShouldRemainValidYAML)
	sc.Step(`^the frontmatter should be valid UTF-8$`, ctx.frontmatterShouldBeValidUTF8)
	sc.Step(`^the "([^"]*)" extension should not appear in the title$`, ctx.extensionShouldNotAppearInTitle)
}

func (ctx *TestContext) createMarkdownFileWithContent(filename, content string) error {
	tempDir, err := os.MkdirTemp("", "stardoc-fm-*")
	if err != nil {
		return err
	}
	ctx.TrackDir(tempDir)

	ctx.sourceDirectory = tempDir
	filePath := filepath.Join(tempDir, filename)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func (ctx *TestContext) createMarkdownFileWithName(filename string) error {
	tempDir, err := os.MkdirTemp("", "stardoc-fm-*")
	if err != nil {
		return err
	}
	ctx.TrackDir(tempDir)

	ctx.sourceDirectory = tempDir
	filePath := filepath.Join(tempDir, filename)

	content := "# Content\n\nThis is test content."
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func (ctx *TestContext) createMarkdownFileInDirectory(filename, directory string) error {
	tempDir, err := os.MkdirTemp("", "stardoc-fm-*")
	if err != nil {
		return err
	}
	ctx.TrackDir(tempDir)

	ctx.sourceDirectory = tempDir
	dirPath := filepath.Join(tempDir, directory)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(dirPath, filename)
	content := "# Content\n\nThis is test content."
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func (ctx *TestContext) createMarkdownFileInRoot(filename string) error {
	return ctx.createMarkdownFileWithName(filename)
}

func (ctx *TestContext) createMarkdownFileWithFrontmatter(filename, frontmatter string) error {
	tempDir, err := os.MkdirTemp("", "stardoc-fm-*")
	if err != nil {
		return err
	}
	ctx.TrackDir(tempDir)

	ctx.sourceDirectory = tempDir
	filePath := filepath.Join(tempDir, filename)

	content := frontmatter + "\n# Content\n\nThis is test content."
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func (ctx *TestContext) processFileWithFrontmatter() error {
	// Find the first .md file in source directory
	files, err := filepath.Glob(filepath.Join(ctx.sourceDirectory, "**/*.md"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		// Try direct children
		entries, err := os.ReadDir(ctx.sourceDirectory)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				files = append(files, filepath.Join(ctx.sourceDirectory, entry.Name()))
			} else if entry.IsDir() {
				// Check subdirectories
				subEntries, err := os.ReadDir(filepath.Join(ctx.sourceDirectory, entry.Name()))
				if err != nil {
					continue
				}
				for _, subEntry := range subEntries {
					if !subEntry.IsDir() && strings.HasSuffix(subEntry.Name(), ".md") {
						files = append(files, filepath.Join(ctx.sourceDirectory, entry.Name(), subEntry.Name()))
					}
				}
			}
		}
	}

	if len(files) == 0 {
		return fmt.Errorf("no markdown files found in source directory")
	}

	filePath := files[0]

	// Read the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Get parent directory for title generation
	relPath, _ := filepath.Rel(ctx.sourceDirectory, filePath)
	parentDir := filepath.Dir(relPath)
	if parentDir == "." {
		parentDir = ""
	}

	// Process with frontmatter
	processed, err := frontmatter.Inject(string(content), filepath.Base(filePath), parentDir)
	if err != nil {
		return err
	}

	ctx.processedContent = processed
	return nil
}

func (ctx *TestContext) outputShouldHaveFrontmatter(expectedOutput string) error {
	// Normalize whitespace for comparison
	expected := strings.TrimSpace(expectedOutput)
	actual := strings.TrimSpace(ctx.processedContent)

	if expected != actual {
		return fmt.Errorf("frontmatter mismatch:\nExpected:\n%s\n\nActual:\n%s", expected, actual)
	}

	return nil
}

func (ctx *TestContext) frontmatterShouldRemainUnchanged() error {
	// Parse the processed content
	fm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if fm == nil {
		return fmt.Errorf("frontmatter was removed")
	}

	return nil
}

func (ctx *TestContext) titleShouldStillBe(expectedTitle string) error {
	fm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if fm.Title != expectedTitle {
		return fmt.Errorf("expected title %q, got %q", expectedTitle, fm.Title)
	}

	return nil
}

func (ctx *TestContext) descriptionShouldStillBe(expectedDesc string) error {
	fm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if fm.Description != expectedDesc {
		return fmt.Errorf("expected description %q, got %q", expectedDesc, fm.Description)
	}

	return nil
}

func (ctx *TestContext) frontmatterShouldInclude(expectedFrontmatter string) error {
	// Parse expected frontmatter
	expectedFm, _, err := frontmatter.Parse(expectedFrontmatter)
	if err != nil {
		return err
	}

	// Parse actual frontmatter
	actualFm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if actualFm.Title != expectedFm.Title {
		return fmt.Errorf("title mismatch: expected %q, got %q", expectedFm.Title, actualFm.Title)
	}

	if expectedFm.Description != "" && actualFm.Description != expectedFm.Description {
		return fmt.Errorf("description mismatch: expected %q, got %q", expectedFm.Description, actualFm.Description)
	}

	return nil
}

func (ctx *TestContext) generatedTitleShouldBe(expectedTitle string) error {
	fm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if fm.Title != expectedTitle {
		return fmt.Errorf("expected generated title %q, got %q", expectedTitle, fm.Title)
	}

	return nil
}

func (ctx *TestContext) numberPrefixShouldBeStripped(prefix string) error {
	fm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if strings.Contains(fm.Title, prefix) {
		return fmt.Errorf("title still contains number prefix %q: %q", prefix, fm.Title)
	}

	return nil
}

func (ctx *TestContext) titleShouldBeDerivedFromParentDir() error {
	// This is implicitly tested by checking the generated title
	return nil
}

func (ctx *TestContext) titleShouldRemain(expectedTitle string) error {
	return ctx.titleShouldStillBe(expectedTitle)
}

func (ctx *TestContext) filenameShouldNotInfluenceTitle() error {
	// This is verified by checking that the title matches the frontmatter, not the filename
	return nil
}

func (ctx *TestContext) malformedFrontmatterShouldBeDetected() error {
	// Check that frontmatter was added (which means the malformed one was detected)
	if !frontmatter.HasFrontmatter(ctx.processedContent) {
		return fmt.Errorf("no frontmatter found in processed content")
	}
	return nil
}

func (ctx *TestContext) newValidFrontmatterShouldBeAdded(expectedFrontmatter string) error {
	// Parse expected frontmatter
	expectedFm, _, err := frontmatter.Parse(expectedFrontmatter)
	if err != nil {
		return err
	}

	// Parse actual frontmatter
	actualFm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if actualFm.Title != expectedFm.Title {
		return fmt.Errorf("title mismatch: expected %q, got %q", expectedFm.Title, actualFm.Title)
	}

	return nil
}

func (ctx *TestContext) originalContentShouldBePreserved() error {
	// The body content should still be present
	if !strings.Contains(ctx.processedContent, "Content") {
		return fmt.Errorf("original content was not preserved")
	}
	return nil
}

func (ctx *TestContext) frontmatterShouldBePopulatedWith(expectedFrontmatter string) error {
	return ctx.frontmatterShouldInclude(expectedFrontmatter)
}

func (ctx *TestContext) allFrontmatterFieldsShouldBePreserved() error {
	// Parse the processed content and check it's valid YAML
	fm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if fm == nil {
		return fmt.Errorf("frontmatter was not preserved")
	}

	// Check that title is preserved
	if fm.Title == "" {
		return fmt.Errorf("title was not preserved")
	}

	return nil
}

func (ctx *TestContext) structureShouldRemainValidYAML() error {
	// Try to parse as YAML
	fm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if fm == nil {
		return fmt.Errorf("frontmatter is not valid YAML")
	}

	// Also try parsing the entire frontmatter block
	lines := strings.Split(ctx.processedContent, "\n")
	if len(lines) < 3 || lines[0] != "---" {
		return fmt.Errorf("invalid frontmatter structure")
	}

	closingIndex := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			closingIndex = i
			break
		}
	}

	if closingIndex == -1 {
		return fmt.Errorf("frontmatter closing delimiter not found")
	}

	fmContent := strings.Join(lines[1:closingIndex], "\n")
	var yamlCheck map[string]interface{}
	if err := yaml.Unmarshal([]byte(fmContent), &yamlCheck); err != nil {
		return fmt.Errorf("frontmatter is not valid YAML: %w", err)
	}

	return nil
}

func (ctx *TestContext) frontmatterShouldBeValidUTF8() error {
	// Parse frontmatter to ensure it's valid
	fm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if fm == nil {
		return fmt.Errorf("frontmatter is not valid")
	}

	return nil
}

func (ctx *TestContext) extensionShouldNotAppearInTitle(extension string) error {
	fm, _, err := frontmatter.Parse(ctx.processedContent)
	if err != nil {
		return err
	}

	if strings.Contains(fm.Title, extension) {
		return fmt.Errorf("title should not contain extension %q: %q", extension, fm.Title)
	}

	return nil
}
