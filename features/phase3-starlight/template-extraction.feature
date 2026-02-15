Feature: Template Extraction
  As a stardoc developer
  I want to embed a Starlight template in the binary
  So that users don't need network access or separate template files

  Background:
    Given the stardoc binary contains an embedded Starlight template
    And a temp workspace exists at "/tmp/stardoc-abc123"

  Scenario: Extract template to temp workspace
    When the template is extracted
    Then the temp workspace should contain a "package.json" file
    And the temp workspace should contain an "astro.config.mjs" file
    And the temp workspace should contain a "tsconfig.json" file
    And the temp workspace should contain a "src" directory

  Scenario: Template contains Starlight dependencies
    When the template is extracted
    And I read the "package.json" file
    Then the dependencies should include "@astrojs/starlight"
    And the dependencies should include "astro"
    And the devDependencies should include necessary TypeScript types

  Scenario: Template package.json has correct structure
    When the template is extracted
    And I read the "package.json" file
    Then it should have a "name" field set to "flashdoc-site"
    And it should have a "type" field set to "module"
    And it should have a "scripts" section with "dev" and "build" scripts

  Scenario: Template astro.config.mjs is valid
    When the template is extracted
    And I read the "astro.config.mjs" file
    Then it should import "starlight" integration
    And it should have a placeholder for site title
    And the placeholder should be "{{SITE_TITLE}}"

  Scenario: Template directory structure matches Starlight conventions
    When the template is extracted
    Then the following structure should exist:
      """
      /tmp/stardoc-abc123/
      ├── package.json
      ├── astro.config.mjs
      ├── tsconfig.json
      ├── public/
      └── src/
          └── content/
              └── docs/
      """

  Scenario: Template uses latest stable Starlight version
    When the template is extracted
    And I read the "package.json" file
    Then the "@astrojs/starlight" version should be "^0.29.0" or newer
    And the version should not be a pre-release version

  Scenario: Template extraction is idempotent
    When the template is extracted
    And I extract the template again
    Then the files should be overwritten with the same content
    And no duplicate files should be created

  Scenario: Handle extraction errors
    Given the temp workspace is not writable
    When the template extraction is attempted
    Then an error should be returned
    And the error should indicate "failed to extract template"
    And the CLI should exit with code 1

  Scenario: Template contains minimal boilerplate
    When the template is extracted
    Then there should be no example content pages
    And there should be no placeholder documentation
    And the "src/content/docs" directory should be empty
    And users' markdown files will populate this directory

  Scenario: Template tsconfig.json extends Starlight defaults
    When the template is extracted
    And I read the "tsconfig.json" file
    Then it should extend "astro/tsconfigs/strict"
    And it should include ".astro/types.d.ts" in the includes
    And it should exclude "dist" directory

  Scenario: Template is embedded using go:embed
    Given the stardoc source code
    When I examine the template package
    Then it should use "//go:embed" directive
    And the embedded files should be from "templates/starlight/*"
    And the embed should include all necessary template files

  Scenario Outline: Template version is tracked
    When the template is extracted
    Then <version_tracking>

    Examples:
      | version_tracking                                           |
      | the package.json should include a comment or field indicating stardoc version |
      | the CLI should log "Using embedded Starlight template v0.29.0" |

  Scenario: Template supports custom Starlight config
    When the template is extracted
    And I read the "astro.config.mjs" file
    Then it should support sidebar configuration
    And it should support custom branding
    And it should have social links disabled by default
    But the structure should allow easy customization
