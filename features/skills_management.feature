@skills
Feature: Manage Skills
  As a professional
  I want to manage my skills
  So that I can track my expertise and proficiency levels

  Background:
    Given I am on the main menu

  # ============================================================================
  # Empty State
  # ============================================================================

   @happy @smoke
   Scenario: View empty skills list
     Given I have no data
     When I select "manage_skills" from the menu
     Then I should see "No skills"
     And I should be able to go back to the menu

  # ============================================================================
  # Skills List Display
  # ============================================================================

  @happy @smoke
  Scenario: View skills list with skills
    Given I have 5 skills in my profile
    When I select "manage_skills" from the menu
    Then I should see a list of skills
    And I should see "5"

   # ============================================================================
   # Skill Detail View
   # ============================================================================

  @happy
  Scenario: View skill details
    Given I have a skill "Go" with category "backend"
    When I select "manage_skills" from the menu
    And I press enter to view details
    Then I should see "Go"
    And I should see "backend"



  # ============================================================================
  # Add Skill
  # ============================================================================

   @happy
   Scenario: Open add skill modal
     Given I have no data
     When I select "manage_skills" from the menu
     And I press "a" to add skill
     Then I should see the add skill form



   @happy
   Scenario: Add a new skill
     Given I have no data
     When I select "manage_skills" from the menu
     And I press "a" to add skill
     And I enter skill name "Python"
     And I select category "backend"
     And I submit the skill form
     Then I should still be on the skills list
     And there should be 1 skill
     And the skill should have name "Python"

  # ============================================================================
  # Edit Skill
  # ============================================================================

  @happy
  Scenario: Open edit skill modal
    Given I have a skill "JavaScript" with category "backend"
    When I select "manage_skills" from the menu
    And I press "e" to edit
    Then I should see the edit skill form



  @happy
  Scenario: Edit skill and save
    Given I have a skill "JavaScript" with category "backend"
    When I select "manage_skills" from the menu
    And I press "e" to edit
    And I clear the skill name field
    And I enter skill name "TypeScript"
    And I submit the skill form
    Then I should still be on the skills list
    And the skill should have name "TypeScript"

  # ============================================================================
  # Delete Skill
  # ============================================================================

  @happy
  Scenario: Open delete skill confirmation
    Given I have a skill "Ruby" with category "backend"
    When I select "manage_skills" from the menu
    And I press "d" to delete
    Then I should see the delete confirmation



  @happy
  Scenario: Confirm delete skill
    Given I have a skill "Ruby" with category "backend"
    When I select "manage_skills" from the menu
    And I press "d" to delete
    And I confirm the deletion
    Then I should still be on the skills list
    And there should be 0 skills

  # ============================================================================
  # Skill Inference
  # ============================================================================

  @happy
  Scenario: Review skill suggestions
    Given I have an event "Built microservices in Go with PostgreSQL database"
    When I select "manage_skills" from the menu
    And I press "i" to infer skills
    And the inference completes
    Then I should see the skill suggestions modal
    And I should see "Go"

  @happy
  Scenario: Accept skill suggestion
    Given I have an event "Built microservices in Go"
    When I select "manage_skills" from the menu
    And I press "i" to infer skills
    And the inference completes
    And I accept the first suggestion
    Then I should see "Success"
    And there should be 1 skill

  @happy
  Scenario: Reject skill suggestion
    Given I have an event "Built microservices in Go"
    When I select "manage_skills" from the menu
    And I press "i" to infer skills
    And the inference completes
    And I reject the first suggestion
    Then I should see "No skills found. Press 'a' to add your first skill."
    And the suggestion should be marked as rejected

  # ============================================================================
  # Skill Inference with Save (PR #172 - event-skill linkage)
  # ============================================================================

  @happy @inference
  Scenario: Infer skill and save links it to the event
    Given I have an event "Built REST API in Go with PostgreSQL"
    And no skills are linked to the event
    When I select "manage_skills" from the menu
    And I infer and accept skill "Go" for the event
    Then "Go" should be linked to the event
    And there should be 1 skill

  @happy @inference
  Scenario: Globally existing skill not linked to event appears as suggestion
    Given I have a skill "Go" with category "backend"
    And I have an event "Built microservices in Go with gRPC"
    And "Go" is not linked to the event
    When I trigger inference for the event
    Then "Go" should be suggested as a new skill
    And "Go" should not be in the existing skills list

  @happy @inference
  Scenario: Skill already linked to event appears as existing
    Given I have a skill "Go" with category "backend"
    And I have an event "Built API in Go" that uses skill "Go"
    When I trigger inference for the event
    Then "Go" should be in the existing skills list
    And "Go" should not be suggested as a new skill

  @happy @inference
  Scenario: Accept inferred skill persists link in event_skills
    Given I have an event "Deployed apps on AWS using Terraform"
    When I select "manage_skills" from the menu
    And I infer and accept skill "AWS" for the event
    And I infer and accept skill "Terraform" for the event
    Then "AWS" should be linked to the event
    And "Terraform" should be linked to the event
    And there should be 2 skills

  @happy @inference
  Scenario: Re-running inference after save shows skill as existing
    Given I have an event "Built REST API in Go"
    When I select "manage_skills" from the menu
    And I infer and accept skill "Go" for the event
    And I trigger inference for the event
    Then "Go" should be in the existing skills list


  # ============================================================================
  # View Events Using Skill
  # ============================================================================

  @happy
  Scenario: View events that use a skill
    Given I have a skill "Go" with category "backend"
    And I have an event "Built API in Go" that uses skill "Go"
    And I have an event "Wrote CLI tool in Go" that uses skill "Go"
    When I select "manage_skills" from the menu
    And I select skill "Go"
    And I press "s" to view events
    Then I should see the skill events modal
    And I should see 2 events
    And I should see "Built API in Go"
    And I should see "Wrote CLI tool in Go"

  @happy
  Scenario: View event detail from skill events modal
    Given I have a skill "Python" with category "backend"
    And I have an event "Built data pipeline" that uses skill "Python"
    When I select "manage_skills" from the menu
    And I select skill "Python"
    And I press "s" to view events
    And I press enter to view event details
    Then I should see "Built data pipeline"
    And I should see the full event description



  # ============================================================================
  # Proficiency Level
  # ============================================================================

  @happy
  Scenario: Add skill with proficiency level
    Given I have no data
    When I select "manage_skills" from the menu
    And I press "a" to add skill
    And I enter "Kubernetes" as skill name
    And I select category "DevOps"
    And I select level "Expert"
    And I submit the skill form
    Then there should be 1 skill
    And the skill should have level "Expert"

  @happy
  Scenario: Edit skill proficiency level
    Given I have a skill "Docker" with category "devops" and level "Intermediate"
    When I select "manage_skills" from the menu
    And I select skill "Docker"
    And I press "e" to edit
    And I change level to "Expert"
    And I submit the skill form
    Then the skill should have level "Expert"

  # ============================================================================
  # Years of Experience
  # ============================================================================

  @happy
  Scenario: Add skill with years of experience
    Given I have no data
    When I select "manage_skills" from the menu
    And I press "a" to add skill
    And I enter "Java" as skill name
    And I select category "Languages"
    And I enter years of experience "8"
    And I submit the skill form
    Then there should be 1 skill
    And the skill should have years "8"

  @happy
  Scenario: Edit skill years of experience
    Given I have a skill "Python" with category "backend" and years "3"
    When I select "manage_skills" from the menu
    And I select skill "Python"
    And I press "e" to edit
    And I change years to "5"
    And I submit the skill form
    Then the skill should have years "5"

