# Task 38: Navigation Pattern Standardization

## Overview
- **Goal**: Standardize keyboard shortcuts, navigation behavior, and selection preservation across all 10 intents using library-based solutions
- **Time Estimate**: 8-12 days (Phases 1-5 core, Phases 6-8 optional)
- **Prerequisites**: All tests passing, understanding of intent architecture

## Session Contract Acknowledgment
- [x] Ran `make session-start` and it passed
- [x] Acknowledge and commit to following all workflow rules
- [ ] Token count: _____ (must be < 50k to start)

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
| **Key bindings** | Follow docs: `q`=quit, `h`=home, `esc`=back, `?`=help modal |
| **DOT graphs** | Skip - unnecessary overhead |
| **Breaking changes** | Acceptable - `q` and `m` behavior will change |

---

## Files to Modify

### Phase 1: KeyMaps
- [ ] Create: `internal/cli/navigation/keymaps.go`
- [ ] Create: `internal/cli/navigation/keymaps_test.go`
- [ ] Modify: `internal/cli/intents/contract.go` (add KeyMap to BaseIntent)

### Phase 2: Help Modal
- [ ] Create: `internal/cli/components/help_modal.go`
- [ ] Create: `internal/cli/components/help_modal_test.go`
- [ ] Modify: `internal/cli/intents/contract.go` (add help modal support)

### Phase 3: Intent Key Handler Standardization
- [ ] `internal/cli/intents/browse_timeline_intent.go` (2 states)
- [ ] `internal/cli/intents/fact_management_intent.go` (3 states)
- [ ] `internal/cli/intents/metadata_editor_intent.go` (3 states)
- [ ] `internal/cli/intents/bulk_operations_intent.go` (4 states)
- [ ] `internal/cli/intents/capture_event_intent.go` (4 states)
- [ ] `internal/cli/intents/import_wizard_intent.go` (5 states)
- [ ] `internal/cli/intents/burst_management_intent.go` (7 states)
- [ ] `internal/cli/intents/configure_system_intent.go` (7 states)
- [ ] `internal/cli/intents/export_artifact_intent.go` (9 states)
- [ ] `internal/cli/intents/generate_cv_intent.go` (13 states)

### Phase 4: Selection Preservation
- [ ] Modify: `internal/cli/intents/contract.go` (add selection state)
- [ ] Modify: `internal/cli/intents/result.go` (add selection metadata)
- [ ] Modify: All list-based intents

### Phase 5: ListNavigationHandler Adoption
- [ ] `internal/cli/intents/metadata_editor_intent.go`
- [ ] `internal/cli/intents/bulk_operations_intent.go`
- [ ] `internal/cli/intents/configure_system_intent.go`
- [ ] `internal/cli/intents/export_artifact_intent.go`

### Phase 6-8: Dependencies and State Machines (Optional)
- [ ] `go.mod` / `go.sum` (add stateless, bubbleo)
- [ ] All intent files (convert to stateless state machines)
- [ ] `internal/cli/app/app.go` (integrate bubbleo shell)

---

## Implementation Checklist

### Phase 1: KeyMaps Foundation (No New Dependencies)
**Time**: 1-2 days | **Risk**: Low | **Status**: IN PROGRESS

#### TDD: RED Phase
- [ ] Create test file: `internal/cli/navigation/keymaps_test.go`
- [ ] Write tests for GlobalKeyMap initialization
- [ ] Write tests for ListKeyMap initialization
- [ ] Write tests for FormKeyMap initialization
- [ ] Write tests for key matching behavior
- [ ] Write tests for help string generation
- [ ] Confirm tests FAIL (no implementation yet)

#### TDD: GREEN Phase
- [ ] Create: `internal/cli/navigation/keymaps.go`
- [ ] Implement GlobalKeyMap with Quit, Help, Back bindings
- [ ] Implement ListKeyMap with Up, Down, Select, Delete, Edit, Filter, Search bindings
- [ ] Implement FormKeyMap with NextField, PrevField, Submit, Cancel, Toggle bindings
- [ ] Implement DefaultGlobalKeyMap() constructor
- [ ] Implement DefaultListKeyMap() constructor
- [ ] Implement DefaultFormKeyMap() constructor
- [ ] Confirm all tests PASS

#### TDD: REFACTOR Phase
- [ ] Ensure consistent naming conventions
- [ ] Add ShortHelp() and FullHelp() methods for help.KeyMap interface
- [ ] Tests still pass

#### Commit
```bash
git commit -m "feat(cli): add standardized KeyMaps using bubbles/key"
```

---

### Phase 2: Help Modal Integration
**Time**: 1 day | **Risk**: Low | **Status**: PENDING

#### TDD: RED Phase
- [ ] Create test file: `internal/cli/components/help_modal_test.go`
- [ ] Write tests for HelpModal creation
- [ ] Write tests for showing/hiding modal
- [ ] Write tests for rendering with KeyMap content
- [ ] Write tests for key handling (esc/? to close)
- [ ] Confirm tests FAIL

#### TDD: GREEN Phase
- [ ] Create: `internal/cli/components/help_modal.go`
- [ ] Implement HelpModal using bubbles/help
- [ ] Implement Show() and Hide() methods
- [ ] Implement Update() for key handling
- [ ] Implement View() for rendering
- [ ] Confirm tests PASS

#### Commit
```bash
git commit -m "feat(components): add HelpModal using bubbles/help"
```

---

### Phase 3: Intent Key Handler Standardization
**Time**: 2-3 days | **Risk**: Medium | **Status**: PENDING

**Order**: Start with simplest (BrowseTimeline - 2 states), end with most complex (GenerateCV - 13 states)

---

### Phase 4: Selection Preservation
**Time**: 1-2 days | **Risk**: Medium | **Status**: PENDING

---

### Phase 5: ListNavigationHandler Adoption
**Time**: 1-2 days | **Risk**: Low | **Status**: PENDING

#### Current Status
Already using ListNavigationHandler:
- [x] `browse_timeline_intent.go`
- [x] `burst_management_intent.go`
- [x] `fact_management_intent.go`

Need to add:
- [ ] `metadata_editor_intent.go`
- [ ] `bulk_operations_intent.go`
- [ ] `configure_system_intent.go`
- [ ] `export_artifact_intent.go`

---

### Phase 6-8: Optional Future Work
- [ ] Add stateless + bubbleo dependencies
- [ ] State machine formalization
- [ ] Navigation stack with bubbleo

---

## Acceptance Criteria

### Phase 1-5 (Core)
- [ ] `q` quits application from all screens
- [ ] `?` shows help modal on all screens
- [ ] `esc` consistently goes back (preserving selection)
- [ ] `h` returns to home from applicable screens
- [ ] All 7 additional intents use `ListNavigationHandler`
- [ ] Selection preserved when navigating back
- [ ] All 2,078+ tests pass
- [ ] Documentation updated

---

## Rollback Plan

Each phase is independently revertible via git revert.

---

## Post-Task Checklist (MUST COMPLETE BEFORE NEXT TASK)
- [ ] `make check-compliance` passes
- [ ] All checkboxes above completed
- [ ] Task marked complete `[x]` in task file
- [ ] Token count: _____ (< 100k to continue)

---

**Last Updated**: 2026-01-10
**Status**: In Progress - Phase 1
