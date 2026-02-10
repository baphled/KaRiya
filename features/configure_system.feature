@configure
Feature: Configure System
  As a user
  I want to configure the application settings
  So that I can customize the app to my preferences

  Background:
    Given I am on the main menu

  # ============================================================================
  # Domain Selection
  # ============================================================================

  @happy @smoke
  Scenario: View configuration domains
    When I select "configure_system" from the menu
    Then I should see the domain selection screen
    And I should see "System"
    And I should see "Profile"
    And I should see "Export"
    And I should see "UI"

  @happy
  Scenario: Navigate domains with vim keys
    When I select "configure_system" from the menu
    And I press "j" to navigate down
    And I press "k" to navigate up
    Then I should still be on domain selection

  @happy
  Scenario: Navigate domains with arrow keys
    When I select "configure_system" from the menu
    And I press down arrow
    And I press up arrow
    Then I should still be on domain selection

  @happy
  Scenario: Jump to first and last domain
    When I select "configure_system" from the menu
    And I press "G" to go to last
    Then I should be at the last domain
    When I press "g" to go to first
    Then I should be at the first domain

  @sad
  Scenario: Exit domain selection returns to menu
    When I select "configure_system" from the menu
    And I press escape
    Then I should be on the main menu

  @sad
  Scenario: Quit from domain selection
    When I select "configure_system" from the menu
    And I press "q" to quit
    Then I should be on the main menu

  # ============================================================================
  # System Domain Settings
  # ============================================================================

  @happy @wip
  Scenario: Open system settings
    When I select "configure_system" from the menu
    And I select "System" domain
    Then I should see the edit settings modal
    And I should see "Log Level"
    And I should see "Data Directory"
    And I should see "Auto Backup"
    And I should see "Backup Count"

  @happy
  Scenario: Edit log level setting
    When I select "configure_system" from the menu
    And I select "System" domain
    And I navigate to "Log Level" field
    Then I should see "debug"
    And I should see "info"
    And I should see "warn"
    And I should see "error"

  @happy
  Scenario: Edit data directory setting
    When I select "configure_system" from the menu
    And I select "System" domain
    And I navigate to "Data Directory" field
    And I enter "/custom/data/path"
    Then the field should show "/custom/data/path"

  @happy
  Scenario: Toggle auto backup setting
    When I select "configure_system" from the menu
    And I select "System" domain
    And I navigate to "Auto Backup" field
    And I toggle the boolean value
    Then the value should change

  @happy
  Scenario: Edit backup count setting
    When I select "configure_system" from the menu
    And I select "System" domain
    And I navigate to "Backup Count" field
    And I enter "10"
    Then the field should show "10"

  # ============================================================================
  # Profile Domain Settings
  # ============================================================================

  @happy @wip
  Scenario: Open profile settings
    When I select "configure_system" from the menu
    And I select "Profile" domain
    Then I should see the edit settings modal
    And I should see "Full Name"
    And I should see "Email"
    And I should see "Professional Title"
    And I should see "Location"

  @happy @wip
  Scenario: View all profile fields
    When I select "configure_system" from the menu
    And I select "Profile" domain
    Then I should see "GitHub URL"
    And I should see "Portfolio URL"
    And I should see "Programming Languages"
    And I should see "Frontend Technologies"
    And I should see "Core Strengths"

  @happy
  Scenario: Edit name and email
    When I select "configure_system" from the menu
    And I select "Profile" domain
    And I enter name "John Doe"
    And I tab to email
    And I enter email "john@example.com"
    And I submit settings
    Then I should see the review modal

  @happy
  Scenario: Edit list fields with comma-separated values
    When I select "configure_system" from the menu
    And I select "Profile" domain
    And I navigate to "Programming Languages" field
    And I enter "Go,Python,TypeScript"
    Then the field should accept comma-separated values

  @happy
  Scenario: Select default role
    When I select "configure_system" from the menu
    And I select "Profile" domain
    And I navigate to "Default Role" field
    Then I should see "Junior IC"
    And I should see "Senior IC"
    And I should see "Staff IC"
    And I should see "Principal IC"
    And I should see "Manager"

  @happy
  Scenario: Select default audience
    When I select "configure_system" from the menu
    And I select "Profile" domain
    And I navigate to "Default Audience" field
    Then I should see "Technical"
    And I should see "Executive"
    And I should see "General"

  # ============================================================================
  # Export Domain Settings
  # ============================================================================

  @happy
  Scenario: Open export settings
    When I select "configure_system" from the menu
    And I select "Export" domain
    Then I should see the edit settings modal
    And I should see "Default Destination"
    And I should see "Auto-Open"

  @happy @wip
  Scenario: Select default destination
    When I select "configure_system" from the menu
    And I select "Export" domain
    And I navigate to "Default Destination" field
    Then I should see "File"
    And I should see "Clipboard"

  @happy
  Scenario: Toggle auto-open setting
    When I select "configure_system" from the menu
    And I select "Export" domain
    And I navigate to "Auto-Open" field
    And I toggle the boolean value
    Then the value should change

  # ============================================================================
  # UI Domain Settings
  # ============================================================================

  @happy
  Scenario: Open UI settings
    When I select "configure_system" from the menu
    And I select "UI" domain
    Then I should see the edit settings modal
    And I should see "Theme"
    And I should see "Animations"

  @happy @wip
  Scenario: Select theme
    When I select "configure_system" from the menu
    And I select "UI" domain
    And I navigate to "Theme" field
    Then I should see "Light"
    And I should see "Dark"

  @happy
  Scenario: Toggle animations setting
    When I select "configure_system" from the menu
    And I select "UI" domain
    And I navigate to "Animations" field
    And I toggle the boolean value
    Then the value should change

  # ============================================================================
  # Edit Settings Modal
  # ============================================================================

  @happy
  Scenario: Cancel edit settings
    When I select "configure_system" from the menu
    And I select "System" domain
    And I press escape
    Then I should be on domain selection

  @happy
  Scenario: Navigate form with tab
    When I select "configure_system" from the menu
    And I select "System" domain
    And I press tab
    Then I should move to the next field
    When I press shift-tab
    Then I should move to the previous field

  @happy
  Scenario: Submit settings with Ctrl+S
    When I select "configure_system" from the menu
    And I select "System" domain
    And I make a change
    And I press Ctrl+S
    Then I should see the review modal

  @happy @wip
  Scenario: Submit settings with Enter
    When I select "configure_system" from the menu
    And I select "System" domain
    And I make a change
    And I complete the form
    Then I should see the review modal

  # ============================================================================
  # Review Changes Modal
  # ============================================================================

  @happy @wip
  Scenario: View pending changes
    When I select "configure_system" from the menu
    And I select "System" domain
    And I change log level to "debug"
    And I submit settings
    Then I should see the review modal
    And I should see "Changes to apply"
    And I should see "log_level"
    And I should see "debug"

  @happy
  Scenario: Confirm changes from review
    When I am on the review changes modal
    And I press enter to confirm
    Then I should see the confirm modal

  @happy
  Scenario: Confirm changes with y key
    When I am on the review changes modal
    And I press "y" to confirm
    Then I should see the confirm modal

  @happy
  Scenario: Cancel review returns to edit
    When I am on the review changes modal
    And I press escape
    Then I should see the edit settings modal

  @happy
  Scenario: Cancel review with n key
    When I am on the review changes modal
    And I press "n" to cancel
    Then I should see the edit settings modal

  @happy
  Scenario: Cancel review with q key
    When I am on the review changes modal
    And I press "q" to cancel
    Then I should see the edit settings modal

  # ============================================================================
  # Confirm Modal
  # ============================================================================

  @happy
  Scenario: View confirm dialog
    When I am on the confirm modal
    Then I should see "Confirm Changes"
    And I should see "Are you sure"

  @happy
  Scenario: Confirm save with enter
    When I am on the confirm modal
    And I press enter to confirm
    Then I should see the saving modal

  @happy
  Scenario: Confirm save with y key
    When I am on the confirm modal
    And I press "y" to confirm
    Then I should see the saving modal

  @happy
  Scenario: Cancel confirm returns to review
    When I am on the confirm modal
    And I press escape
    Then I should see the review modal

  @happy
  Scenario: Cancel confirm with n key
    When I am on the confirm modal
    And I press "n" to cancel
    Then I should see the review modal

  # ============================================================================
  # Saving and Result
  # ============================================================================

  @happy @wip
  Scenario: View saving progress
    When I confirm save
    Then I should see the saving modal
    And I should see "Saving"
    And I should see a spinner

  @happy
  Scenario: View success result
    When the save completes successfully
    Then I should see the success modal
    And I should see "Configuration saved"

  @happy
  Scenario: Success modal auto-dismisses
    When the save completes successfully
    Then the success modal should auto-dismiss
    And I should be on the main menu

  @happy
  Scenario: Dismiss success with enter
    When the save completes successfully
    And I press enter
    Then I should be on the main menu

  @sad
  Scenario: View error result
    When the save fails
    Then I should see the error modal
    And I should see "Save Failed"
    And I should see error details

  @sad
  Scenario: Dismiss error returns to edit
    When the save fails
    And I dismiss the error modal
    Then I should see the edit settings modal

  # ============================================================================
  # Help Toggle
  # ============================================================================

  @happy
  Scenario: Toggle help in domain selection
    When I select "configure_system" from the menu
    And I press "?" to toggle help
    Then I should see help information

  @happy
  Scenario: Toggle help in edit settings
    When I select "configure_system" from the menu
    And I select "System" domain
    And I press "?" to toggle help
    Then I should see help information

  # ============================================================================
  # Full Workflow
  # ============================================================================

  @happy @wip
  Scenario: Complete configuration save workflow
    When I select "configure_system" from the menu
    And I select "System" domain
    And I change log level to "debug"
    And I submit settings
    And I confirm the review
    And I confirm the save
    And the save completes
    Then I should be on the main menu

  @happy @wip
  Scenario: Full escape navigation path
    When I select "configure_system" from the menu
    And I select "System" domain
    And I make a change
    And I submit settings
    And I confirm the review
    And I press escape from confirm
    Then I should see the review modal
    When I press escape
    Then I should see the edit settings modal
    When I press escape
    Then I should be on domain selection
    When I press escape
    Then I should be on the main menu

  # ============================================================================
  # Different Domains
  # ============================================================================

  @happy
  Scenario: Select each domain
    When I select "configure_system" from the menu
    And I select "System" domain
    And I press escape
    And I select "Profile" domain
    And I press escape
    And I select "Export" domain
    And I press escape
    And I select "UI" domain
    Then I should see the edit settings modal
