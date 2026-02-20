@bug @vhs-only
Feature: Bug #55 — Profile Export Shows Hardcoded "John Doe" Placeholder
  As a user
  I want to export my profile
  So that I can share my actual career data

  # GitHub issue: https://github.com/baphled/kariya/issues/55

  Background:
    Given I am on the main menu

  @bug-055
  Scenario: Profile export displays hardcoded placeholder name instead of actual profile data
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I press Ctrl+S
    And I press enter
