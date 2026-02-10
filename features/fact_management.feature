@facts
Feature: Manage Career Facts
  As a professional
  I want to manage my career facts
  So that I can curate impactful achievements for my CV

  Background:
    Given I am on the main menu

  # ============================================================================
  # Empty State
  # ============================================================================

  @happy @smoke
  Scenario: View empty fact list
    Given the database is empty
    When I select "fact_management" from the menu
    Then I should see "No facts"
    And I should see "n"
    And I should be able to go back to the menu

  # ============================================================================
  # Fact List Display
  # ============================================================================

  @happy @smoke
  Scenario: View fact list with facts
    Given I have 5 facts in my profile
    When I select "fact_management" from the menu
    Then I should see a list of facts
    And I should see fact text
    And I should see strength signals
    And I should see categories

  @happy
  Scenario: Navigate through fact list with vim keys
    Given I have 5 facts in my profile
    When I select "fact_management" from the menu
    And I press "j" to navigate down
    And I press "k" to navigate up
    Then I should still be on the fact list

  @happy
  Scenario: Navigate through fact list with arrow keys
    Given I have 5 facts in my profile
    When I select "fact_management" from the menu
    And I press down arrow
    And I press up arrow
    Then I should still be on the fact list

  @happy
  Scenario: Page through fact list
    Given I have 20 facts in my profile
    When I select "fact_management" from the menu
    And I press page down
    Then I should see different facts
    When I press page up
    Then I should see the original facts

  @happy
  Scenario: Jump to first and last fact
    Given I have 20 facts in my profile
    When I select "fact_management" from the menu
    And I press "G" to go to last
    Then I should be at the last fact
    When I press "g" to go to first
    Then I should be at the first fact

  @happy
  Scenario: Refresh fact list
    Given I have 3 facts in my profile
    When I select "fact_management" from the menu
    And I press "r" to refresh
    Then I should still be on the fact list
    And the facts should be reloaded

  # ============================================================================
  # Fact Detail View
  # ============================================================================

  @happy
  Scenario: View fact details
    Given I have a fact "Led team of 5 engineers to deliver project on time"
    When I select "fact_management" from the menu
    And I press enter to view details
    Then I should see the fact detail view
    And I should see "Led team of 5 engineers"
    And I should see competency categories
    And I should see role fit
    And I should see audience relevance

  @happy
  Scenario: Close fact detail view
    Given I have a fact "Led team of 5 engineers to deliver project on time"
    When I select "fact_management" from the menu
    And I press enter to view details
    And I press escape
    Then I should still be on the fact list

  # ============================================================================
  # Create New Fact
  # ============================================================================

  @happy @wip
  Scenario: Open new fact form
    Given the database is empty
    When I select "fact_management" from the menu
    And I press "n" to create new fact
    Then I should see the fact editor form
    And I should see "Fact Text"
    And I should see "Competency Categories"
    And I should see "Role Fit"
    And I should see "Audience Relevance"

  @happy
  Scenario: Cancel new fact
    Given the database is empty
    When I select "fact_management" from the menu
    And I press "n" to create new fact
    And I press escape
    Then I should still be on the fact list
    And there should be 0 facts

  @happy @wip
  Scenario: Create new fact with all fields
    Given the database is empty
    When I select "fact_management" from the menu
    And I press "n" to create new fact
    And I enter fact text "Reduced deployment time by 50% through CI/CD automation"
    And I tab to competency categories
    And I select competency category "Technical"
    And I tab to role fit
    And I select role fit "Senior IC"
    And I tab to audience relevance
    And I select audience "Hiring Manager"
    And I submit the fact form
    Then I should still be on the fact list
    And there should be 1 fact
    And the fact should have text "Reduced deployment time by 50%"

  @happy @wip
  Scenario: Create fact with multiple competency categories
    Given the database is empty
    When I select "fact_management" from the menu
    And I press "n" to create new fact
    And I enter fact text "Led cross-functional team to deliver microservices architecture"
    And I tab to competency categories
    And I select competency category "Technical"
    And I select competency category "Leadership"
    And I submit the fact form
    Then the fact should have categories "Technical,Leadership"

  @happy @wip
  Scenario: Create fact with multiple audience types
    Given the database is empty
    When I select "fact_management" from the menu
    And I press "n" to create new fact
    And I enter fact text "Mentored 3 junior developers to senior level"
    And I tab to audience relevance
    And I select audience "Hiring Manager"
    And I select audience "Recruiter"
    And I submit the fact form
    Then the fact should have audiences "Hiring Manager,Recruiter"

  # ============================================================================
  # Edit Fact
  # ============================================================================

  @happy
  Scenario: Open edit fact form from list
    Given I have a fact "Original fact text here"
    When I select "fact_management" from the menu
    And I press "e" to edit
    Then I should see the fact editor form
    And I should see "Original fact text"

  @happy
  Scenario: Open edit fact form from detail
    Given I have a fact "Original fact text here"
    When I select "fact_management" from the menu
    And I press enter to view details
    And I press "e" to edit
    Then I should see the fact editor form

  @happy
  Scenario: Cancel edit fact returns to list
    Given I have a fact "Original fact text here"
    When I select "fact_management" from the menu
    And I press "e" to edit
    And I press escape
    Then I should still be on the fact list
    And the fact should have text "Original fact text"

  @happy
  Scenario: Cancel edit fact from detail returns to detail
    Given I have a fact "Original fact text here"
    When I select "fact_management" from the menu
    And I press enter to view details
    And I press "e" to edit
    And I press escape
    Then I should see the fact detail view

  @happy @wip
  Scenario: Edit fact text and save
    Given I have a fact "Original fact text here"
    When I select "fact_management" from the menu
    And I press "e" to edit
    And I clear the fact text field
    And I enter fact text "Updated fact text here"
    And I submit the fact form
    Then I should still be on the fact list
    And the fact should have text "Updated fact text"

  @happy @wip
  Scenario: Edit fact competency categories
    Given I have a fact with category "Technical"
    When I select "fact_management" from the menu
    And I press "e" to edit
    And I tab to competency categories
    And I deselect competency category "Technical"
    And I select competency category "Leadership"
    And I submit the fact form
    Then the fact should have categories "Leadership"

  @happy @wip
  Scenario: Navigate form with tab and shift-tab
    Given the database is empty
    When I select "fact_management" from the menu
    And I press "n" to create new fact
    And I press tab
    Then I should be on competency categories field
    When I press tab
    Then I should be on role fit field
    When I press shift-tab
    Then I should be on competency categories field

  # ============================================================================
  # Delete Fact
  # ============================================================================

  @happy
  Scenario: Open delete confirmation from list
    Given I have a fact "Fact to delete"
    When I select "fact_management" from the menu
    And I press "d" to delete
    Then I should see the delete confirmation

  @happy
  Scenario: Open delete confirmation from detail
    Given I have a fact "Fact to delete"
    When I select "fact_management" from the menu
    And I press enter to view details
    And I press "d" to delete
    Then I should see the delete confirmation

  @happy
  Scenario: Cancel delete fact
    Given I have a fact "Fact to delete"
    When I select "fact_management" from the menu
    And I press "d" to delete
    And I press "n" to cancel
    Then I should still be on the fact list
    And there should be 1 fact

  @happy
  Scenario: Cancel delete with escape
    Given I have a fact "Fact to delete"
    When I select "fact_management" from the menu
    And I press "d" to delete
    And I press escape
    Then I should still be on the fact list
    And there should be 1 fact

  @happy
  Scenario: Confirm delete fact with y
    Given I have a fact "Fact to delete"
    When I select "fact_management" from the menu
    And I press "d" to delete
    And I press "y" to confirm
    Then I should still be on the fact list
    And there should be 0 facts

  @happy
  Scenario: Confirm delete fact with enter
    Given I have a fact "Fact to delete"
    When I select "fact_management" from the menu
    And I press "d" to delete
    And I confirm the deletion
    Then I should still be on the fact list
    And there should be 0 facts

  # ============================================================================
  # Validation
  # ============================================================================

  @sad
  Scenario: Fact text minimum length validation
    Given the database is empty
    When I select "fact_management" from the menu
    And I press "n" to create new fact
    And I enter fact text "Too short"
    And I submit the fact form
    Then I should see a validation error
    And I should see "10"

  @sad @wip
  Scenario: Fact text maximum length validation
    Given the database is empty
    When I select "fact_management" from the menu
    And I press "n" to create new fact
    And I enter fact text with 2001 characters
    And I submit the fact form
    Then I should see a validation error
    And I should see "2000"

  # ============================================================================
  # Help Toggle
  # ============================================================================

  @happy
  Scenario: Toggle help in list state
    Given I have 3 facts in my profile
    When I select "fact_management" from the menu
    And I press "?" to toggle help
    Then I should see help information
    And I should see available shortcuts

  @happy
  Scenario: Toggle help in detail state
    Given I have a fact "My fact"
    When I select "fact_management" from the menu
    And I press enter to view details
    And I press "?" to toggle help
    Then I should see help information

  @happy
  Scenario: Toggle help in editor state
    Given the database is empty
    When I select "fact_management" from the menu
    And I press "n" to create new fact
    And I press "?" to toggle help
    Then I should see help information

  # ============================================================================
  # Navigation and Exit
  # ============================================================================

  @sad
  Scenario: Exit fact list returns to menu
    Given the database is empty
    When I select "fact_management" from the menu
    And I press escape
    Then I should be on the main menu

  @sad
  Scenario: Full navigation escape path
    Given I have a fact "My fact"
    When I select "fact_management" from the menu
    And I press enter to view details
    And I press "e" to edit
    And I press escape
    Then I should see the fact detail view
    When I press escape
    Then I should still be on the fact list
    When I press escape
    Then I should be on the main menu
