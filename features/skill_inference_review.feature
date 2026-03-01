Feature: Review inferred skills during event capture
  As a professional
  I want to review and manage inferred skills during event capture
  So that I can ensure accurate skill tracking for my career profile

  Background:
    Given the database is empty
    And I am on the main menu

  @happy @smoke @wip
  Scenario: User can view suggested skills with confidence scores
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Built microservices in Go with PostgreSQL database"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    And I should see "s" key badge for skills
    When I open the skill review modal
    Then I should see suggested skills with confidence scores
    And I should see "Go" skill with confidence

  @happy @enrichment @wip
  Scenario: User can accept suggested skills
    Given I have an event "Implemented REST API with Python and Flask" at company "TechCorp"
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Added authentication layer using JWT tokens"
    And I set event company to "TechCorp"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    When I accept all inferred skills
    And I confirm the review
    Then I should be on the main menu
    And there should be skills including "Python"

  @happy @enrichment @wip
  Scenario: User can reject suggested skills
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Worked with Docker and Kubernetes for container orchestration"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    When I reject all suggested skills
    And I confirm the review
    Then I should be on the main menu
    And the event should have no skills

  @happy @enrichment @wip
  Scenario: Accepted skills are persisted in result
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Deployed applications on AWS using CloudFormation and Terraform"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    When I accept all inferred skills
    And I confirm the review
    Then I should be on the main menu
    And there should be skills including "AWS"
    And there should be skills including "Terraform"

  @happy @enrichment @wip
  Scenario: Skill review modal is accessible via 's' key
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Optimized database queries using PostgreSQL and Redis caching"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    When I press the "s" key
    Then I should see the skill review modal
    And I should see suggested skills with confidence scores

  @happy @enrichment @wip
  Scenario: User can selectively accept and reject skills
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Developed frontend with React and TypeScript, backend with Node.js"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    When I accept the "React" skill
    And I reject the "TypeScript" skill
    And I confirm the review
    Then I should be on the main menu
    And there should be skills including "React"
    And the event should not have "TypeScript" skill

  @happy @enrichment @wip
  Scenario: Navigate between review sections including skills
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Architected microservices using Go and gRPC protocols"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    When I press the "e" key
    Then I should see the metadata editor
    When I press the "s" key
    Then I should see the skill review modal
    When I press the "b" key
    Then I should see the burst editor

  @sad @wip
  Scenario: Rejecting all skills results in no skills attached
    When I select "capture_event" from the menu
    And I select quick capture strategy
    And I enter event description "Managed CI/CD pipelines with Jenkins and GitLab"
    And I submit the event
    And I dismiss the success modal
    Then I should be on the enrichment review screen
    When I reject all suggested skills
    And I confirm the review
    Then I should be on the main menu
    And the event should have no skills
