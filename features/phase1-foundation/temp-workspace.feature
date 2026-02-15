Feature: Temporary Workspace Management
  As a stardoc user
  I want the tool to create a temporary workspace for the Starlight site
  So that my file system remains clean and I don't need to manage build artifacts

  Background:
    Given the stardoc CLI is available
    And a source directory "./test-docs" exists with markdown files

  Scenario: Create temporary directory on startup
    When I run "stardoc ./test-docs"
    Then a temporary directory should be created in the system temp location
    And the temp directory name should start with "stardoc-"
    And the temp directory name should include a unique identifier
    And the temp directory should have write permissions

  Scenario: Temp directory contains Starlight structure
    When I run "stardoc ./test-docs"
    Then the temp directory should contain a "src/content/docs" subdirectory
    And the temp directory should contain a "src/content.config.ts" file
    And the temp directory should contain a "public" subdirectory
    And the temp directory should contain a "package.json" file
    And the temp directory should contain an "astro.config.mjs" file

  Scenario: Temp directory location is deterministic
    When I run "stardoc ./test-docs"
    Then the temp directory should be created under the OS temp directory
    And the temp directory path should be logged to the terminal
    And the log message should be "Workspace: /tmp/stardoc-XXXXX"

  Scenario: Reuse existing temp directory if still valid
    Given stardoc has created a temp directory for "./test-docs"
    And the temp directory still exists
    When I run "stardoc ./test-docs" again
    Then a new temporary directory should be created
    And the new directory should have a different unique identifier

  Scenario: Handle temp directory creation failure
    Given the system temp directory is not writable
    When I run "stardoc ./test-docs"
    Then the CLI should display an error "failed to create workspace"
    And the CLI should exit with code 1

  Scenario: Temp directory isolation
    Given I run "stardoc ./docs-a" in one terminal
    When I run "stardoc ./docs-b" in another terminal
    Then each instance should have its own isolated temp directory
    And the directories should not interfere with each other

  Scenario: Temp directory cleanup on successful exit
    Given I run "stardoc ./test-docs"
    And the dev server starts successfully
    When I press Ctrl+C to stop the server
    Then the temp directory should be deleted
    And no stardoc-* directories should remain in the system temp

  Scenario: Temp directory cleanup on error during setup
    Given stardoc creates a temp directory for "./test-docs"
    When an error occurs during dependency installation
    Then the temp directory should be deleted before exit
    And no orphaned temp directories should remain
