Feature: Terminal Output
  As a stardoc user
  I want clear and informative terminal output
  So that I understand what the tool is doing and can diagnose issues

  Background:
    Given the stardoc CLI is available

  Scenario: Clean startup output
    Given I run "stardoc ./docs"
    Then the output should follow this sequence:
      """
      🚀 Stardoc - Ephemeral Documentation Viewer

      📁 Source: ./docs
      📦 Workspace: /tmp/stardoc-abc123
      🔍 Found 12 markdown files

      📥 Installing dependencies...
      ✅ Dependencies installed in 8s

      🚀 Starting dev server...
      ✅ Server ready at http://localhost:4321

      🌐 Opening browser...

      Press Ctrl+C to stop
      """

  Scenario: Use emoji consistently
    When stardoc produces output
    Then emojis should be used for visual clarity:
      | Symbol | Meaning |
      | 🚀     | Starting/launching |
      | 📁     | Directory/folder |
      | 📦     | Workspace/package |
      | 🔍     | Scanning/searching |
      | 📥     | Installing/downloading |
      | ✅     | Success |
      | ❌     | Error |
      | ⚠️     | Warning |
      | 🌐     | Browser/web |
      | 🧹     | Cleanup |

  Scenario: Progress indicators for long operations
    Given dependency installation is in progress
    When the operation takes more than 2 seconds
    Then a spinner or progress indicator should be shown
    And the indicator should update to show activity

  Scenario: Streaming output from subprocesses
    Given the dev server is running
    When Astro logs output
    Then the output should be streamed in real-time
    And it should be prefixed with appropriate labels
    And Astro's colored output should be preserved

  Scenario: Error output is distinct
    Given an error occurs during setup
    When stardoc displays the error
    Then the error should be prefixed with "❌ Error:"
    And the error message should be in red (if colors are supported)
    And the error should be clearly visible

  Scenario: Warning output is distinct
    Given a non-fatal issue occurs
    When stardoc displays a warning
    Then the warning should be prefixed with "⚠️  Warning:"
    And the warning message should be in yellow (if colors are supported)

  Scenario: Verbose mode
    Given I run "stardoc ./docs --verbose"
    When stardoc operates
    Then additional debug information should be logged
    And internal operations should be visible
    And timing information should be included

  Scenario: Quiet mode
    Given I run "stardoc ./docs --quiet"
    When stardoc operates
    Then only essential messages should be displayed
    And progress indicators should be suppressed
    And subprocess output should be suppressed
    And only the final URL and errors should be shown

  Scenario: Color support detection
    Given the terminal supports colors
    When stardoc produces output
    Then colors should be used for emphasis
    And ANSI color codes should be applied

  Scenario: No color support
    Given the terminal does not support colors (dumb terminal)
    When stardoc produces output
    Then plain text should be used
    And ANSI color codes should not be applied

  Scenario: NO_COLOR environment variable
    Given the environment variable "NO_COLOR" is set
    When stardoc produces output
    Then colors should be disabled
    And all output should be plain text

  Scenario: Interactive vs non-interactive detection
    Given the output is being piped or redirected
    When stardoc operates
    Then it should detect non-interactive mode
    And it should disable progress spinners
    And it should use simpler output format

  Scenario Outline: Logging format options
    When stardoc logs messages
    Then messages should follow <format_type>

    Examples:
      | format_type                                      |
      | a structured format "[timestamp] [level] message" |
      | a simple user-facing format "emoji message"      |

  Scenario: Timestamp option
    Given I run "stardoc ./docs --timestamps"
    When stardoc logs messages
    Then each message should include a timestamp
    And the format should be "HH:MM:SS"

  Scenario: File processing summary
    Given 50 markdown files are processed
    When processing completes
    Then the output should include a summary:
      """
      📊 Summary:
         • 50 files processed
         • 3 without frontmatter (fixed)
         • 12 directories created
      """

  Scenario: Cleanup messages
    Given I press Ctrl+C to exit
    When cleanup occurs
    Then the output should show:
      """
      🧹 Cleaning up...
      ✅ Workspace removed
      ✅ Server stopped
      👋 Goodbye!
      """

  Scenario: Error context and suggestions
    Given an error occurs because the directory doesn't exist
    When the error is displayed
    Then it should include context and a suggestion:
      """
      ❌ Error: directory not found: ./docs

      Make sure the path exists and try again.
      """

  Scenario: Help text formatting
    Given I run "stardoc --help"
    When the help text is displayed
    Then it should include command usage
    And it should include a description
    And it should include available flags with descriptions
    And it should include examples
    And the formatting should be clear and readable

  Scenario: Version output
    Given I run "stardoc --version"
    Then the output should display:
      """
      stardoc version 0.1.0
      """

  Scenario: Output buffering
    When stardoc produces output
    Then output should be flushed immediately
    And there should be no buffering delays
    And real-time feedback should be visible

  Scenario: Multiline error messages
    Given a complex error occurs with a stack trace
    When stardoc displays the error
    Then the error should be formatted with proper indentation
    And the stack trace should be readable
    And the most important information should be at the top

  Scenario: Server output passthrough
    Given the dev server is running
    When I access the site in a browser
    Then HTTP request logs should be visible in the terminal
    And Astro's build output should be visible
    And hot reload messages should be visible

  Scenario: Handle terminal width
    Given the terminal has a narrow width (e.g., 80 columns)
    When stardoc produces output
    Then long lines should be wrapped appropriately
    And the output should remain readable
    And ASCII art or tables should adapt to the width

  Scenario: Exit message clarity
    Given the dev server is running normally
    When I press Ctrl+C
    Then the output should clearly indicate shutdown:
      """

      ^C received
      🛑 Shutting down gracefully...
      """
