@onboarding
Feature: User Onboarding
  As a new user
  I want to complete the onboarding wizard
  So that I can set up my profile and start using KaRiya

  Background:
    Given I start the onboarding wizard

  @happy @smoke
  Scenario: A user successfully completes onboarding with required fields
    When I enter "Jane Developer" as my name
    And I press enter
    And I enter "jane@example.com" as my email
    And I press tab
    And I enter "London, UK" as my location
    And I press enter
    And I skip the optional fields
    Then the onboarding wizard should be complete
    And my profile should have name "Jane Developer"
    And my profile should have email "jane@example.com"
    And my profile should have location "London, UK"

  # ============================================================================
  # Step 3 Completion (Optional Fields)
  # ============================================================================

  @happy
  Scenario: Complete Step 3 with professional title
    When I enter "Test User" as my name
    And I press enter
    And I enter "test@example.com" as my email
    And I press tab
    And I enter "London" as my location
    And I press enter
    And I enter "Senior Software Engineer" as my title
    And I skip the optional fields
    Then the onboarding wizard should be complete
    And my profile should have title "Senior Software Engineer"

  @happy
  Scenario: Complete onboarding with all optional fields
    When I enter "Jane Developer" as my name
    And I press enter
    And I enter "jane@example.com" as my email
    And I press tab
    And I enter "London, UK" as my location
    And I press enter
    And I enter "Principal Engineer" as my title
    And I press tab
    And I enter "janedev" as my GitHub username
    And I press tab
    And I enter "https://janedev.com" as my portfolio
    And I press enter
    Then the onboarding wizard should be complete
    And my profile should have title "Principal Engineer"
    And my profile should have GitHub username "janedev"
    And my profile should have portfolio "https://janedev.com"

  @happy
  Scenario: Skip all optional fields in Step 3
    When I enter "Test User" as my name
    And I press enter
    And I enter "test@example.com" as my email
    And I press tab
    And I enter "London" as my location
    And I press enter
    And I skip the optional fields
    Then the onboarding wizard should be complete
    And my profile should not have a title

  # ============================================================================
  # Validation Scenarios
  # ============================================================================

  # ============================================================================
  # Navigation Scenarios
  # ============================================================================

