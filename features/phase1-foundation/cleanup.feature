Feature: Cleanup
  As a stardoc user
  I want the tool to clean up all temporary files and resources on exit
  So that my system doesn't accumulate orphaned files and processes

  Background:
    Given the stardoc CLI is available
    And a source directory "./test-docs" exists with markdown files

  Scenario: Normal cleanup on exit
    Given I run "stardoc ./test-docs"
    And the dev server starts successfully
    And a temp directory has been created at "/tmp/stardoc-abc123"
    When I press Ctrl+C to stop the server
    Then the temp directory "/tmp/stardoc-abc123" should be completely removed
    And no files should remain in "/tmp/stardoc-abc123"
    And the directory itself should not exist

  Scenario: Cleanup on early error
    Given I run "stardoc ./test-docs"
    And a temp directory has been created
    When an error occurs before the server starts
    Then the temp directory should be cleaned up
    And the CLI should exit with a non-zero code
    And the error message should be displayed before cleanup

  Scenario: Cleanup child processes
    Given I run "stardoc ./test-docs"
    And the Astro dev server is running with PID 12345
    When I exit stardoc
    Then the process with PID 12345 should be terminated
    And no child processes should remain running
    And no zombie processes should be created

  Scenario: Cleanup with multiple running processes
    Given I run "stardoc ./test-docs"
    And the dev server spawns multiple worker processes
    When I exit stardoc
    Then all worker processes should be terminated
    And the entire process tree should be cleaned up

  Scenario: Cleanup on panic
    Given I run "stardoc ./test-docs"
    And a temp directory has been created
    When a panic occurs in the Go code
    Then the defer cleanup handlers should execute
    And the temp directory should be removed
    And child processes should be terminated

  Scenario: Partial cleanup on force exit
    Given I run "stardoc ./test-docs"
    And I press Ctrl+C
    When I press Ctrl+C again to force exit
    Then the CLI should attempt best-effort cleanup
    And the temp directory removal should be attempted
    And child process termination should be attempted
    But the CLI should not wait more than 2 seconds before force exiting

  Scenario: Cleanup verification logging
    Given I run "stardoc ./test-docs"
    When I exit stardoc
    Then the CLI should log "Stopping dev server..."
    And the CLI should log "Cleaning up workspace..."
    And the CLI should log "Cleanup complete" on success

  Scenario: Handle cleanup errors gracefully
    Given I run "stardoc ./test-docs"
    And the temp directory has become read-only
    When I exit stardoc
    Then the CLI should log a warning about cleanup failure
    And the warning should include the temp directory path
    And the CLI should still attempt to stop child processes
    And the CLI should exit with code 0 (user's exit was successful)

  Scenario: No orphaned temp directories from previous runs
    Given stardoc was forcefully killed in a previous run
    And orphaned temp directory "/tmp/stardoc-old123" exists
    When I run "stardoc ./test-docs"
    Then the CLI should not clean up "/tmp/stardoc-old123"
    And the CLI should create its own new temp directory
    And the CLI should log its own temp directory path

  Scenario: Cleanup node_modules if present
    Given I run "stardoc ./test-docs"
    And npm has installed dependencies in the temp directory
    And "node_modules" folder exists in temp workspace
    When I exit stardoc
    Then the entire temp directory including "node_modules" should be removed
    And the removal should not hang or take excessive time
