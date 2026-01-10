# Task 38: Navigation Pattern Standardization

## Overview
- **Goal**: Standardize keyboard shortcuts, navigation behavior, and selection preservation across all 10 intents using library-based solutions
- **Time Estimate**: 8-12 days (Phases 1-5 core, Phases 6-8 optional)
- **Prerequisites**: All tests passing, understanding of intent architecture

## Session Contract Acknowledgment
- [x] Ran `make session-start` and it passed
- [x] Acknowledge and commit to following all workflow rules
- [x] Token count: Within limits

## Pre-Task Checklist (MUST COMPLETE BEFORE STARTING)
- [x] `make check-compliance` passes
- [x] Reviewed existing patterns in:
  - `internal/cli/navigation/list_navigator.go`
  - `internal/cli/intents/contract.go` (BaseIntent)
  - `docs/KEYBOARD_REFERENCE.md`
  - `docs/TUI_STANDARDS.md`
- [x] Confirmed this is ONE atomic task (not multiple changes)
- [x] Identified which test files will be created/modified

---

## Problem Statement

### Current Issues
1. **`q` key inconsistency**: Currently cancels intent (should quit application per docs)
2. **`m` key vs `h` key**: Code uses `m` for main menu, docs specify `h` for home
3. **No `?` help modal**: Docs say `?` should show context-sensitive help
4. **Selection not preserved**: When navigating back, list selection resets to 0
5. **ListNavigationHandler underutilized**: Only 3/10 intents use centralized navigation
6. **Duplicated key handling**: Each intent has copy-pasted key handling code

### Documentation Standards (Source of Truth)
From `docs/KEYBOARD_REFERENCE.md`:
| Key | Action | Expected Behavior |
|-----|--------|-------------------|
| `q` | Quit application | Exit app from any screen |
| `Ctrl+C` | Quit application | Standard interrupt |
| `?` | Show help | Context-sensitive help modal |
| `Esc` | Go back / Cancel | Primary back navigation key |
| `h` | Go home | Return to main menu (global context) |

---

## Decisions

| Item | Decision |
|------|----------|
| **Libraries** | `qmuntal/stateless` + `kevm/bubbleo` (add in Phase 6) |
| **Priority** | Keyboard first (Phases 1-5), state machines later (Phases 6-8) |
| **Key bindings** | Follow docs: `q`=quit, `?`=help modal, `esc`=back |
| **DOT graphs** | Skip - unnecessary overhead |
| **Breaking changes** | Acceptable - `q` and `m` behavior changed |
| **`m` key** | **REMOVED** - no longer mapped (use `esc` for back navigation) |
| **`h` key** | **REMOVED** - conflicted with vim navigation, use `esc` instead |

---

## Files Modified

### Phase 1: KeyMaps Foundation - COMPLETED
- [x] Created: `internal/cli/navigation/keymaps.go`
- [x] Created: `internal/cli/navigation/keymaps_test.go`
- [x] Commit: `496fede feat(cli): add standardized KeyMaps using bubbles/key`

### Phase 2: Help Modal - COMPLETED
- [x] Created: `internal/cli/components/help_modal.go`
- [x] Created: `internal/cli/components/help_modal_test.go`
- [x] Commit: `6d6c3d8 feat(components): add HelpModal using bubbles/help`

### Phase 2b: Remove Conflicting Keys - COMPLETED
- [x] Modified: Multiple intents to remove `h` key bindings
- [x] Commit: `92491b7 fix(cli): remove conflicting 'h' key bindings per KEYBOARD_REFERENCE.md`

### Phase 3: Intent Key Handler Standardization - COMPLETED
- [x] `internal/cli/intents/browse_timeline_intent.go` (2 states)
- [x] `internal/cli/intents/fact_management_intent.go` (3 states)
- [x] `internal/cli/intents/metadata_editor_intent.go` (3 states)
- [x] `internal/cli/intents/bulk_operations_intent.go` (4 states)
- [x] `internal/cli/intents/capture_event_intent.go` (4 states)
- [x] `internal/cli/intents/import_wizard_intent.go` (5 states)
- [x] `internal/cli/intents/burst_management_intent.go` (7 states)
- [x] `internal/cli/intents/configure_system.go` (7 states)
- [x] `internal/cli/intents/export_artifact.go` (9 states)
- [x] `internal/cli/intents/generate_cv_intent.go` (13 states)
- [x] Commit: `88eaef9 feat(intents): standardize BrowseTimeline key handling with HandleGlobalKeys`
- [x] Commit: `87d5795 feat(intents): standardize all intents with HandleGlobalKeys`

### Tests Updated - COMPLETED
- [x] All `q` key tests now expect `tea.Quit` instead of `Cancelled` status
- [x] All `m` key tests removed (key no longer mapped)
- [x] E2E tests updated for new navigation behavior
- [x] App unit tests updated

---

## Implementation Checklist

### Phase 1: KeyMaps Foundation - COMPLETED
**Time**: 1-2 days | **Risk**: Low | **Status**: COMPLETED

- [x] Created KeyMaps using bubbles/key library
- [x] GlobalKeyMap with Quit, Help, Back bindings
- [x] ListKeyMap with navigation bindings
- [x] FormKeyMap with form-specific bindings
- [x] All tests passing

### Phase 2: Help Modal - COMPLETED
**Time**: 1 day | **Risk**: Low | **Status**: COMPLETED

- [x] Created HelpModal using bubbles/help
- [x] Show/Hide functionality
- [x] Key handling (esc/? to close)
- [x] All tests passing

### Phase 3: Intent Key Handler Standardization - COMPLETED
**Time**: 2-3 days | **Risk**: Medium | **Status**: COMPLETED

- [x] Created `HandleGlobalKeys()` helper function in `view_helpers.go`
- [x] Returns `GlobalKeyResult` enum: `KeyNotHandled`, `KeyQuit`, `KeyHelp`, `KeyBack`
- [x] All 10 intents updated to use `HandleGlobalKeys()` pattern
- [x] New behavior:
  - `q` / `ctrl+c` → Returns `tea.Quit` (quits application)
  - `?` → Help modal placeholder (TODO for BaseIntent integration)
  - `esc` → Go back to previous state (or cancel at root)
- [x] Removed all `case "m":` handlers (m key no longer mapped)
- [x] All 2,000+ tests passing

### Phase 4: Selection Preservation - PENDING
**Time**: 1-2 days | **Risk**: Medium | **Status**: NOT STARTED

- [ ] Modify: `internal/cli/intents/contract.go` (add selection state)
- [ ] Modify: `internal/cli/intents/result.go` (add selection metadata)
- [ ] Modify: All list-based intents

### Phase 5: ListNavigationHandler Adoption - IN PROGRESS
**Time**: 1-2 days | **Risk**: Low | **Status**: IN PROGRESS

Already using ListNavigationHandler:
- [x] `browse_timeline_intent.go`
- [x] `burst_management_intent.go`
- [x] `fact_management_intent.go`
- [x] `bulk_operations_intent.go` (migrated 2026-01-10)

Need to add:
- [ ] `configure_system.go`
- [ ] `export_artifact.go`

Note: `metadata_editor_intent.go` does not need migration - it's a form-based editor, not a list navigator.

### Phase 6: Help Modal BaseIntent Integration - PENDING
**Time**: 1 day | **Risk**: Low | **Status**: NOT STARTED

- [ ] Add help modal state to BaseIntent
- [ ] Integrate help modal toggle with `?` key
- [ ] Remove TODO placeholders from all intents

### Phase 7-8: Optional Future Work - PENDING
- [ ] Add stateless + bubbleo dependencies
- [ ] State machine formalization
- [ ] Navigation stack with bubbleo

---

## Git Log (Commits Made)

```
988375a feat(intents): migrate BulkOperations to ListNavigationHandler
87d5795 feat(intents): standardize all intents with HandleGlobalKeys
88eaef9 feat(intents): standardize BrowseTimeline key handling with HandleGlobalKeys
92491b7 fix(cli): remove conflicting 'h' key bindings per KEYBOARD_REFERENCE.md
6d6c3d8 feat(components): add HelpModal using bubbles/help
cb9c11a docs(docs): add Task 38 navigation standardization plan
496fede feat(cli): add standardized KeyMaps using bubbles/key
```

---

## Acceptance Criteria

### Phases 1-3 (COMPLETED)
- [x] `q` quits application from all screens
- [x] `ctrl+c` quits application from all screens
- [x] `esc` consistently goes back to previous state
- [x] `?` key handled (placeholder for help modal)
- [x] `m` key removed (no longer mapped)
- [x] `h` key removed (conflicted with navigation)
- [x] All 2,000+ tests pass
- [x] No race conditions

### Phases 4-5 (PENDING)
- [ ] Selection preserved when navigating back
- [ ] All 7 additional intents use `ListNavigationHandler`
- [ ] Documentation updated

### Phase 6 (PENDING)
- [ ] `?` shows help modal on all screens (integrated with BaseIntent)

---

## Rollback Plan

Each phase is independently revertible via git revert.

---

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [x] `make check-compliance` passes (for Phases 1-3)
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

---

**Last Updated**: 2026-01-10
**Status**: Phases 1-3 COMPLETED, Phases 4-6 PENDING
