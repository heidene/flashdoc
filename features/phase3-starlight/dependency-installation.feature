Feature: Dependency Installation
  As a stardoc user
  I want dependencies to be installed automatically
  So that I can view my documentation without manual setup

  Background:
    Given the stardoc CLI is available
    And a temp workspace exists at "/tmp/stardoc-abc123"
    And the Starlight template has been extracted
    And a package manager has been detected

  Scenario: Install dependencies with pnpm
    Given "pnpm" is the selected package manager
    When dependency installation is triggered
    Then stardoc should execute "pnpm install" in the temp workspace
    And the command should run with output streaming
    And node_modules should be created

  Scenario: Install dependencies with bun
    Given "bun" is the selected package manager
    When dependency installation is triggered
    Then stardoc should execute "bun install" in the temp workspace
    And the command should run with output streaming
    And node_modules should be created

  Scenario: Install dependencies with npm
    Given "npm" is the selected package manager
    When dependency installation is triggered
    Then stardoc should execute "npm install" in the temp workspace
    And the command should run with output streaming
    And node_modules should be created

  Scenario: Display installation progress
    Given dependency installation is in progress
    When npm is installing packages
    Then the CLI should display "Installing dependencies..."
    And package manager output should be visible to the user
    And the output should update in real-time

  Scenario: Successful installation
    Given dependency installation completes successfully
    When the install command exits with code 0
    Then the CLI should log "Dependencies installed successfully"
    And the CLI should proceed to the next step (server startup)

  Scenario: Installation failure
    Given dependency installation is triggered
    When the install command exits with code 1
    Then the CLI should log an error "dependency installation failed"
    And the package manager's error output should be displayed
    And the CLI should exit with code 1
    And the temp workspace should be cleaned up

  Scenario: Network error during installation
    Given the network is unavailable
    When dependency installation is triggered
    Then the install command should fail
    And the error should indicate a network issue
    And the CLI should suggest checking internet connectivity
    And the CLI should exit with code 1

  Scenario: Corrupted package.json
    Given the package.json in the temp workspace is malformed
    When dependency installation is triggered
    Then the install command should fail
    And the error should mention "invalid package.json"
    And the CLI should exit with code 1

  Scenario: Install with frozen lockfile (pnpm)
    Given "pnpm" is the selected package manager
    And the template includes a "pnpm-lock.yaml" file
    When dependency installation is triggered
    Then stardoc should execute "pnpm install --frozen-lockfile"
    And the lockfile should not be modified

  Scenario: Install with frozen lockfile (npm)
    Given "npm" is the selected package manager
    And the template includes a "package-lock.json" file
    When dependency installation is triggered
    Then stardoc should execute "npm ci"
    And the lockfile should not be modified

  Scenario: Install with frozen lockfile (bun)
    Given "bun" is the selected package manager
    And the template includes a "bun.lockb" file
    When dependency installation is triggered
    Then stardoc should execute "bun install --frozen-lockfile"
    And the lockfile should not be modified

  Scenario: No lockfile present
    Given the template does not include a lockfile
    When dependency installation is triggered
    Then stardoc should use standard "install" command
    And a new lockfile should be generated
    And this is acceptable for a temporary workspace

  Scenario: Installation timeout
    Given dependency installation is taking longer than 5 minutes
    When the timeout is reached
    Then the CLI should log a warning "installation is taking longer than expected"
    But the installation should continue
    And the CLI should not force-kill the process

  Scenario: Handle SIGINT during installation
    Given dependency installation is in progress
    When the user presses Ctrl+C
    Then the install process should be terminated
    And the CLI should log "Installation cancelled"
    And the temp workspace should be cleaned up
    And the CLI should exit with code 0 (user-initiated)

  Scenario: Retry on transient failure
    Given dependency installation fails with a transient error
    When the error is detected as potentially transient
    Then the CLI may retry the installation once
    And it should log "Retrying installation..."
    But it should not retry on permanent errors (e.g., malformed JSON)

  Scenario: Silent install option
    Given I run "stardoc ./docs --silent"
    When dependency installation is triggered
    Then package manager output should be suppressed
    And only stardoc's summary messages should be shown
    And the install command should run with silent/quiet flags

  Scenario: Install creates node_modules
    Given dependency installation completes successfully
    When I check the temp workspace
    Then a "node_modules" directory should exist
    And it should contain "@astrojs/starlight"
    And it should contain "astro"

  Scenario: Install duration logging
    Given dependency installation takes 23 seconds
    When the installation completes
    Then the CLI should log "Dependencies installed in 23s"

  Scenario: Parallel installations are isolated
    Given I run "stardoc ./docs-a" in terminal 1
    And I run "stardoc ./docs-b" in terminal 2
    When both instances install dependencies simultaneously
    Then each should use its own temp workspace
    And the installations should not conflict
    And both should succeed independently

  Scenario: Installation with npm audit warnings
    Given "npm" is the selected package manager
    When "npm install" produces audit warnings
    Then the warnings should be visible to the user
    But stardoc should still proceed (warnings are not errors)
    And the dev server should start normally
