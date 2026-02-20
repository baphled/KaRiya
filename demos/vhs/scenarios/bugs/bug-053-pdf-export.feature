@bug @vhs-only
Feature: Bug #53 — PDF Export Shows Placeholder Text
  As a user
  I want to export my CV as a PDF
  So that I can share a professionally formatted document

  # GitHub issue: https://github.com/baphled/kariya/issues/53

  Background:
    Given I am on the main menu

  @bug-053
  Scenario: PDF export shows placeholder text instead of actual content
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I press Ctrl+S
    And I navigate down
    And I press enter
