@bursts
Feature: Manage Career Bursts
  As a professional
  I want to manage my career bursts
  So that I can organize related events into meaningful career narratives

  Background:
    Given I am on the main menu

  # ============================================================================
  # Empty State
  # ============================================================================

  @happy @smoke
  Scenario: View empty burst list
    Given the database is empty
    When I select "burst_management" from the menu
    Then I should see "No bursts"
    And I should be able to go back to the menu

  # ============================================================================
  # Burst List Display
  # ============================================================================

  @happy @smoke
  Scenario: View burst list with bursts
    Given I have 3 bursts in my profile
    When I select "burst_management" from the menu
    Then I should see a list of bursts
    And I should see "3"

  @happy
  Scenario: Navigate through burst list with vim keys
    Given I have 5 bursts in my profile
    When I select "burst_management" from the menu
    And I press "j" to navigate down
    And I press "k" to navigate up
    Then I should still be on the burst list

  @happy
  Scenario: Navigate through burst list with arrow keys
    Given I have 5 bursts in my profile
    When I select "burst_management" from the menu
    And I press down arrow
    And I press up arrow
    Then I should still be on the burst list

  @happy
  Scenario: Page through burst list
    Given I have 20 bursts in my profile
    When I select "burst_management" from the menu
    And I press page down
    Then I should see different bursts
    When I press page up
    Then I should see the original bursts

  @happy
  Scenario: Jump to first and last burst
    Given I have 20 bursts in my profile
    When I select "burst_management" from the menu
    And I press "G" to go to last
    Then I should be at the last burst
    When I press "g" to go to first
    Then I should be at the first burst

  # ============================================================================
  # Burst Detail Modal
  # ============================================================================

  @happy
  Scenario: View burst details
    Given I have a burst "Backend API Development" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    Then I should see the burst detail modal
    And I should see "Backend API Development"
    And I should see "3"

  @happy
  Scenario: Close burst detail modal
    Given I have a burst "Backend API Development" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press escape
    Then I should still be on the burst list

  @happy
  Scenario: View burst events from detail modal
    Given I have a burst "Backend API Development" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "v" to view events
    Then I should see the burst events modal
    And I should see event details

  @happy
  Scenario: Close burst events modal
    Given I have a burst "Backend API Development" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "v" to view events
    And I press escape
    Then I should see the burst detail modal

  @happy @wip
  Scenario: View burst facts from detail modal
    Given I have a confirmed burst "Backend API Development" with facts
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "f" to view facts
    Then I should see the burst facts modal

  @happy @wip
  Scenario: Close burst facts modal
    Given I have a confirmed burst "Backend API Development" with facts
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "f" to view facts
    And I press escape
    Then I should see the burst detail modal

  @happy @wip
  Scenario: View burst skills from detail modal
    Given I have a confirmed burst "Backend API Development" with skills
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "s" to view skills
    Then I should see the burst skills modal

  @happy @wip
  Scenario: Close burst skills modal
    Given I have a confirmed burst "Backend API Development" with skills
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "s" to view skills
    And I press escape
    Then I should see the burst detail modal

  # ============================================================================
  # Edit Burst
  # ============================================================================

  @happy
  Scenario: Open edit burst modal from list
    Given I have a burst "Original Name" with 2 events
    When I select "burst_management" from the menu
    And I press "e" to edit
    Then I should see the edit burst form
    And I should see "Original Name"

  @happy
  Scenario: Open edit burst modal from detail
    Given I have a burst "Original Name" with 2 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "e" to edit
    Then I should see the edit burst form

  @happy
  Scenario: Cancel edit burst
    Given I have a burst "Original Name" with 2 events
    When I select "burst_management" from the menu
    And I press "e" to edit
    And I press escape
    Then I should still be on the burst list
    And the burst should have name "Original Name"

  @happy
  Scenario: Edit burst name and save
    Given I have a burst "Original Name" with 2 events
    When I select "burst_management" from the menu
    And I press "e" to edit
    And I clear the burst name field
    And I enter burst name "Updated Name"
    And I submit the burst form
    Then I should still be on the burst list
    And the burst should have name "Updated Name"

  @happy @wip
  Scenario: Edit burst description and save
    Given I have a burst "My Burst" with 2 events
    When I select "burst_management" from the menu
    And I press "e" to edit
    And I tab to description field
    And I enter burst description "A detailed description"
    And I submit the burst form
    Then the burst should have description "A detailed description"

  # ============================================================================
  # Delete Burst
  # ============================================================================

  @happy
  Scenario: Open delete confirmation from list
    Given I have a burst "Burst to Delete" with 2 events
    When I select "burst_management" from the menu
    And I press "d" to delete
    Then I should see the delete confirmation

  @happy
  Scenario: Open delete confirmation from detail
    Given I have a burst "Burst to Delete" with 2 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "d" to delete
    Then I should see the delete confirmation

  @happy
  Scenario: Cancel delete burst
    Given I have a burst "Burst to Delete" with 2 events
    When I select "burst_management" from the menu
    And I press "d" to delete
    And I cancel the confirmation
    Then I should still be on the burst list
    And there should be 1 burst

  @happy
  Scenario: Confirm delete burst
    Given I have a burst "Burst to Delete" with 2 events
    When I select "burst_management" from the menu
    And I press "d" to delete
    And I confirm the deletion
    Then I should still be on the burst list
    And there should be 0 bursts

  # ============================================================================
  # Confirm Burst (Trigger Fact Extraction)
  # ============================================================================

  @happy
  Scenario: Confirm unconfirmed burst
    Given I have an unconfirmed burst "New Burst" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "c" to confirm
    Then I should see the confirm burst modal
    And I should see "confirm"

  @happy
  Scenario: Cancel confirm burst
    Given I have an unconfirmed burst "New Burst" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "c" to confirm
    And I press escape
    Then I should see the burst detail modal
    And the burst should not be confirmed

  @happy @wip
  Scenario: Confirm burst triggers fact extraction
    Given I have an unconfirmed burst "New Burst" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "c" to confirm
    And I confirm the action
    Then I should see the loading modal
    And I should see "Extracting"

  @happy @wip
  Scenario: Re-confirm burst shows re-extraction prompt
    Given I have a confirmed burst "Confirmed Burst" with facts
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "c" to confirm
    Then I should see "re-extract"

  # ============================================================================
  # Burst Suggestion (AI Detection)
  # ============================================================================

  @happy @wip
  Scenario: Trigger burst suggestion
    Given I have 5 unassigned events
    When I select "burst_management" from the menu
    And I press "s" to suggest bursts
    Then I should see the loading modal
    And I should see "Detecting"

  @happy @wip
  Scenario: Cancel burst suggestion
    Given I have 5 unassigned events
    When I select "burst_management" from the menu
    And I press "s" to suggest bursts
    And I press escape
    Then I should still be on the burst list

  @happy @wip
  Scenario: View burst suggestions after detection
    Given I have 5 unassigned events
    When I select "burst_management" from the menu
    And I press "s" to suggest bursts
    And the detection completes
    Then I should see the burst suggestion modal
    And I should see suggested burst names
    And I should see confidence scores

  @happy @wip
  Scenario: Navigate through burst suggestions
    Given I have burst suggestions available
    When I am on the burst suggestion modal
    And I press "j" to navigate down
    And I press "k" to navigate up
    Then I should see different suggestions highlighted

  @happy @wip
  Scenario: View events for burst suggestion
    Given I have burst suggestions available
    When I am on the burst suggestion modal
    And I press enter to view events
    Then I should see the suggestion events modal
    And I should see event details

  @happy @wip
  Scenario: Accept burst suggestion
    Given I have burst suggestions available
    When I am on the burst suggestion modal
    And I press "a" to accept
    Then I should see success message
    And there should be 1 burst

  @happy @wip
  Scenario: Reject burst suggestion
    Given I have burst suggestions available
    When I am on the burst suggestion modal
    And I press "r" to reject
    Then the suggestion should be marked as rejected
    And there should be 0 bursts

  @happy @wip
  Scenario: Cancel burst suggestion review
    Given I have burst suggestions available
    When I am on the burst suggestion modal
    And I press escape
    Then I should still be on the burst list
    And there should be 0 bursts

  # ============================================================================
  # Skill Inference from Burst
  # ============================================================================

  @happy @wip
  Scenario: Trigger skill inference from confirmed burst
    Given I have a confirmed burst "Backend Development" with 5 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "i" to infer skills
    Then I should see the loading modal
    And I should see "Inferring"

  @happy @wip
  Scenario: Skill inference not available for unconfirmed burst
    Given I have an unconfirmed burst "New Burst" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "i" to infer skills
    Then I should still be on the burst detail modal
    And I should not see the loading modal

  @happy @wip
  Scenario: Cancel skill inference
    Given I have a confirmed burst "Backend Development" with 5 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "i" to infer skills
    And I press escape
    Then I should see the burst detail modal

  @happy @wip
  Scenario: View skill suggestions after inference
    Given I have a confirmed burst "Backend Development" with 5 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "i" to infer skills
    And the inference completes
    Then I should see the skill suggestion modal
    And I should see suggested skills
    And I should see skill categories

  @happy @wip
  Scenario: Navigate through skill suggestions
    Given I have skill suggestions from burst
    When I am on the skill suggestion modal
    And I press "j" to navigate down
    And I press "k" to navigate up
    Then I should see different skills highlighted

  @happy @wip
  Scenario: View events for skill suggestion
    Given I have skill suggestions from burst
    When I am on the skill suggestion modal
    And I press enter to view events
    Then I should see events that led to this skill

  @happy @wip
  Scenario: Accept skill suggestion
    Given I have skill suggestions from burst
    When I am on the skill suggestion modal
    And I press "a" to accept
    Then I should see success message
    And there should be 1 skill

  @happy @wip
  Scenario: Reject skill suggestion
    Given I have skill suggestions from burst
    When I am on the skill suggestion modal
    And I press "r" to reject
    Then the skill should be marked as rejected

  # ============================================================================
  # Navigation and Exit
  # ============================================================================

  @sad
  Scenario: Exit burst list returns to menu
    Given the database is empty
    When I select "burst_management" from the menu
    And I press escape
    Then I should be on the main menu

  @sad @wip
  Scenario: Navigate back through modal stack
    Given I have a burst "My Burst" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "v" to view events
    And I press escape
    Then I should see the burst detail modal
    When I press escape
    Then I should still be on the burst list
    When I press escape
    Then I should be on the main menu
