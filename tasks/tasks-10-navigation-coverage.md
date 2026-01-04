# Task List: Navigation Test Coverage Enhancement for KaRiya TUI

**PRD Reference**: `docs/features/10-navigation-coverage.md`

**Purpose**: Ensure robust, explicit, and comprehensive test coverage for all navigation flows within the KaRiya TUI application.

**Status**: 🚀 READY FOR IMPLEMENTATION

**Version**: 1.0 - Initial Task Generation (2026-01-03)

**Effort Estimate**: 15 hours

**Priority**: HIGH

---

## Relevant Files

### Files to be Created
- `internal/cli/app/navigation_integration_test.go` - Integration tests for app-level navigation, including cross-intent scenarios and edge cases.
- `internal/cli/intents/navigation_helpers_test.go` - Tests for intent-level navigation helpers.
- `docs/NAVIGATION_TEST_STRATEGY.md` - Documentation outlining the navigation test strategy and intended coverage targets.

### Files to be Modified
- `internal/cli/intents/router_result_handler_test.go` - Extend tests to cover navigation edge cases, such as restoring context after back navigation.
- `internal/cli/intents/[intent_name]_test.go` (for each multi-step intent e.g., `capture_event_test.go`, `bulk_operations_test.go`) - Add intent-specific navigation tests.
- `internal/cli/components/navigation_menu_test.go` - Verify navigation actions triggered by menu selection.

### Reference Files
- `docs/NAVIGATION_TEST_COVERAGE_REPORT.md` - Source audit for navigation gaps.

---

## Tasks

- [x] 1.0 App-Level Navigation Tests (Phase 1)
  - [x] 1.1 Create `navigation_integration_test.go` to simulate main menu navigation and cross-intent flow.
  - [x] 1.2 Write tests to cover back, forward, and cancel navigation at the app level.
  - [x] 1.3 Test edge cases like immediate back from first step or resuming context after back navigation.
  - [x] 1.4 Assert proper context restoration after using the router's stack.
  - [x] 1.5 Run tests and verify no regressions.

- [x] 2.0 Intent-Level Navigation Tests (Phase 2)
  - [x] 2.1 For each multi-step intent (`capture_event`, `bulk_operations`, etc.), write tests for forward, backward, and cancel flows.
  - [x] 2.2 Test edge cases like cancel after submission and nested modal flows.
  - [x] 2.3 Validate state preservation and correct transitions.
  - [x] 2.4 Run tests and verify no regressions.

- [x] 3.0 Navigation Component Tests (Phase 3)
  - [x] 3.1 Extend `navigation_menu_test.go` to simulate menu-based navigation actions.
  - [x] 3.2 Add tests to verify correct intent activation from menu selections.
  - [x] 3.3 Run tests and verify no regressions.

- [ ] 4.0 Documentation and Final Verification (Phase 4)
  - [ ] 4.1 Create `NAVIGATION_TEST_STRATEGY.md` documenting the strategy, coverage expectations, and edge cases tested.
  - [ ] 4.2 Perform a full test suite run and ensure 100% reliability.
  - [ ] 4.3 Confirm 90%+ code coverage on navigation paths.
  - [ ] 4.4 Final commit and handoff.

---

## High-Level Overview

### What This Task Accomplishes

This task list ensures that all navigation flows—both within and across intents—are explicitly tested, reducing the risk of undetected regressions and improving the robustness of the user experience.

**Phase 1: App-Level Navigation Tests (5 hours)**
- Simulate menu and cross-intent navigation.
- Test edge cases like context restoration after back navigation.

**Phase 2: Intent-Level Tests (5 hours)**
- Test forward/backward/cancel flows for all intents.
- Cover edge cases like nested workflows and cancellation after submission.

**Phase 3: Navigation Menu Tests (3 hours)**
- Extend existing menu tests to verify proper integration with intent activation.

**Phase 4: Documentation & Verification (2 hours)**
- Document navigation testing strategy.
- Achieve 90%+ code coverage and verify full test suite reliability.

### Success Criteria

✅ All app-level navigation paths tested.
✅ Multi-step intents have comprehensive navigation tests.
✅ Navigation menu integration verified.
✅ Strategy documented and test reliability confirmed.
✅ Code coverage exceeds 90%.
✅ Zero regressions or crashes.

---

**Document Version**: 1.0

