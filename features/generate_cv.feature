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

  @happy @wip
  Scenario: Navigate wizard with tab
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I press tab
    Then I should be on audience field
    When I press tab
    Then I should be on audience field
  @happy @wip
  Scenario: Cancel wizard with escape
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I press escape
    Then I should be on the main menu

  # ============================================================================
  # CV Configuration Wizard - Step 2 (TECH)
  # ============================================================================

  @happy @wip
  Scenario: View step 2 technology options
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    Then I should see "Technology Focus"
    And I should see "Language Agnostic"
    And I should see "Generalist"
    And I should see "Specialist"

  @happy @wip
  Scenario: Select language agnostic focus
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    And I select "Language Agnostic" technology focus
    Then I should skip technology selection
    And I should see "Skills Presentation"
  @happy @wip
  Scenario: Select generalist focus shows multi-select
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    And I select "Generalist" technology focus
    Then I should see "Skills Presentation"
  @happy @wip
  Scenario: Select specialist focus shows single-select
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    And I select "Specialist" technology focus
    Then I should see "Skills Presentation"

  @happy @wip
  Scenario: View focus area options
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    And I select technology focus
    Then I should see "Skills Presentation"
  @happy @wip
  Scenario: Navigate back to step 1 with escape
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    And I press escape
    Then I should be back on step 1

  # ============================================================================
  # CV Configuration Wizard - Step 3 (FORMAT)
  # ============================================================================

  @happy @wip
  Scenario: View step 3 format options
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    And I complete step 2
    Then I should see "Skills Presentation"
    And I should see "Skills Limit"
    And I should see "CV Length"

  @happy @wip
  Scenario: Select skills format
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    And I complete step 2
    Then I should see "Grouped"
    And I should see "Flat"

  @happy @wip
  Scenario: Set skills limit
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    And I complete step 2
    And I tab to skills limit
    And I enter skills limit "10"
    Then the skills limit should be 10

  @happy @wip
  Scenario: Select CV length
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete step 1
    And I complete step 2
    And I tab to CV length
    Then I should see "CV Length"
    And I should see "Target length for the CV"

  @happy @wip
  Scenario: Skip wizard with Ctrl+S
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I press Ctrl+S to skip
    Then I should see the progress modal
    And I should see "Generating"

  # ============================================================================
  # Technology Extraction
  # ============================================================================

  @happy @wip
  Scenario: View technology extraction progress
    Given I have a complete profile with skills
    When I select "generate_cv" from the menu
    Then I should see the progress modal
    And I should see "Extracting Technologies"

  @happy @wip
  Scenario: Cancel technology extraction
    Given I have a complete profile with skills
    When I select "generate_cv" from the menu
    And I see the extracting progress
    And I press escape
    Then I should be on the main menu

  # ============================================================================
  # CV Generation Progress
  # ============================================================================

  @happy @wip
  Scenario: View CV generation progress
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete the wizard
    Then I should see the progress modal
    And I should see "Generating CV"

  @happy @wip
  Scenario: Cancel CV generation
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I complete the wizard
    And I see the generating progress
    And I press escape
    Then I should be on the main menu

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
  Scenario: Review screen shows selected audience
    Given I have generated a CV
    When I am on the CV review screen
    Then I should see "hiring_manager"

  @happy
  Scenario: Preview screen shows audience-specific highlights
    Given I have generated a CV
    When I navigate to the CV preview screen
    Then I should see "Key Highlights"

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
  Scenario: Export Text to file
    Given I have generated a CV
    When I open the export options modal
    And I select format "Text"
    And I select location "File"
    And I confirm export
    Then I should see export progress
    And the export should complete
    And CV is exported as a text file

  @happy
  Scenario: Export Markdown to file
    Given I have generated a CV
    When I open the export options modal
    And I select format "Markdown"
    And I select location "File"
    And I confirm export
    Then I should see export progress
    And the export should complete
    And CV is exported as a markdown file

  @happy
  Scenario: Export YAML to file
    Given I have generated a CV
    When I open the export options modal
    And I select format "YAML"
    And I select location "File"
    And I confirm export
    Then I should see export progress
    And the export should complete
    And CV is exported as a YAML file

  @happy
  Scenario: Export as Text to clipboard
    Given I have generated a CV
    When I open the export options modal
    And I select format "Text"
    And I select location "Clipboard"
    And I confirm export
    Then I should see export progress
    And the export should complete
    And the clipboard should contain the CV as text

  @happy
  Scenario: Export as Markdown to clipboard
    Given I have generated a CV
    When I open the export options modal
    And I select format "Markdown"
    And I select location "Clipboard"
    And I confirm export
    Then I should see export progress
    And the export should complete
    And the clipboard should contain the CV as markdown

  @happy
  Scenario: Export as YAML to clipboard
    Given I have generated a CV
    When I open the export options modal
    And I select format "YAML"
    And I select location "Clipboard"
    And I confirm export
    Then I should see export progress
    And the export should complete
    And the clipboard should contain the CV as YAML

  # ============================================================================
  # Error Handling
  # ============================================================================

  @sad @wip
  Scenario: Handle generation error
    Given CV generation will fail
    When I select "generate_cv" from the menu
    And I complete the wizard
    Then I should see an error modal
    And I should see error details

  @sad @wip
  Scenario: Handle export error
    Given I have generated a CV
    And export will fail
    When I start an export
    Then I should see an error modal
    And I should be able to retry

  # ============================================================================
  # Navigation and Exit
  # ============================================================================

