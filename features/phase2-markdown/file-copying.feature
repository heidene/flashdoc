Feature: File Copying
  As a stardoc user
  I want markdown files to be copied to the correct Starlight structure
  So that the documentation site can be generated properly

  Background:
    Given the stardoc CLI is available
    And a source directory "./docs" exists
    And a temp workspace exists at "/tmp/stardoc-abc123"

  Scenario: Copy files to Starlight content directory
    Given the source directory contains:
      """
      ./docs/
      ├── README.md
      └── guide.md
      """
    When files are processed and copied
    Then the temp workspace should contain:
      """
      /tmp/stardoc-abc123/src/content/docs/
      ├── index.md       (from README.md)
      └── guide.md
      """

  Scenario: Preserve directory structure
    Given the source directory contains:
      """
      ./docs/
      ├── intro.md
      └── api/
          ├── overview.md
          └── endpoints/
              └── users.md
      """
    When files are processed and copied
    Then the temp workspace should contain:
      """
      /tmp/stardoc-abc123/src/content/docs/
      ├── intro.md
      └── api/
          ├── overview.md
          └── endpoints/
              └── users.md
      """

  Scenario: Rename README.md to index.md
    Given a source file "README.md"
    When the file is copied
    Then it should be copied as "index.md" in the target directory
    And the frontmatter title should be derived appropriately

  Scenario: Rename nested README files
    Given the source directory contains:
      """
      ./docs/
      ├── README.md
      └── guides/
          └── README.md
      """
    When files are processed and copied
    Then the temp workspace should contain:
      """
      /tmp/stardoc-abc123/src/content/docs/
      ├── index.md
      └── guides/
          └── index.md
      """

  Scenario: Create necessary parent directories
    Given a source file "deep/nested/path/file.md"
    When the file is copied
    Then the directories "deep/nested/path/" should be created in the target
    And the file should be copied to "deep/nested/path/file.md"

  Scenario: Handle file copy errors
    Given a source file "document.md"
    And the target directory is not writable
    When the file copy is attempted
    Then an error should be logged "failed to copy document.md"
    And the error should include the underlying system error
    And the CLI should exit with code 1

  Scenario: Process files before copying
    Given a source file "guide.md" without frontmatter
    When the file is copied
    Then frontmatter should be injected first
    And then the processed content should be written to the target
    And the original source file should remain unchanged

  Scenario: Preserve file timestamps
    Given a source file "doc.md" with modification time "2024-01-15 10:30:00"
    When the file is copied
    Then the target file should have the same modification time

  Scenario: Handle special characters in filenames
    Given a source file "FAQ's & Tips.md"
    When the file is copied
    Then the file should be copied as "FAQ's & Tips.md"
    And the special characters should be preserved

  Scenario: Copy files in correct order
    Given the source directory contains 100 markdown files
    When files are copied
    Then all 100 files should be copied successfully
    And a progress indicator should show "Copied 100/100 files"

  Scenario: Skip non-markdown files
    Given the source directory contains:
      """
      ./docs/
      ├── doc.md
      ├── image.png
      └── data.json
      """
    When files are copied
    Then only "doc.md" should be copied to the target
    And "image.png" and "data.json" should be ignored

  Scenario: Handle duplicate filenames in different directories
    Given the source directory contains:
      """
      ./docs/
      ├── api/
      │   └── overview.md
      └── guides/
          └── overview.md
      """
    When files are copied
    Then both "overview.md" files should be copied
    And they should maintain their separate directory paths:
      """
      /tmp/stardoc-abc123/src/content/docs/
      ├── api/
      │   └── overview.md
      └── guides/
          └── overview.md
      """

  Scenario: Report copy statistics
    Given the source directory contains 42 markdown files
    When files are copied
    Then the CLI should log "Processing 42 files..."
    And the CLI should log "Copied 42 files successfully"

  Scenario: Handle unicode filenames
    Given a source file "文档.md"
    When the file is copied
    Then the file should be copied as "文档.md"
    And the unicode characters should be preserved correctly

  Scenario: Atomic copy prevents partial files
    Given a source file "large.md" with 10MB of content
    When the file is being copied
    And the process is interrupted mid-copy
    Then the temp workspace should not contain a partial "large.md"

  Scenario: Interrupted copy cleanup
    Given a source file "large.md" with 10MB of content
    When the file is being copied
    And the process is interrupted mid-copy
    And cleanup runs on shutdown
    Then any partial files should be removed
