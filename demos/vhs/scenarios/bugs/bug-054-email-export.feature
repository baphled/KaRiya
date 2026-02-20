@bug @vhs-only
Feature: Bug #54 — Email Export Fails or Shows Missing Implementation
  As a user
  I want to export my CV via email
  So that I can send it directly to recruiters

  # GitHub issue: https://github.com/baphled/kariya/issues/54

  Background:
    Given I am on the main menu

  @bug-054
  Scenario: Email export shows failure or missing implementation
    Given I have a complete profile with events and facts
    When I select "generate_cv" from the menu
    And I press Ctrl+S
    And I navigate down
    And I navigate down
    And I press enter
