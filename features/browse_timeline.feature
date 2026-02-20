@browse
Feature: Browse Career Timeline
  As a professional
  I want to browse my career events
  So that I can review and manage my work history

  Background:
    Given I am on the main menu

  # ============================================================================
  # Empty State
  # ============================================================================

   @happy @smoke
   Scenario: View empty timeline
     Given I have no data
     When I select "browse_timeline" from the menu
     Then I should see "No events"
     And I should be able to go back to the menu

  # ============================================================================
  # Event List Display
  # ============================================================================

  @happy @smoke
  Scenario: View timeline with events
    Given I have 5 events in my timeline
    When I select "browse_timeline" from the menu
    Then I should see a list of events
    And I should see "5"

   # ============================================================================
   # Event Detail View
   # ============================================================================

  @happy
  Scenario: View event details
    Given I have an event "Built REST API" at company "Acme Corp"
    When I select "browse_timeline" from the menu
    And I press enter to view details
    Then I should see "Built REST API"
    And I should see "Acme Corp"



  # ============================================================================
  # Quick Add Event
  # ============================================================================

   @happy
   Scenario: Open quick add modal
     Given I have no data
     When I select "browse_timeline" from the menu
     And I press "a" to add event
     Then I should see the add event form

   @happy
   Scenario: Add event from timeline
     Given I have no data
     When I select "browse_timeline" from the menu
     And I press "a" to add event
     And I enter "New event from timeline" as description
     And I submit the form
     Then I should still be on the timeline
     And there should be 1 event

  # ============================================================================
  # Edit Event
  # ============================================================================

  @happy
  Scenario: Open edit modal for selected event
    Given I have an event "Original text" at company "Acme Corp"
    When I select "browse_timeline" from the menu
    And I press "e" to edit
    Then I should see the edit event form
    And I should see "Original text"

  @happy
  Scenario: Edit event metadata and save
    Given I have an event "Original text" at company "Acme Corp"
    When I select "browse_timeline" from the menu
    And I press "e" to edit
    And I navigate to company field
    And I clear the company field
    And I enter "New Company" as company
    And I submit the form
    Then I should still be on the timeline
    And the event should have company "New Company"

  # ============================================================================
  # Delete Event
  # ============================================================================

  @happy
  Scenario: Open delete confirmation
    Given I have an event "Event to delete" at company "Acme Corp"
    When I select "browse_timeline" from the menu
    And I press "d" to delete
    Then I should see the delete confirmation


  @happy
  Scenario: Confirm delete event
    Given I have an event "Event to delete" at company "Acme Corp"
    When I select "browse_timeline" from the menu
    And I press "d" to delete
    And I confirm the deletion
    Then I should still be on the timeline
    And there should be 0 events

  # ============================================================================
  # Navigation and Exit
  # ============================================================================


   # ============================================================================
   # View Skills from Event Detail
   # ============================================================================

  @happy
  Scenario: View skills from event detail modal
    Given I have an event "Built API with Go" at company "TechCorp"
    And the event has skills "Go,PostgreSQL,REST"
    When I select "browse_timeline" from the menu
    And I press enter to view details
    And I press "s" to view skills
    Then I should see the skills detail modal
    And I should see "Go"
    And I should see "PostgreSQL"
    And I should see "REST"


