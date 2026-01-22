Feature: Signal Handling
  As a stardoc user
  I want the tool to handle interrupt signals gracefully
  So that temporary files are cleaned up and child processes are terminated properly

  Background:
    Given the stardoc CLI is available
    And a source directory "./test-docs" exists with markdown files
    And I have started "stardoc ./test-docs"
    And the dev server is running

  Scenario: Handle SIGINT (Ctrl+C)
    When I press Ctrl+C
    Then the CLI should catch the SIGINT signal
    And the CLI should display "Shutting down gracefully..."
    And the dev server process should be terminated
    And the temp directory should be cleaned up
    And the CLI should exit with code 0

  Scenario: Handle SIGTERM
    When I send a SIGTERM signal to the stardoc process
    Then the CLI should catch the SIGTERM signal
    And the CLI should display "Shutting down gracefully..."
    And the dev server process should be terminated
    And the temp directory should be cleaned up
    And the CLI should exit with code 0

  Scenario: Multiple interrupt signals
    When I press Ctrl+C
    And I press Ctrl+C again within 1 second
    Then the CLI should force exit immediately
    And the CLI should display "Force stopping..."
    And the CLI should exit with code 1
    And best-effort cleanup should be attempted

  Scenario: Interrupt during dependency installation
    Given stardoc is installing npm dependencies
    When I press Ctrl+C
    Then the npm install process should be terminated
    And the temp directory should be cleaned up
    And the CLI should exit with code 0

  Scenario: Interrupt during server startup
    Given stardoc is starting the Astro dev server
    And the server has not finished starting yet
    When I press Ctrl+C
    Then the Astro process should be terminated
    And the temp directory should be cleaned up
    And the CLI should exit with code 0

  Scenario: Interrupt with child processes running
    Given the Astro dev server is running with PID 12345
    When I press Ctrl+C
    Then the CLI should send SIGTERM to PID 12345
    And the CLI should wait up to 5 seconds for graceful shutdown
    And if the process doesn't stop, send SIGKILL
    And the CLI should not exit until all child processes are stopped

  Scenario: Cleanup timeout handling
    Given cleanup operations are taking longer than expected
    When I press Ctrl+C
    Then the CLI should start cleanup
    And the CLI should wait up to 10 seconds for cleanup completion
    And if cleanup doesn't complete, force exit with a warning
    And the warning should mention "cleanup may be incomplete"

  Scenario: Signal handling does not interfere with normal operation
    Given the dev server is running normally
    And no interrupt signal has been sent
    When the user browses the documentation site
    Then the server should respond normally
    And the signal handler should not consume excessive resources
