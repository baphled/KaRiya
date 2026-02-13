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

  @sad
  Scenario: A user cannot proceed without entering a name
    When I press enter
    Then I should see "Step 1 of 3"

  @smoke
  Scenario: View initial onboarding screen
    Then I should see "Profile Setup"
    And I should see "Step 1 of 3"
    And I should see "Welcome to KaRiya"

  @smoke
  Scenario: Complete Step 1 with name
    When I enter "Test User" as my name
    And I press enter
    Then I should see "Step 2 of 3"

  Scenario: View Step 2 fields
    When I enter "Test User" as my name
    And I press enter
    Then I should see one of:
      | Email    |
      | email    |
    And I should see one of:
      | Location |
      | location |

  Scenario: Complete Step 2 with email and location
    When I enter "Test User" as my name
    And I press enter
    And I enter "test@example.com" as my email
    And I press tab
    And I enter "London" as my location
    And I press enter
    Then I should see "Step 3 of 3"

  Scenario: View Step 3 professional details
    When I enter "Test User" as my name
    And I press enter
    And I enter "test@example.com" as my email
    And I press tab
    And I enter "London" as my location
    And I press enter
    Then I should see one of:
      | Professional |
      | Title        |
      | GitHub       |
      | Portfolio    |

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

  @sad
  Scenario: Invalid email format shows validation error
    When I enter "Test User" as my name
    And I press enter
    And I enter "not-an-email" as my email
    And I press enter
    Then I should still be on step 2
    And I should see a validation error

  @sad
  Scenario: Empty email shows validation error
    When I enter "Test User" as my name
    And I press enter
    And I press enter
    Then I should still be on step 2

  # ============================================================================
  # Navigation Scenarios
  # ============================================================================

  @sad
  Scenario: Escape key is blocked during onboarding
    When I press escape
    Then I should still see the onboarding wizard
     And I should see "Step 1 of 3"

