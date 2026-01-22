Feature: Browser Open
  As a stardoc user
  I want my default browser to open automatically when the server is ready
  So that I can immediately view my documentation without manual steps

  Background:
    Given the stardoc CLI is available
    And the dev server has started successfully
    And the server is listening on "http://localhost:4321"

  Scenario: Open browser on macOS
    Given the operating system is "darwin"
    When the server is ready
    Then stardoc should execute "open http://localhost:4321"
    And the default browser should launch
    And the documentation site should load

  Scenario: Open browser on Linux
    Given the operating system is "linux"
    When the server is ready
    Then stardoc should execute "xdg-open http://localhost:4321"
    And the default browser should launch
    And the documentation site should load

  Scenario: Open browser on Windows
    Given the operating system is "windows"
    When the server is ready
    Then stardoc should execute "cmd /c start http://localhost:4321"
    And the default browser should launch
    And the documentation site should load

  Scenario: Wait for server before opening browser
    Given the dev server is starting but not ready yet
    When stardoc attempts to open the browser
    Then it should wait until the server is ready
    And it should not open the browser prematurely
    And the URL should be accessible when the browser opens

  Scenario: Browser opens after server ready log
    Given the dev server logs "Local    http://localhost:4321/"
    When stardoc detects this message
    Then it should wait 1 additional second
    And then open the browser
    And this ensures the server is fully ready

  Scenario: Do not open browser with --no-open flag
    Given I run "stardoc ./docs --no-open"
    When the server is ready
    Then the browser should not be opened
    And the CLI should log "Server ready at http://localhost:4321"
    And the CLI should log "(browser not opened due to --no-open flag)"

  Scenario: Handle browser open failure
    Given the "open" command is not available
    When stardoc attempts to open the browser
    Then the command should fail
    And stardoc should log a warning "failed to open browser"
    And the CLI should still continue running
    And the server should remain running
    And the URL should be displayed for manual opening

  Scenario: Browser is already running
    Given the default browser is already open with other tabs
    When stardoc opens the URL
    Then a new tab should open in the existing browser window
    And the existing tabs should not be affected

  Scenario: Open with custom port
    Given the dev server is running on port 8080
    When the browser is opened
    Then the URL should be "http://localhost:8080"
    And not the default port 4321

  Scenario: Open URL with base path
    Given the Astro config has a base path "/docs"
    When the browser is opened
    Then the URL should be "http://localhost:4321/docs"
    And the base path should be included

  Scenario: Log browser open action
    When the browser is opened
    Then the CLI should log "Opening browser..."
    And the CLI should log the exact URL being opened

  Scenario: Network interface binding
    Given the dev server binds to "0.0.0.0"
    And the server logs show "Network  http://192.168.1.100:4321/"
    When the browser is opened
    Then stardoc should use "http://localhost:4321"
    And not the network IP address

  Scenario: Browser open timeout
    Given the "open" command hangs
    When stardoc attempts to open the browser
    Then it should timeout after 5 seconds
    And it should log a warning about the timeout
    But the server should continue running normally

  Scenario: Concurrent browser opens
    Given I run "stardoc ./docs-a" in terminal 1
    And I run "stardoc ./docs-b" in terminal 2
    When both servers are ready
    Then two browser tabs should open
    And each tab should point to the correct port
    And there should be no confusion between instances

  Scenario Outline: Browser preference detection
    Given the environment variable "BROWSER" is set to "firefox"
    When stardoc opens the browser
    Then <browser_selection>

    Examples:
      | browser_selection                              |
      | it should respect the BROWSER environment variable |
      | it should use the system default browser       |

  Scenario: Headless environment detection
    Given the environment has no display (CI/headless server)
    And the DISPLAY environment variable is not set (Linux)
    When stardoc attempts to open the browser
    Then it should detect the headless environment
    And it should skip opening the browser
    And it should log "Headless environment detected, skipping browser open"
    And the server should continue running

  Scenario: Open browser exactly once
    Given the dev server is ready
    When stardoc opens the browser
    Then the browser should open exactly once
    And subsequent hot reloads should not open additional tabs
