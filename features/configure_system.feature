@configure
Feature: Configure System
  As a user
  I want to configure the application settings
  So that I can customise the app to my preferences

  Background:
    Given I am on the main menu

  @happy @smoke
  Scenario: Open settings modal with comma key
    When I press "," to open settings
    Then I should see the settings modal
    And I should see "System"
    And I should see "Profile"
    And I should see "Export"
    And I should see "UI"

  @happy
  Scenario: Settings modal shows system settings by default
    When I press "," to open settings
    Then I should see the settings modal
    And I should see "Log Level"

  @happy
  Scenario: Navigate between sections with j key
    When I press "," to open settings
    And I navigate to the "Profile" section
    Then I should see "Full Name"

  @happy
  Scenario: Cancel settings modal with Escape
    When I press "," to open settings
    Then I should see the settings modal
    When I press escape
    Then I should be on the main menu

  @happy
  Scenario: Settings modal shows export settings
    When I press "," to open settings
    And I navigate to the "Export" section
    Then I should see "Default Destination"

  @happy
  Scenario: Settings modal shows UI settings
    When I press "," to open settings
    And I navigate to the "UI" section
    Then I should see "Theme"

  @happy
  Scenario: Save settings with Ctrl+S
    When I press "," to open settings
    And I press Ctrl+S
    Then I should be on the main menu

  @sad
  Scenario: Quit from settings modal returns to menu
    When I press "," to open settings
    Then I should see the settings modal
    When I press escape
    Then I should be on the main menu
