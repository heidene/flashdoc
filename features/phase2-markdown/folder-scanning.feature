Feature: Folder Scanning
  As a stardoc user
  I want the tool to discover all markdown files in my source directory
  So that they can be included in the generated documentation site

  Background:
    Given the stardoc CLI is available

  Scenario: Scan directory with markdown files
    Given a directory "./docs" with the following structure:
      """
      ./docs/
      ├── README.md
      ├── getting-started.md
      └── guides/
          └── installation.md
      """
    When I run "stardoc ./docs"
    Then the scanner should find 3 markdown files
    And the files should be:
      | path                        |
      | README.md                   |
      | getting-started.md          |
      | guides/installation.md      |

  Scenario: Scan nested directories
    Given a directory "./docs" with the following structure:
      """
      ./docs/
      ├── index.md
      ├── api/
      │   ├── overview.md
      │   └── endpoints/
      │       ├── users.md
      │       └── posts.md
      └── guides/
          └── tutorial.md
      """
    When I run "stardoc ./docs"
    Then the scanner should find 5 markdown files
    And the scanner should preserve the directory structure
    And "api/endpoints/users.md" should maintain its nested path

  Scenario: Handle empty directory
    Given an empty directory "./empty-docs"
    When I run "stardoc ./empty-docs"
    Then the scanner should find 0 markdown files
    And the CLI should display a warning "no markdown files found in ./empty-docs"
    And the CLI should exit with code 1

  Scenario: Ignore non-markdown files
    Given a directory "./mixed" with the following structure:
      """
      ./mixed/
      ├── doc.md
      ├── image.png
      ├── script.js
      └── data.json
      """
    When I run "stardoc ./mixed"
    Then the scanner should find 1 markdown file
    And the scanner should only include "doc.md"
    And "image.png", "script.js", and "data.json" should be ignored

  Scenario: Support various markdown extensions
    Given a directory "./docs" with files:
      """
      ./docs/
      ├── file1.md
      ├── file2.markdown
      ├── file3.mdown
      ├── file4.mkd
      └── file5.txt
      """
    When I run "stardoc ./docs"
    Then the scanner should find 4 markdown files
    And ".md", ".markdown", ".mdown", ".mkd" extensions should be included
    And ".txt" should be excluded

  Scenario: Ignore hidden files and directories
    Given a directory "./docs" with the following structure:
      """
      ./docs/
      ├── visible.md
      ├── .hidden.md
      └── .git/
          └── config.md
      """
    When I run "stardoc ./docs"
    Then the scanner should find 1 markdown file
    And the scanner should only include "visible.md"
    And ".hidden.md" should be ignored
    And files in ".git/" should be ignored

  Scenario: Ignore common exclude patterns
    Given a directory "./docs" with the following structure:
      """
      ./docs/
      ├── index.md
      ├── node_modules/
      │   └── package.md
      ├── .obsidian/
      │   └── workspace.md
      └── dist/
          └── build.md
      """
    When I run "stardoc ./docs"
    Then the scanner should find 1 markdown file
    And the scanner should only include "index.md"
    And "node_modules/", ".obsidian/", "dist/" should be ignored

  Scenario: Preserve file order for navigation
    Given a directory "./docs" with files:
      """
      ./docs/
      ├── 01-intro.md
      ├── 02-setup.md
      ├── 03-usage.md
      """
    When I run "stardoc ./docs"
    Then the files should be scanned in alphabetical order
    And the navigation should reflect this order

  Scenario: Handle symbolic links
    Given a directory "./docs" with a symbolic link "linked.md" pointing to "../other/real.md"
    When I run "stardoc ./docs"
    Then the scanner should follow the symbolic link
    And "linked.md" should be included in the scan
    And the content should be read from "../other/real.md"

  Scenario: Handle scan errors gracefully
    Given a directory "./docs" with a subdirectory "forbidden"
    And "forbidden" has no read permissions
    When I run "stardoc ./docs"
    Then the scanner should log a warning about "forbidden"
    And the scanner should continue scanning other accessible directories
    And other markdown files should still be included

  Scenario: Report scan summary
    Given a directory "./docs" with 15 markdown files in various subdirectories
    When I run "stardoc ./docs"
    Then the CLI should log "Found 15 markdown files"
    And the log should include a summary of scanned directories
