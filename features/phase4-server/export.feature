Feature: Static Site Export
  As a stardoc user
  I want to export a static build of the documentation site
  So that I can host it on a web server or share it with others

  Background:
    Given the stardoc CLI is available
    And a source directory "./test-docs" exists with markdown files

  Scenario: Export to default directory
    When I run "stardoc ./test-docs --export"
    Then the CLI should build the static site
    And the static files should be exported to "./export-doc"
    And the export directory should contain "index.html"
    And the export directory should contain "_astro/" subdirectory
    And the CLI should display "✅ Exported to ./export-doc"
    And the CLI should exit with code 0

  Scenario: Export to custom directory
    When I run "stardoc ./test-docs --export ../docs"
    Then the CLI should build the static site
    And the static files should be exported to "../docs"
    And the export directory should contain "index.html"
    And the CLI should display "✅ Exported to ../docs"
    And the CLI should exit with code 0

  Scenario: Export to absolute path
    Given a temporary directory "/tmp/stardoc-export-test"
    When I run "stardoc ./test-docs --export /tmp/stardoc-export-test"
    Then the static files should be exported to "/tmp/stardoc-export-test"
    And the export directory should contain "index.html"
    And the CLI should exit with code 0

  Scenario: Export to existing directory
    Given a directory "./existing-export" exists
    And "./existing-export" contains a file "old-file.txt"
    When I run "stardoc ./test-docs --export ./existing-export"
    Then the CLI should display a warning about overwriting existing files
    And the static files should be exported to "./existing-export"
    And the old file should be replaced with new static files
    And the CLI should exit with code 0

  Scenario: Export to directory that cannot be created
    Given the parent directory "/forbidden-path" does not exist
    And I do not have write permissions for "/forbidden-path"
    When I run "stardoc ./test-docs --export /forbidden-path/export"
    Then the CLI should display an error "failed to create export directory"
    And the CLI should exit with code 1

  Scenario: Export includes all documentation pages
    Given a directory "./docs" with the following structure:
      """
      ./docs/
      ├── README.md
      ├── guide.md
      └── api/
          └── reference.md
      """
    When I run "stardoc ./docs --export ./output"
    Then the export directory should contain "index.html"
    And the export directory should contain "guide/index.html"
    And the export directory should contain "api/reference/index.html"
    And each page should be accessible as a static HTML file

  Scenario: Export with custom title
    When I run "stardoc ./test-docs --title 'My Project' --export ./docs-output"
    Then the exported site should use the title "My Project"
    And the index.html should contain "My Project" in the title tag
    And the CLI should exit with code 0

  Scenario: Export build failure handling
    Given the Astro build process will fail
    When I run "stardoc ./test-docs --export ./output"
    Then the CLI should display an error "build failed"
    And the CLI should display the build error output
    And the export directory should not be created
    And the CLI should exit with code 1

  Scenario: Export with no-open flag
    When I run "stardoc ./test-docs --export --no-open"
    Then the CLI should not start the dev server
    And the CLI should not open a browser
    And the CLI should only build and export
    And the CLI should exit with code 0

  Scenario: Export progress indication
    When I run "stardoc ./test-docs --export ./output"
    Then the CLI should display "Building static site..."
    And the CLI should display a progress indicator during build
    And the CLI should display "Copying files to ./output..."
    And the CLI should display the total number of files exported
    And the CLI should exit with code 0

  Scenario: Export creates necessary parent directories
    When I run "stardoc ./test-docs --export ./nested/path/to/docs"
    Then the CLI should create all parent directories
    And the static files should be exported to "./nested/path/to/docs"
    And the CLI should exit with code 0

  Scenario: Export without starting dev server
    When I run "stardoc ./test-docs --export"
    Then the CLI should not start the Astro dev server
    And the CLI should not display "Press Ctrl+C to stop"
    And the CLI should build the static site and exit
    And the CLI should exit with code 0

  Scenario: Export cleanup on failure
    Given the export will fail during file copy
    When I run "stardoc ./test-docs --export ./output"
    Then the CLI should attempt to clean up partial files
    And the CLI should display an error message
    And the CLI should exit with code 1

  Scenario: Export validation
    When I run "stardoc ./test-docs --export ./validated-output"
    Then the export should contain valid HTML files
    And the export should contain all CSS and JavaScript assets
    And the export should contain all referenced images
    And the navigation structure should be preserved
    And the CLI should exit with code 0

  Scenario: Export respects ignored files
    Given a directory "./docs" with files:
      """
      ./docs/
      ├── public.md
      ├── .hidden.md
      └── node_modules/
          └── package.md
      """
    When I run "stardoc ./docs --export ./output"
    Then the export should include "public.md" content
    And the export should not include ".hidden.md" content
    And the export should not include "node_modules" content
    And the CLI should exit with code 0

  Scenario: Export with verbose output
    When I run "stardoc ./test-docs --export --verbose"
    Then the CLI should display detailed build logs
    And the CLI should display each file being copied
    And the CLI should display build statistics
    And the CLI should exit with code 0

  Scenario: Multiple exports to same directory
    When I run "stardoc ./test-docs --export ./output"
    And I modify a markdown file in "./test-docs"
    And I run "stardoc ./test-docs --export ./output" again
    Then the export directory should be updated with new content
    And the old content should be replaced
    And the CLI should exit with code 0

  Scenario: Export displays build duration
    When I run "stardoc ./test-docs --export"
    Then the CLI should display the build duration
    And the duration should be in a human-readable format
    And the CLI should exit with code 0
