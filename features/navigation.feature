@navigation
Feature: Application Navigation
  As a user
  I want to navigate through the application
  So that I can access different features efficiently

  # ============================================================================
  # Main Menu Navigation
  # ============================================================================

  @happy @smoke
  Scenario: View main menu on startup
    Given I start the application
    Then I should see the main menu
    And I should see one of:
      | Browse Timeline |
      | Capture Event   |
      | Manage Skills   |

  @happy
  Scenario: Navigate menu with vim keys
    Given I am on the main menu
    When I press "j" to move down
    And I press "k" to move up
    Then I should still be on the main menu

  @happy
  Scenario: Navigate menu with arrow keys
    Given I am on the main menu
    When I press down arrow
    And I press up arrow
    Then I should still be on the main menu

  @happy
  Scenario: Select menu item with enter
    Given I am on the main menu
    When I select "browse_timeline" from the menu
    Then I should not be on the main menu

  # ============================================================================
  # Help System
  # ============================================================================

  @happy @smoke
  Scenario: View help from main menu
    Given I am on the main menu
    When I press "?" for help
    Then I should see help information
    And I should see "Keyboard Reference" or "Toggle this help"

  @happy
  Scenario: Close help modal
    Given I am on the main menu
    When I press "?" for help
    And I press escape
    Then I should be on the main menu

  @happy
  Scenario: Context-sensitive help in browse timeline
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "?" for help
    Then I should see "Move up" or "Navigate"
    And I should see "Enter" or "Confirm"

  # ============================================================================
  # Quit Application
  # ============================================================================

  @happy @smoke
  Scenario: Quit from main menu with q
    Given I am on the main menu
    When I press "q" to quit
    Then the application should exit

  @happy
  Scenario: Quit from main menu with Ctrl+C
    Given I am on the main menu
    When I press Ctrl+C
    Then the application should exit

  @sad
  Scenario: Cannot quit with q while in intent
    Given I am in the middle of capturing an event
    When I press "q" to quit
    Then the application should not exit

  # ============================================================================
  # Back Navigation
  # ============================================================================

  @happy @smoke
  Scenario: Go back from intent to menu
    Given I am on the main menu
    When I select "browse_timeline" from the menu
    And I press escape
    Then I should be on the main menu

  @happy
  Scenario: Go back from nested modal
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "/" to search
    And I press escape
    Then I should still be on the timeline
    When I press escape
    Then I should be on the main menu

  @happy
  Scenario: Multiple back navigations
    Given I have an event "Test event" at company "Test Co"
    When I select "browse_timeline" from the menu
    And I press enter to view details
    And I press escape
    And I press escape
    Then I should be on the main menu

  # ============================================================================
  # Keyboard Shortcuts
  # ============================================================================

  @happy
  Scenario: Global keyboard shortcuts work everywhere
    Given I am on the main menu
    Then pressing "?" should show help
    And pressing "q" should quit
    And pressing "Ctrl+C" should quit

  @happy
  Scenario: Intent-specific shortcuts
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    Then pressing "/" should open search
    And pressing "f" should open filter
    And pressing "s" should open sort

  # ============================================================================
  # Focus Management
  # ============================================================================

  @happy
  Scenario: Modal captures focus
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "/" to search
    Then typing should go to the search input
     And pressing escape should close the modal

