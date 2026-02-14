@capture
Feature: Capture Career Events
  As a professional
  I want to capture my career events
  So that I can build my work history for CV generation

  Background:
    Given I have no data
    And I am on the main menu

  @happy @smoke
  Scenario: Quick capture with minimal input
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Led daily standup meetings for the platform team"
    And I submit the event
    Then I should see the success message
    And there should be 1 event
    And the event should have description "Led daily standup meetings for the platform team"

  @happy
  Scenario: Quick capture with company and project
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I capture an event:
      | field       | value                                    |
      | description | Built REST API with authentication       |
      | company     | Acme Corp                                |
      | project     | User Platform                            |
    And I submit the event
    Then I should see the success message
    And there should be 1 event
    And the event should have company "Acme Corp"
    And the event should have project "User Platform"

  @happy
  Scenario: Manual capture with full metadata
    When I select "capture_event" from the menu
    And I select manual capture strategy
    And I capture an event:
      | field       | value                                    |
      | description | Architected microservices migration      |
      | company     | Tech Innovations Ltd                     |
      | project     | Cloud Migration                          |
      | tags        | technical,project                        |
      | categories  | technical,architecture                   |
    And I submit the event
    Then I should see the success message
    And there should be 1 event
    And the event should have tags "technical,project"
    And the event should have categories "technical,architecture"

  @happy @enrichment
  Scenario: Accept suggested burst during review
    Given I have an event "Built user authentication service" at company "Acme Corp"
    Given I have an event "Implemented OAuth2 integration" at company "Acme Corp"
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Added JWT token refresh mechanism"
    And I set event company to "Acme Corp"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    And I should see "b" key badge for bursts
    When I accept the suggested burst
    And I confirm the review
    Then I should be on the main menu
    And there should be 1 burst

  @happy @enrichment
  Scenario: Accept inferred skill during review
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Built microservices in Go with PostgreSQL database"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    When I accept all inferred skills
    And I confirm the review
    Then I should be on the main menu
    And there should be skills including "Go"

  @happy @enrichment
  Scenario: Edit suggested burst before accepting
    Given I have an event "Deployed Kubernetes cluster" at company "CloudCo"
    Given I have an event "Set up CI/CD pipeline" at company "CloudCo"
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Implemented auto-scaling policies"
    And I set event company to "CloudCo"
    And I submit the event
    And I dismiss the success modal
    When I edit the suggested burst
    And I change burst name to "Infrastructure Automation Initiative"
    And I save the burst edit
    And I confirm the review
    Then there should be 1 burst with name "Infrastructure Automation Initiative"

  @happy @enrichment
  Scenario: Edit event metadata during review
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Initial event description"
    And I submit the event
    And I dismiss the success modal
    When I open the metadata editor
    And I change event company to "Updated Company"
    And I save metadata changes
    And I confirm the review
    Then the event should have company "Updated Company"

  @sad
  Scenario: Cancel at strategy selection returns to menu
    When I select "capture_event" from the menu
    And I cancel
    Then I should be on the main menu
    And there should be 0 events

  @sad
  Scenario: Cancel at form entry returns to strategy
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I cancel
    Then I should see the strategy selection
    And there should be 0 events

  @sad
  Scenario: Escape from form navigates back through states
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I cancel
    Then I should see the strategy selection
    When I cancel
    Then I should be on the main menu

  @sad @enrichment
  Scenario: Reject all suggestions and submit raw event
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Event with rejected enrichments"
    And I submit the event
    And I dismiss the success modal
    When I reject all suggestions
    And I confirm the review
    Then I should be on the main menu
    And there should be 1 event
    And there should be 0 bursts

  @sad
  Scenario: Empty description shows validation error
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I try to submit without description
    Then I should see a capture validation error

  Scenario: Event with all available tags
    When I select "capture_event" from the menu
    And I select manual capture strategy
    And I capture an event:
      | field       | value                                    |
      | description | Comprehensive tagged event               |
      | tags        | project,achievement,leadership,technical,consulting,research,product,mentoring |
    And I submit the event
    Then there should be 1 event
    And the event should have 8 tags

  Scenario: Event with all available categories
    When I select "capture_event" from the menu
    And I select manual capture strategy
    And I capture an event:
      | field       | value                                    |
      | description | Comprehensive categorized event          |
      | categories  | technical,leadership,product,consulting,research,mentoring,communication,collaboration,problem-solving,project-management,architecture |
    And I submit the event
    Then there should be 1 event
    And the event should have 11 categories

  @happy
  Scenario: Quick submit with Ctrl+S from form
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Quick submit test"
    And I press Ctrl+S
    Then I should see the success message
    And there should be 1 event

  # ============================================================================
  # Date Input Scenarios
  # ============================================================================

  @happy
  Scenario: Capture event with specific date
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Event on specific date"
    And I set event date to "2025-06-15"
    And I submit the event
    Then I should see the success message
    And the event should have date "2025-06-15"

  @happy
  Scenario: Capture event with today keyword
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Event happening today"
    And I set event date to "today"
    And I submit the event
    Then I should see the success message
    And the event should have today's date

  @happy
  Scenario: Capture event with relative date
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Event from last week"
    And I set event date to "-7d"
    And I submit the event
    Then I should see the success message
    And the event should have a date 7 days ago

  @sad
  Scenario: Invalid date format shows validation error
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Event with bad date"
    And I set event date to "not-a-date"
    And I submit the event
    Then I should see a date validation error

  # ============================================================================
  # Fact Editing Scenarios (Review State)
  # ============================================================================

  @happy @enrichment
  Scenario: Edit inferred fact during review
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Reduced API latency by 40% through caching optimization"
    And I submit the event
    And I dismiss the success modal
    When I open the facts editor
    And I edit the first fact
    And I change fact text to "Improved API performance by 40%"
    And I save the fact edit
    And I confirm the review
    Then there should be a fact with text "Improved API performance by 40%"

  @happy @enrichment
  Scenario: Reject inferred fact during review
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Built a new feature with some achievements"
    And I submit the event
    And I dismiss the success modal
    When I open the facts editor
    And I reject all suggested facts
    And I confirm the review
    Then there should be 0 facts

  @happy @enrichment
  Scenario: Add manual fact during review
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Completed a project"
    And I submit the event
    And I dismiss the success modal
    When I open the facts editor
    And I add a new fact "Delivered project 2 weeks ahead of schedule"
    And I confirm the review
    Then there should be a fact with text "Delivered project 2 weeks ahead of schedule"

  # ============================================================================
  # Review Navigation Scenarios
  # ============================================================================

  @happy @enrichment
  Scenario: Navigate through review sections
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Built microservices architecture"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    When I press 'm' to open metadata editor
    Then I should see the metadata modal
    When I close the modal
    And I press 'b' to open bursts editor
    Then I should see the bursts modal
    When I close the modal
    And I press 'f' to open facts editor
    Then I should see the facts modal

  # ============================================================================
  # Skills Selection (Manual Capture)
  # ============================================================================

  @happy
  Scenario: Manual capture with skills selection
    When I select "capture_event" from the menu
    And I select manual capture strategy
    And I capture an event:
      | field       | value                                    |
      | description | Built data pipeline for analytics        |
      | company     | DataCo                                   |
      | skills      | Python,Apache Spark,AWS                  |
    And I submit the event
    Then I should see the success message
    And the event should have skills "Python,Apache Spark,AWS"

  # ============================================================================
  # Description Validation
  # ============================================================================

  @sad
  Scenario: Description too short shows validation error
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Short"
    And I submit the event
    Then I should see a validation error about minimum length

  @sad
  Scenario: Description at minimum length is accepted
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Ten chars!"
    And I submit the event
    Then I should see the success message

  # ============================================================================
  # Review Screen Key Badges
  # ============================================================================

  @happy @enrichment
  Scenario: Review screen shows burst badge when bursts inferred
    Given I have an event "Built API gateway" at company "TechCo"
    Given I have an event "Deployed API gateway" at company "TechCo"
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Scaled API gateway for production load"
    And I set event company to "TechCo"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    And I should see "b" key badge for editing bursts

  @happy @enrichment
  Scenario: Review screen shows fact badge when facts inferred
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Reduced API latency by 40% through optimization"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    And I should see "f" key badge for editing facts

  @happy @enrichment
  Scenario: Review screen shows both badges when bursts and facts inferred
    Given I have an event "Started migration project" at company "TechCo"
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Completed database migration reducing query time by 50%"
    And I set event company to "TechCo"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    And I should see "b" key badge for editing bursts
    And I should see "f" key badge for editing facts
