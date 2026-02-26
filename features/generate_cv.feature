@cv
Feature: Generate CV
  As a professional
  I want to generate a CV from my career data
  So that I can present my experience effectively to different audiences

  Background:
    Given I am on the main menu

  # ============================================================================
  # Empty State / Prerequisites
  # ============================================================================

  @sad @smoke
  Scenario: Cannot generate CV without profile
    Given I have no profile configured
    When I select "generate_cv" from the menu
    Then I should see "Career Events"
    And I should see an error or warning

   @sad @smoke
   Scenario: Cannot generate CV without events
     Given I have a profile configured
     But I have no data
     When I select "generate_cv" from the menu
     Then I should see "events"
     And I should see an error or warning

  # ============================================================================
  # CV Configuration Wizard - Step 1 (WHO)
  # ============================================================================

  @happy @smoke
  Scenario: Open CV wizard modal
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    Then I should see the CV wizard modal
    And I should see "Profile"
    And I should see "Audience"

  @happy
  Scenario: Select profile in wizard
    Given I have multiple profiles
    When I select "generate_cv" from the menu
    And I navigate down in the profile selector
    And I confirm selection
    Then I should move to the next field

  @happy
  Scenario: Select audience in wizard
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I tab to audience field
    Then I should see "Hiring Manager"
    And I should see "Recruiter"
    And I should see "Peer"

   # ============================================================================
   # Technology Extraction
   # ============================================================================

  # ============================================================================
  # CV Generation Progress
  # ============================================================================

  # ============================================================================
  # CV Review Screen
  # ============================================================================

  @happy
  Scenario: View CV review after generation
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete the wizard
    And the generation completes
    Then I should see the CV review screen
    And I should see CV metadata
    And I should see statistics

  @happy
  Scenario: Review screen shows section summary
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete the wizard
    And the generation completes
    And I press enter to preview
    Then I should see section names

  @happy
  Scenario: Navigate to preview from review
    Given I have generated a CV
    When I am on the CV review screen
    And I press enter to preview
    Then I should see the CV preview screen

  @happy
  Scenario: Export directly from review
    Given I have generated a CV
    When I am on the CV review screen
    And I press "x" to export
    Then I should see the export options modal

  @happy
  Scenario: Edit CV from review
    Given I have generated a CV
    When I am on the CV review screen
    And I press "e" to edit
    Then I should see the CV wizard modal

  @happy
  Scenario: Review screen shows highlights
    Given I have generated a CV
    When I am on the CV review screen
    Then I should see highlights
  # ============================================================================
  # CV Preview Screen
  # ============================================================================

  @happy
  Scenario: View full CV preview
    Given I have generated a CV
    When I navigate to the CV preview screen
    Then I should see personal details
    And I should see all CV sections
     And I should see bullet points

   @happy
   Scenario: Confirm CV from preview
    Given I have generated a CV
    When I navigate to the CV preview screen
    And I press enter to confirm
    Then the CV generation should complete
    And I should be on the main menu

  @happy
  Scenario: Export from preview
    Given I have generated a CV
    When I navigate to the CV preview screen
    And I press "x" to export
    Then I should see the export options modal

  @happy
  Scenario: Edit from preview
    Given I have generated a CV
    When I navigate to the CV preview screen
    And I press "e" to edit
    Then I should see the CV wizard modal

  @happy
  Scenario: Preview screen shows key highlights
    Given I have generated a CV
    When I navigate to the CV preview screen
    Then I should see key highlights
  # ============================================================================
  # Export Options Modal
  # ============================================================================

  @happy
  Scenario: View export format options
    Given I have generated a CV
    When I open the export options modal
    Then I should see "Text"
    And I should see "Markdown"
    And I should see "YAML"

  @happy
  Scenario: View export location options
    Given I have generated a CV
    When I open the export options modal
    And I tab to location
    Then I should see "File"
    And I should see "Clipboard"

  @happy
  Scenario: Export to file as text
    Given I have generated a CV
    When I open the export options modal
    And I select format "Text"
    And I select location "File"
    And I confirm export
    Then I should see export progress
    And the export should complete

  @happy
  Scenario: Export to clipboard as markdown
    Given I have generated a CV
    When I open the export options modal
    And I select format "Markdown"
    And I select location "Clipboard"
    And I confirm export
    Then I should see export progress
    And the export should complete

  # ============================================================================
  # Export Progress
  # ============================================================================

  @happy
  Scenario: Export success shows confirmation
    Given I have generated a CV
    When I complete an export
    Then I should see success message
    And I should see export location

  # ============================================================================
  # Error Handling
  # ============================================================================

  # ============================================================================
  # Navigation and Exit
  # ============================================================================

