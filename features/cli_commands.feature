@cli
Feature: CLI Commands
  As a user
  I want to use command-line arguments
  So that I can quickly access specific functionality

  # ============================================================================
  # Version Information
  # ============================================================================

  @happy @smoke
  Scenario: Display version information
    Given the application is installed
    When I run "kariya --version"
    Then I should see version information
    And the exit code should be 0

  @happy
  Scenario: Display version with short flag
    Given the application is installed
    When I run "kariya -v"
    Then I should see version information
    And the exit code should be 0

  # ============================================================================
  # Help Information
  # ============================================================================

  @happy @smoke
  Scenario: Display help information
    Given the application is installed
    When I run "kariya --help"
    Then I should see help information
    And I should see available commands
    And the exit code should be 0

  @happy
  Scenario: Display help with short flag
    Given the application is installed
    When I run "kariya -h"
    Then I should see help information
    And the exit code should be 0

  # ============================================================================
  # Direct Intent Launch
  # ============================================================================

  @happy
  Scenario: Launch browse timeline directly
    Given the application is installed
    When I run "kariya browse"
    Then the application should start in browse timeline

  @happy
  Scenario: Launch capture event directly
    Given the application is installed
    When I run "kariya capture"
    Then the application should start in capture event

  @happy
  Scenario: Launch skills management directly
    Given the application is installed
    When I run "kariya skills"
    Then the application should start in skills management

  @happy
  Scenario: Launch configuration directly
    Given the application is installed
    When I run "kariya config"
    Then the application should start in configuration

  # ============================================================================
  # Quick Capture Commands
  # ============================================================================

  @happy
  Scenario: Quick capture event from CLI
    Given the application is installed
    When I run "kariya add 'Built REST API' --company 'Acme Corp'"
    Then the event should be saved
    And I should see confirmation message
    And the exit code should be 0

  @happy
  Scenario: Quick capture with multiple flags
    Given the application is installed
    When I run "kariya add 'Deployed Kubernetes' --company 'TechCo' --category 'technical' --project 'Infrastructure'"
    Then the event should be saved with all metadata
    And the exit code should be 0

  # ============================================================================
  # Export Commands
  # ============================================================================

  @happy
  Scenario: Export CV to JSON
    Given the application is installed
    And I have events in my timeline
    When I run "kariya export --format json"
    Then I should receive JSON output
    And the exit code should be 0

  @happy
  Scenario: Export CV to YAML
    Given the application is installed
    And I have events in my timeline
    When I run "kariya export --format yaml"
    Then I should receive YAML output
    And the exit code should be 0

  @happy
  Scenario: Export CV to file
    Given the application is installed
    And I have events in my timeline
    When I run "kariya export --format json --output cv.json"
    Then the file "cv.json" should exist
    And it should contain valid JSON
    And the exit code should be 0

  # ============================================================================
  # Import Commands
  # ============================================================================

  @happy
  Scenario: Import events from JSON
    Given the application is installed
    And I have a valid JSON file "events.json"
    When I run "kariya import events.json"
    Then the events should be imported
    And I should see import summary
    And the exit code should be 0

  @sad
  Scenario: Import from invalid file
    Given the application is installed
    And I have an invalid JSON file "bad.json"
    When I run "kariya import bad.json"
    Then I should see an error message
    And the exit code should be non-zero

  # ============================================================================
  # Configuration Commands
  # ============================================================================

  @happy
  Scenario: Display current configuration
    Given the application is installed
    When I run "kariya config show"
    Then I should see current configuration
    And the exit code should be 0

  @happy
  Scenario: Set configuration value
    Given the application is installed
    When I run "kariya config set log.level debug"
    Then the configuration should be updated
    And I should see confirmation
    And the exit code should be 0

  # ============================================================================
  # Database Commands
  # ============================================================================

  @happy
  Scenario: Initialize database
    Given the application is installed
    And no database exists
    When I run "kariya init"
    Then the database should be created
    And migrations should be applied
    And the exit code should be 0

  @happy
  Scenario: Database status
    Given the application is installed
    When I run "kariya db status"
    Then I should see database status
    And the exit code should be 0

  # ============================================================================
  # Error Handling
  # ============================================================================

  @sad
  Scenario: Unknown command
    Given the application is installed
    When I run "kariya unknowncommand"
    Then I should see an error message about unknown command
    And the exit code should be non-zero

  @sad
  Scenario: Missing required argument
    Given the application is installed
    When I run "kariya add"
    Then I should see an error about missing description
    And the exit code should be non-zero

  @sad
  Scenario: Invalid flag
    Given the application is installed
    When I run "kariya --invalidflag"
    Then I should see an error about unknown flag
    And the exit code should be non-zero
