@workflows
Feature: Chained Workflows
  As a user
  I want to navigate between different features seamlessly
  So that I can manage my career data efficiently across the application

  Background:
    Given I am on the main menu

  # ============================================================================
  # Browse After Data Population
  # ============================================================================

  @happy @smoke
  Scenario: Show events in Browse after populating data
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    Then I should see a list of events

  @happy
  Scenario: Show facts in FactManagement after populating data
    Given I have 3 events in my timeline
    And I have 2 facts
    When I select "fact_management" from the menu
    Then I should see a list of facts

  @happy
  Scenario: Show bursts in BurstManagement after populating data
    Given I have 5 events in my timeline
    And I have 2 bursts
    When I select "burst_management" from the menu
    Then I should see a list of bursts

  # ============================================================================
  # Multi-Intent Navigation
  # ============================================================================

  @happy
  Scenario: Navigate from Browse to Generate CV
    Given I have 5 events in my timeline
    And I have 2 bursts
    And I have 3 facts
    When I select "browse_timeline" from the menu
    And I press escape
    And I select "generate_cv" from the menu
    Then I should see one of:
      | Profile |
      | Select  |
      | CV      |

  @happy
  Scenario: Navigate through multiple intents in sequence
    Given I have 5 events in my timeline
    When I select "browse_timeline" from the menu
    And I press escape
    And I select "generate_cv" from the menu
    And I press escape
    And I select "configure_system" from the menu
    And I press escape
    Then I should be on the main menu

  @happy
  Scenario: Navigate from BurstManagement to FactManagement
    Given I have 5 events in my timeline
    And I have 2 bursts
    And I have 3 facts
    When I select "burst_management" from the menu
    And I press escape
    And I select "fact_management" from the menu
    Then I should see a list of facts

  # ============================================================================
  # Session Persistence Across Intents
  # ============================================================================

  @happy
  Scenario: Persist data across intent navigation
    Given I have 3 events in my timeline
    And I have 2 facts
    When I select "browse_timeline" from the menu
    And I press escape
    And I select "fact_management" from the menu
    And I press escape
    Then there should be 3 events
    And there should be 2 facts

  @happy @wip
  Scenario: Persist data across simulated restart and intent navigation
    Given I have 4 events in my timeline
    When I restart the application
    And I select "browse_timeline" from the menu
    Then I should see a list of events

  # ============================================================================
  # Data Visibility Across Intents
  # ============================================================================

  @happy
  Scenario: See events in Browse
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    Then I should see a list of events

  @happy
  Scenario: See facts in FactManagement and GenerateCV context
    Given I have 3 events in my timeline
    And I have 3 facts
    When I select "fact_management" from the menu
    And I press escape
    And I select "generate_cv" from the menu
    Then I should see one of:
      | Profile |
      | Select  |
      | CV      |

  # ============================================================================
  # Workflow: Browse then Configure
  # ============================================================================

  @happy
  Scenario: Browse data then configure system
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press escape
    And I select "configure_system" from the menu
    Then I should see one of:
      | System  |
      | Profile |
      | Domain  |
    When I press escape
    Then there should be 3 events

  # ============================================================================
  # Return to Menu Between Intents
  # ============================================================================

  @happy
  Scenario: Return to menu after each intent
    When I select "browse_timeline" from the menu
    And I press escape
    Then I should be on the main menu
    When I select "generate_cv" from the menu
    And I press escape
    Then I should be on the main menu
    When I select "configure_system" from the menu
    And I press escape
    Then I should be on the main menu
