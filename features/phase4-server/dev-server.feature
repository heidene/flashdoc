Feature: Dev Server
  As a stardoc user
  I want the Astro dev server to start automatically
  So that I can immediately view my documentation in a browser

  Background:
    Given the stardoc CLI is available
    And a temp workspace exists with all dependencies installed
    And markdown files have been processed and copied

  Scenario: Start dev server with default port
    When the dev server is started
    Then stardoc should execute "npm run dev" (or equivalent) in the temp workspace
    And the server should listen on port 4321
    And the CLI should log "Dev server starting..."

  Scenario: Dev server starts successfully
    Given the dev server is starting
    When the server is ready
    Then the CLI should detect "astro dev" ready message in output
    And the CLI should log "Server ready at http://localhost:4321"
    And the server process should be running

  Scenario: Stream server output
    Given the dev server is running
    When Astro logs messages
    Then the output should be streamed to stardoc's terminal
    And build messages should be visible
    And request logs should be visible
    And error messages should be visible

  Scenario: Server startup with custom port
    Given I run "stardoc ./docs --port 8080"
    When the dev server is started
    Then the server should listen on port 8080
    And the CLI should log "Server ready at http://localhost:8080"

  Scenario: Port already in use
    Given port 4321 is already occupied by another process
    When the dev server is started
    Then Astro should automatically try the next available port (e.g., 4322)
    And the CLI should log "Port 4321 in use, trying 4322..."
    And the actual port should be extracted from Astro output
    And the browser should open with the correct port

  Scenario: Server startup failure
    Given the astro.config.mjs has a syntax error
    When the dev server is started
    Then the startup should fail
    And the error message should be displayed
    And the CLI should log "Failed to start dev server"
    And the CLI should exit with code 1

  Scenario: Monitor server process
    Given the dev server is running with PID 12345
    When stardoc checks the process status
    Then the process should be alive
    And stardoc should maintain a reference to the process

  Scenario: Server crashes during runtime
    Given the dev server is running
    When the Astro process crashes unexpectedly
    Then stardoc should detect the crash
    And the CLI should log "Dev server crashed unexpectedly"
    And stardoc should display the crash reason from stderr
    And the CLI should exit with code 1

  Scenario: Server responds to requests
    Given the dev server is running on port 4321
    When I access "http://localhost:4321" in a browser
    Then the Starlight documentation site should load
    And the site should display the processed markdown content
    And navigation should reflect the file structure

  Scenario: Hot reload works
    Given the dev server is running
    When a markdown file is modified in the temp workspace
    Then Astro should detect the change
    And the browser should hot-reload the page
    And the updated content should be visible

  Scenario: Server shutdown on exit
    Given the dev server is running with PID 12345
    When I press Ctrl+C
    Then stardoc should send SIGTERM to PID 12345
    And the server should shut down gracefully
    And the CLI should log "Stopping dev server..."

  Scenario: Force kill server if graceful shutdown fails
    Given the dev server is running with PID 12345
    When stardoc sends SIGTERM
    And the process doesn't stop within 5 seconds
    Then stardoc should send SIGKILL
    And the process should be terminated forcefully

  Scenario: Server uses correct package manager
    Given "pnpm" was detected as the package manager
    When the dev server is started
    Then stardoc should execute "pnpm run dev"
    And not "npm run dev"

  Scenario: Server starts with environment variables
    When the dev server is started
    Then environment variables should be set:
      | Variable | Value |
      | NODE_ENV | development |
    And these variables should be passed to the Astro process

  Scenario: Parse server URL from output
    Given the dev server starts and logs:
      """
      🚀 astro  v4.0.0 started in 345ms

      ┃ Local    http://localhost:4321/
      ┃ Network  use --host to expose
      """
    When stardoc parses the output
    Then it should extract "http://localhost:4321/" as the server URL
    And this URL should be used for opening the browser

  Scenario: Handle non-standard Astro output
    Given the dev server logs output in an unexpected format
    When stardoc tries to detect server readiness
    Then it should fall back to a timeout-based approach
    And it should wait 10 seconds and assume the server is ready
    And it should log "Server detection timed out, assuming ready"

  Scenario: Server keeps running until interrupted
    Given the dev server is running
    When no interrupt signal is sent
    Then the server should continue running indefinitely
    And stardoc should wait for user input or signals
    And the CLI should not exit prematurely

  Scenario: Multiple dev servers in parallel
    Given I run "stardoc ./docs-a" in terminal 1
    And I run "stardoc ./docs-b" in terminal 2
    When both dev servers start
    Then each should use a different port
    And both should run independently
    And there should be no port conflicts

  Scenario: Dev server uses production-like settings
    When the dev server is started
    Then Astro should run in development mode
    And hot reload should be enabled
    And source maps should be enabled
    But the site should closely resemble the production build
