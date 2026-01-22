Feature: Frontmatter Injection
  As a stardoc user
  I want the tool to automatically add or fix frontmatter in my markdown files
  So that Starlight can properly display titles and metadata

  Background:
    Given the stardoc CLI is available
    And a temp workspace has been created

  Scenario: Add frontmatter to file without any
    Given a markdown file "guide.md" with content:
      """
      # My Guide

      This is the content.
      """
    When the file is processed
    Then the output should have frontmatter:
      """
      ---
      title: Guide
      ---
      # My Guide

      This is the content.
      """

  Scenario: Preserve existing frontmatter
    Given a markdown file "api.md" with content:
      """
      ---
      title: Custom API Title
      description: API documentation
      ---
      # API Reference

      Content here.
      """
    When the file is processed
    Then the frontmatter should remain unchanged
    And the title should still be "Custom API Title"
    And the description should still be "API documentation"

  Scenario: Add missing title to existing frontmatter
    Given a markdown file "guide.md" with content:
      """
      ---
      description: A helpful guide
      ---
      # Content

      More content.
      """
    When the file is processed
    Then the frontmatter should include:
      """
      ---
      title: Guide
      description: A helpful guide
      ---
      """

  Scenario: Generate title from filename
    Given a markdown file "getting-started.md"
    When the file is processed
    Then the generated title should be "Getting Started"

  Scenario: Generate title from kebab-case filename
    Given a markdown file "api-reference.md"
    When the file is processed
    Then the generated title should be "Api Reference"

  Scenario: Generate title from snake_case filename
    Given a markdown file "user_authentication.md"
    When the file is processed
    Then the generated title should be "User Authentication"

  Scenario: Handle numbered prefixes in filenames
    Given a markdown file "01-introduction.md"
    When the file is processed
    Then the generated title should be "Introduction"
    And the number prefix "01-" should be stripped

  Scenario: Handle special characters in filenames
    Given a markdown file "FAQ's & Tips.md"
    When the file is processed
    Then the generated title should be "FAQ's & Tips"

  Scenario: Handle index files
    Given a markdown file "index.md" in directory "guides"
    When the file is processed
    Then the generated title should be "Guides"
    And the title should be derived from the parent directory name

  Scenario: Handle README files
    Given a markdown file "README.md" in directory "api"
    When the file is processed
    Then the generated title should be "Api"
    And the title should be derived from the parent directory name

  Scenario: Handle README in root
    Given a markdown file "README.md" in the root source directory
    When the file is processed
    Then the generated title should be "Home"

  Scenario: Preserve existing title even with misleading filename
    Given a markdown file "temp-file-123.md" with frontmatter:
      """
      ---
      title: Important Document
      ---
      """
    When the file is processed
    Then the title should remain "Important Document"
    And the filename should not influence the title

  Scenario: Handle malformed frontmatter
    Given a markdown file "broken.md" with content:
      """
      ---
      title: Unclosed frontmatter
      # Heading

      Content
      """
    When the file is processed
    Then the malformed frontmatter should be detected
    And new valid frontmatter should be added at the top:
      """
      ---
      title: Broken
      ---
      """
    And the original malformed content should be preserved as body content

  Scenario: Handle empty frontmatter
    Given a markdown file "empty.md" with content:
      """
      ---
      ---
      # Content
      """
    When the file is processed
    Then the frontmatter should be populated with:
      """
      ---
      title: Empty
      ---
      """

  Scenario: Preserve other frontmatter fields
    Given a markdown file "advanced.md" with frontmatter:
      """
      ---
      title: Advanced Topics
      sidebar:
        order: 2
        badge: New
      tags:
        - advanced
        - tutorial
      ---
      """
    When the file is processed
    Then all frontmatter fields should be preserved
    And the structure should remain valid YAML

  Scenario: Handle unicode in titles
    Given a markdown file "日本語.md"
    When the file is processed
    Then the generated title should be "日本語"
    And the frontmatter should be valid UTF-8

  Scenario: Strip file extension from generated titles
    Given a markdown file "document.md"
    When the file is processed
    Then the generated title should be "Document"
    And the ".md" extension should not appear in the title
