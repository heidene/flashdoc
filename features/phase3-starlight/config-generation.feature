Feature: Config Generation
  As a stardoc user
  I want the Astro config to be dynamically generated with my custom title
  So that the documentation site reflects my project name

  Background:
    Given the stardoc CLI is available
    And a temp workspace exists at "/tmp/stardoc-abc123"
    And the Starlight template has been extracted

  Scenario: Generate config with default title
    Given I run "stardoc ./docs"
    When the astro.config.mjs is generated
    Then the config should contain:
      """
      title: 'Docs'
      """
    And the title should be derived from the directory name

  Scenario: Generate config with custom title flag
    Given I run "stardoc ./docs --title 'API Documentation'"
    When the astro.config.mjs is generated
    Then the config should contain:
      """
      title: 'API Documentation'
      """

  Scenario: Replace template placeholder
    Given the template astro.config.mjs contains:
      """
      export default defineConfig({
        integrations: [
          starlight({
            title: '{{SITE_TITLE}}',
          }),
        ],
      });
      """
    When the config is generated with title "My Project"
    Then the {{SITE_TITLE}} placeholder should be replaced with "My Project"
    And the final config should contain:
      """
      title: 'My Project',
      """

  Scenario: Generate title from directory name with hyphens
    Given I run "stardoc ./my-awesome-project"
    When the astro.config.mjs is generated
    Then the default title should be "My Awesome Project"

  Scenario: Generate title from directory name with underscores
    Given I run "stardoc ./my_api_docs"
    When the astro.config.mjs is generated
    Then the default title should be "My Api Docs"

  Scenario: Handle special characters in auto-generated title
    Given I run "stardoc ./docs & guides"
    When the astro.config.mjs is generated
    Then the title should properly escape special characters
    And the title should be "Docs & Guides"

  Scenario: Generate config with sidebar autogeneration
    When the astro.config.mjs is generated
    Then it should not include explicit sidebar configuration
    And the sidebar should use Starlight's default autogeneration
    And the sidebar should automatically reflect the file structure

  Scenario Outline: Generate config without social links
    When the astro.config.mjs is generated
    Then <social_handling>

    Examples:
      | social_handling                         |
      | the config should not include a "social" section |
      | the "social" section should be empty    |

  Scenario Outline: Generate config with default locale handling
    When the astro.config.mjs is generated
    Then <locale_handling>

    Examples:
      | locale_handling                                    |
      | the config should include defaultLocale 'en'       |
      | the locale should be omitted to use Starlight defaults |

  Scenario: Config is valid JavaScript
    When the astro.config.mjs is generated
    Then the file should be valid JavaScript/ES6 syntax
    And it should be parseable by Node.js
    And running "node -c astro.config.mjs" should succeed

  Scenario Outline: Handle title with quotes
    Given I run "stardoc ./docs --title 'John's Documentation'"
    When the astro.config.mjs is generated
    Then the title should be properly escaped
    And the config should contain <quote_format>

    Examples:
      | quote_format                            |
      | single quotes with escaping like 'John\'s Documentation' |
      | double quotes like "John's Documentation" |

  Scenario: Handle very long titles
    Given I run "stardoc ./docs --title 'This is a very long title that exceeds normal length expectations for documentation sites'"
    When the astro.config.mjs is generated
    Then the full title should be preserved in the config
    And no truncation should occur

  Scenario: Generate config preserves template structure
    Given the template astro.config.mjs has custom integrations
    When the config is generated
    Then the custom integrations should be preserved
    And only the title placeholder should be replaced
    And other configuration options should remain intact

  Scenario: Config generation error handling
    Given the template astro.config.mjs is malformed
    When config generation is attempted
    Then an error should be logged
    And the error should indicate "failed to generate config"
    And the CLI should exit with code 1

  Scenario: Config includes documentation URL
    When the astro.config.mjs is generated
    Then the config may include a comment with Starlight docs URL
    And the comment should be "// Learn more: https://starlight.astro.build/"

  Scenario: Config supports future customization
    When the astro.config.mjs is generated
    Then the file should include helpful comments
    And comments should guide users on common customizations
    But the file should remain minimal by default
