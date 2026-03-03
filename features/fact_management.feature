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
     Given I have no data
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


  # ============================================================================
  # Create New Fact
  # ============================================================================

   @happy
   Scenario: Open new fact form
     Given I have no data
     When I select "fact_management" from the menu
     And I press "n" to create new fact
     Then I should see the fact editor form
     And I should see "Fact Text"
     And I should see "Competency Categories"


   @happy @wip
   Scenario: Create new fact with all fields
     Given I have no data
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
    Given I have no data
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
    Given I have no data
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


  # ============================================================================
  # Help Toggle
  # ============================================================================


  # ============================================================================
  # Navigation and Exit
  # ============================================================================

