Feature: CLI Argument Parsing
  As a user
  I want to provide a source directory path to stardoc
  So that I can view its markdown files as a Starlight documentation site

  Background:
    Given the stardoc CLI is available

  Scenario: Valid directory path provided
    When I run "stardoc ./my-docs"
    Then the CLI should parse the path "./my-docs"
    And the CLI should validate that the path exists
    And the CLI should proceed with site generation

  Scenario: Absolute path provided
    When I run "stardoc /Users/me/documents/guides"
    Then the CLI should parse the path "/Users/me/documents/guides"
    And the CLI should validate that the path exists
    And the CLI should proceed with site generation

  Scenario: No arguments provided
    When I run "stardoc" without arguments
    Then the CLI should display usage information
    And the CLI should exit with code 1
    And the error message should contain "Usage: stardoc <directory>"

  Scenario: Too many arguments provided
    When I run "stardoc ./docs ./more-docs"
    Then the CLI should display an error "too many arguments"
    And the CLI should display usage information
    And the CLI should exit with code 1

  Scenario: Directory does not exist
    Given the directory "./nonexistent" does not exist
    When I run "stardoc ./nonexistent"
    Then the CLI should display an error "directory not found: ./nonexistent"
    And the CLI should exit with code 1

  Scenario: Path is a file not a directory
    Given the file "./README.md" exists
    When I run "stardoc ./README.md"
    Then the CLI should display an error "path is not a directory: ./README.md"
    And the CLI should exit with code 1

  Scenario: Help flag
    When I run "stardoc --help"
    Then the CLI should display usage information
    And the CLI should display available flags
    And the CLI should exit with code 0

  Scenario: Version flag
    When I run "stardoc --version"
    Then the CLI should display "stardoc version X.Y.Z"
    And the CLI should exit with code 0

  Scenario: Custom title flag
    When I run "stardoc ./docs --title 'My Custom Docs'"
    Then the CLI should parse the title as "My Custom Docs"
    And the CLI should use this title in the generated site

  Scenario: Port flag
    When I run "stardoc ./docs --port 4321"
    Then the CLI should use port 4321 for the dev server
    And the CLI should validate that the port is between 1024 and 65535

  Scenario: Invalid port flag
    When I run "stardoc ./docs --port 99999"
    Then the CLI should display an error "invalid port: must be between 1024 and 65535"
    And the CLI should exit with code 1

  Scenario: No-open flag
    When I run "stardoc ./docs --no-open"
    Then the CLI should not attempt to open a browser
    And the dev server should still start normally
