# Navigation Test Coverage Report

**Date:** 2026-01-03

## Overview

A thorough audit was performed of all navigation-related test code throughout the KaRiya codebase (TUI application). The goal was to determine whether navigation flows—in particular, flows beyond the main menu, such as within or between intents—are comprehensively tested.

## Findings

### 1. Main Menu Navigation
- The main menu's presence and functionality are apparently implicitly tested and at least partially exercised as part of the root model.

### 2. Intent-Level Navigation
- There are _no explicit tests_ for user navigation events ("back", "forward", "cancel", etc) within any intent modules. Searches for relevant keywords and code patterns turned up empty in all major intent test files.
- Example intent tests checked: `capture_event_integration_test.go`, `burst_management_test.go`, `browse_timeline_test.go`, `export_artifact_test.go`, `configure_system_test.go`, `bulk_operations_test.go`, `import_wizard_test.go`, `fact_management_test.go`, and `metadata_editor_test.go`.

### 3. Cross-Intent/Screen Navigation
- No explicit cross-intent navigation scenarios (e.g. activating one intent from another, or exercising the router's back and context stack) are tested in any integration or router test files.
- Components and navigation helper test files (e.g. `navigation_menu_test.go`, `key_handler_test.go`) do not cover application-level navigation, only key handling or control rendering logic.

### 4. Field/Row Navigation
- Narrow navigation behaviors such as field focus cycling within a form (via j/k or Tab/Shift+Tab) _are_ tested (`form_jk_navigation_test.go`), but these operate at a field or component level, **not at the application or intent navigation level**.

## Gaps and Recommendations

- **Lack of coverage** for navigation within intents: There are no tests that simulate or assert a user's ability to go back, move forward, or cancel out of multi-step flows in any intent.
- **No tests** cover switching between intents or returning to preserved context after navigating back—the core advantage of the intent-driven architecture.
- **Indirect/no coverage** for edge-case navigation (e.g. exiting during confirmation steps, preserving state after back navigation).

### Next Steps / Suggestions

1. **Add Integration Tests:** Create tests at the app/model level that:
    - Simulate selecting menu options, activating intents
    - Step through various screens within an intent, using back/forward/cancel navigation
    - Exercise the router's stack/context preservation features
    - Assert state transitions, context preservation, and correct rendering/views after navigation actions

2. **Intent-Level Navigation Tests:** For each intent module, add tests that:
    - Start at the initial state, progress through intermediate states, use navigation actions to move between states, and confirm end state matches expectations
    - Test edge cases (e.g. pressing back on first or last step, cancelling, or nested modal flows)

3. **Cross-Intent/Advanced:** Add test cases for routing between intents (e.g. invoking one intent from another, or restoring prior context after back navigation)

## Conclusion

Navigation is only reliably tested for the main menu and form field focus cycling. All broader app and intent navigation patterns are untested.

**Improvement is strongly recommended by adding explicit integration and intent-level navigation tests to ensure the robustness of the user navigation experience.**

