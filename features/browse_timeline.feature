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
    Given the database is empty
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

  @happy
  Scenario: Close event detail modal
    Given I have an event "Built REST API" at company "Acme Corp"
    When I select "browse_timeline" from the menu
    And I press enter to view details
    And I press escape
    Then I should still be on the timeline

  @happy
  Scenario: View and close details multiple times
    Given I have an event "Built REST API" at company "Acme Corp"
    When I select "browse_timeline" from the menu
    And I press enter to view details
    And I press escape
    And I press enter to view details
    And I press escape
    Then I should still be on the timeline

  # ============================================================================
  # Search Functionality
  # ============================================================================

  @happy
  Scenario: Open search modal
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "/" to search
    Then I should see the search modal

  @happy
  Scenario: Cancel search modal
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "/" to search
    And I press escape
    Then I should still be on the timeline

  @happy
  Scenario: Search for events by text
    Given I have an event "Built REST API" at company "Acme Corp"
    And I have an event "Deployed Kubernetes" at company "CloudCo"
    When I select "browse_timeline" from the menu
    And I press "/" to search
    And I type "REST" in the search
    And I submit the search
    Then I should see "Built REST API"
    And I should not see "Kubernetes"

  # ============================================================================
  # Filter Functionality
  # ============================================================================

  @happy
  Scenario: Open filter modal
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    Then I should see the filter modal

  @happy
  Scenario: Cancel filter modal
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    And I press escape
    Then I should still be on the timeline

  @happy
  Scenario: Filter modal shows company options
    Given I have an event "Built REST API" at company "Acme Corp"
    And I have an event "Deployed Kubernetes" at company "CloudCo"
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    Then I should see "Acme Corp"
    And I should see "CloudCo"
    When I press escape
    Then I should still be on the timeline

  @happy
  Scenario: Filter modal shows sort options
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    Then I should see "Sort By"
    And I should see "Sort Order"

  # ============================================================================
  # Sort Functionality
  # ============================================================================

  @happy
  Scenario: Open sort modal
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "s" to sort
    Then I should see the sort modal

  @happy
  Scenario: Cancel sort modal
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "s" to sort
    And I press escape
    Then I should still be on the timeline

  # ============================================================================
  # Quick Add Event
  # ============================================================================

  @happy
  Scenario: Open quick add modal
    Given the database is empty
    When I select "browse_timeline" from the menu
    And I press "a" to add event
    Then I should see the add event form

  @happy
  Scenario: Cancel quick add
    Given the database is empty
    When I select "browse_timeline" from the menu
    And I press "a" to add event
    And I press escape
    Then I should still be on the timeline
    And there should be 0 events

  @happy
  Scenario: Add event from timeline
    Given the database is empty
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
  Scenario: Cancel edit event
    Given I have an event "Original text" at company "Acme Corp"
    When I select "browse_timeline" from the menu
    And I press "e" to edit
    And I press escape
    Then I should still be on the timeline
    And the event should have description "Original text"

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
  Scenario: Cancel delete event
    Given I have an event "Event to delete" at company "Acme Corp"
    When I select "browse_timeline" from the menu
    And I press "d" to delete
    And I cancel the confirmation
    Then I should still be on the timeline
    And there should be 1 event

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

  @sad
  Scenario: Exit timeline returns to menu
    Given I have 3 events in my timeline
    When I select "browse_timeline" from the menu
    And I press escape
    Then I should be on the main menu

  @sad
  Scenario: Go back from empty timeline
    Given the database is empty
    When I select "browse_timeline" from the menu
    And I press escape
    Then I should be on the main menu

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

  @happy
  Scenario: Close skills modal returns to event detail
    Given I have an event "Built API with Go" at company "TechCorp"
    And the event has skills "Go,PostgreSQL"
    When I select "browse_timeline" from the menu
    And I press enter to view details
    And I press "s" to view skills
    And I press escape
    Then I should see the event detail modal

  # ============================================================================
  # Advanced Filtering
  # ============================================================================

  @happy
  Scenario: Filter by multiple companies
    Given I have an event "Event at Acme" at company "Acme Corp"
    And I have an event "Event at Tech" at company "TechCorp"
    And I have an event "Event at Other" at company "OtherCo"
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    And I select companies "Acme Corp,TechCorp"
    And I confirm filter
    Then I should see 2 events
    And I should see "Event at Acme"
    And I should see "Event at Tech"
    And I should not see "Event at Other"

  @happy
  Scenario: Filter by category
    Given I have an event "Technical event" with category "technical"
    And I have an event "Leadership event" with category "leadership"
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    And I select category "technical"
    And I confirm filter
    Then I should see 1 event
    And I should see "Technical event"

  @happy
  Scenario: Filter by project
    Given I have an event "Project A work" with project "Project Alpha"
    And I have an event "Project B work" with project "Project Beta"
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    And I select project "Project Alpha"
    And I confirm filter
    Then I should see 1 event
    And I should see "Project A work"

  @happy
  Scenario: Filter by date range
    Given I have an event "Old event" dated "2024-01-01"
    And I have an event "Recent event" dated "2025-06-01"
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    And I set date from "2025-01-01"
    And I confirm filter
    Then I should see 1 event
    And I should see "Recent event"

  @happy
  Scenario: Combine company and category filter
    Given I have an event "Acme tech" at company "Acme Corp" with category "technical"
    And I have an event "Acme lead" at company "Acme Corp" with category "leadership"
    And I have an event "Tech tech" at company "TechCorp" with category "technical"
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    And I select company "Acme Corp"
    And I select category "technical"
    And I confirm filter
    Then I should see 1 event
    And I should see "Acme tech"

  @happy
  Scenario: Clear all filters
    Given I have 5 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    And I select company "Acme Corp"
    And I confirm filter
    And I press "x" to clear filters
    Then I should see 5 events

  # ============================================================================
  # Search with Sort by Relevance
  # ============================================================================

  @happy
  Scenario: Sort search results by relevance
    Given I have an event "Go programming basics"
    And I have an event "Advanced Go techniques and patterns"
    And I have an event "Python programming"
    When I select "browse_timeline" from the menu
    And I press "/" to search
    And I enter "Go" as search text
    And I press "s" to sort
    And I select sort by "relevance"
    And I confirm sort
    Then I should see "Advanced Go techniques" before "Go programming basics"

  # ============================================================================
  # Filter Display
  # ============================================================================

  @happy
  Scenario: Active filters are displayed
    Given I have 5 events in my timeline
    When I select "browse_timeline" from the menu
    And I press "f" to filter
    And I select company "Acme Corp"
    And I confirm filter
    Then I should see "Acme Corp" as an active filter indicator
