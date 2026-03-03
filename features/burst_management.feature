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
     Given I have no data
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
  Scenario: View burst events from detail modal
    Given I have a burst "Backend API Development" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "v" to view events
    Then I should see the burst events modal
    And I should see event details

  @happy
  Scenario: View burst facts from detail modal
    Given I have a confirmed burst "Backend API Development" with facts
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "f" to view facts
    Then I should see the burst facts modal

  @happy
  Scenario: View burst skills from detail modal
    Given I have a confirmed burst "Backend API Development" with skills
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "s" to view skills
    Then I should see the burst skills modal

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
  Scenario: Edit burst name and save
    Given I have a burst "Original Name" with 2 events
    When I select "burst_management" from the menu
    And I press "e" to edit
    And I clear the burst name field
    And I enter burst name "Updated Name"
    And I submit the burst form
    Then I should still be on the burst list
    And the burst should have name "Updated Name"

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
  Scenario: Confirm burst triggers fact extraction
    Given I have an unconfirmed burst "New Burst" with 3 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "c" to confirm
    And I confirm the action
    Then I should still be on the burst list
    And the burst should be confirmed

  @happy
  Scenario: Re-confirm burst shows re-extraction prompt
    Given I have a confirmed burst "Confirmed Burst" with facts
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "c" to confirm
    Then I should see the confirm burst modal

  # ============================================================================
  # Burst Suggestion (AI Detection)
  # ============================================================================

  @happy
  Scenario: Trigger burst suggestion
    Given I have 5 unassigned events
    When I select "burst_management" from the menu
    And I press "s" to suggest bursts
    And the detection completes
    Then I should see the burst suggestion modal

  @happy
  Scenario: View burst suggestions after detection
    Given I have 5 unassigned events
    When I select "burst_management" from the menu
    And I press "s" to suggest bursts
    And the detection completes
    Then I should see the burst suggestion modal
    And I should see suggested burst names
     And I should see confidence scores

   @happy
   Scenario: View events for burst suggestion
    Given I have burst suggestions available
    When I am on the burst suggestion modal
    And I press enter to view events
    Then I should see the suggestion events modal
    And I should see event details

  @happy
  Scenario: Accept burst suggestion
    Given I have burst suggestions available
    When I am on the burst suggestion modal
    And I press "a" to accept
    Then I should see the skill suggestions modal
    And I close the modal
    And there should be 1 burst

  @happy
  Scenario: Reject burst suggestion
    Given I have burst suggestions available
    When I am on the burst suggestion modal
    And I press "r" to reject
    Then I should still be on the burst list
    And there should be 0 bursts

  # ============================================================================
  # Skill Inference from Burst
  # ============================================================================

  @happy
  Scenario: Trigger skill inference from confirmed burst
    Given I have a confirmed burst "Backend Development" with 5 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "i" to infer skills
    And the inference completes
    Then I should see the skill suggestion modal

  @happy
  Scenario: View skill suggestions after inference
    Given I have a confirmed burst "Backend Development" with 5 events
    When I select "burst_management" from the menu
    And I press enter to view details
    And I press "i" to infer skills
    And the inference completes
    Then I should see the skill suggestion modal
    And I should see suggested skills
    And I should see skill categories

   @happy
  Scenario: View events for skill suggestion
    Given I have skill suggestions from burst
    When I am on the skill suggestion modal
    And I press enter to view events
    Then I should see events that led to this skill

  @happy
  Scenario: Accept skill suggestion
    Given I have skill suggestions from burst
    When I am on the skill suggestion modal
    And I press "a" to accept
    Then I should still be on the skill suggestions modal
    And there should be 1 skill

  @happy
  Scenario: Reject skill suggestion
    Given I have skill suggestions from burst
    When I am on the skill suggestion modal
    And I press "r" to reject
    Then the skill should be marked as rejected
