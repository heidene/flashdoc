Feature: Package Manager Detection
  As a stardoc user
  I want the tool to automatically detect my preferred package manager
  So that dependencies are installed using my existing tools

  Background:
    Given the stardoc CLI is available

  Scenario: Detect pnpm when available
    Given "pnpm" is installed and available in PATH
    And "bun" is not available
    And "npm" is available
    When stardoc detects the package manager
    Then it should select "pnpm"
    And it should log "Using package manager: pnpm"

  Scenario: Detect bun when pnpm is not available
    Given "pnpm" is not available
    And "bun" is installed and available in PATH
    And "npm" is available
    When stardoc detects the package manager
    Then it should select "bun"
    And it should log "Using package manager: bun"

  Scenario: Fall back to npm
    Given "pnpm" is not available
    And "bun" is not available
    And "npm" is available
    When stardoc detects the package manager
    Then it should select "npm"
    And it should log "Using package manager: npm"

  Scenario: Detection priority order
    Given "pnpm", "bun", and "npm" are all available
    When stardoc detects the package manager
    Then it should select "pnpm"
    And the priority order should be: pnpm > bun > npm

  Scenario: No package manager available
    Given "pnpm" is not available
    And "bun" is not available
    And "npm" is not available
    When stardoc detects the package manager
    Then an error should be returned
    And the error should state "no package manager found (tried: pnpm, bun, npm)"
    And the CLI should exit with code 1

  Scenario: Package manager detection via version check
    Given "pnpm" is available
    When stardoc detects the package manager
    Then it should execute "pnpm --version" to verify
    And if the command succeeds, pnpm should be selected
    And if the command fails, it should try the next package manager

  Scenario: Respect package manager override flag
    Given I run "stardoc ./docs --package-manager npm"
    And "pnpm" is available
    When stardoc detects the package manager
    Then it should use "npm" instead of "pnpm"
    And it should skip the detection process

  Scenario: Invalid package manager override
    Given I run "stardoc ./docs --package-manager yarn"
    When stardoc validates the package manager
    Then an error should be returned
    And the error should state "unsupported package manager: yarn"
    And the error should list supported managers: "pnpm, bun, npm"
    And the CLI should exit with code 1

  Scenario: Package manager detection caching
    Given stardoc has detected "pnpm" for the current run
    When dependency installation occurs multiple times
    Then the detection should only happen once
    And subsequent operations should reuse the detected manager

  Scenario: Package manager in PATH but not executable
    Given "pnpm" exists in PATH but has no execute permissions
    When stardoc detects the package manager
    Then it should treat "pnpm" as unavailable
    And it should continue to the next option (bun or npm)

  Scenario: Detect package manager version compatibility
    Given "pnpm" is available with version "9.0.0"
    When stardoc detects the package manager
    Then it should check if the version is compatible
    And versions >= 8.0.0 should be supported for pnpm
    And versions >= 1.0.0 should be supported for bun
    And versions >= 8.0.0 should be supported for npm

  Scenario: Warn about old package manager versions
    Given "npm" is available with version "6.14.0"
    When stardoc detects the package manager
    Then it should log a warning "npm version 6.14.0 is old, consider upgrading to >= 8.0.0"
    But it should still proceed with npm

  Scenario: Package manager detection on different platforms
    Given the operating system is "darwin", "linux", or "windows"
    When stardoc detects the package manager
    Then the detection should work consistently across platforms
    And the command execution should use platform-appropriate paths

  Scenario: Log detection process
    Given I run "stardoc ./docs"
    When stardoc detects the package manager
    Then it should log "Detecting package manager..."
    And it should log "Checking for pnpm..."
    And if not found, "Checking for bun..."
    And if not found, "Checking for npm..."
    And finally, "Using package manager: <selected>"

  Scenario Outline: Detection with user agent environment variable
    Given the environment variable "npm_config_user_agent" contains "pnpm"
    When stardoc detects the package manager
    Then <detection_behavior>

    Examples:
      | detection_behavior                              |
      | it should prefer "pnpm" based on the user agent |
      | it should use the standard detection process    |
