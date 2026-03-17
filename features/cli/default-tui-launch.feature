Feature: Default TUI Launch
  As a user
  I want to launch KaRiya from the CLI
  So that I get the interactive TUI when no arguments are provided

  Scenario: Running kariya with no arguments launches TUI
    When I run "kariya" with no arguments
    Then the TUI should launch and display the main menu

  Scenario: Running kariya with --help displays help (not TUI)
    When I run "kariya --help"
    Then the CLI should display help text
    And the TUI should not launch

  Scenario: Running kariya with subcommand runs subcommand
    When I run "kariya version"
    Then the CLI should display the version
    And the TUI should not launch
