# CaptureEvent Intent Integration Status

**Last Updated**: 2026-01-14  
**Status**: 🚧 **Phase 2 Complete** - Screens created, handlers implemented, wiring pending

---

## Overview

CaptureEvent intent migration to screens architecture. This document tracks the integration status and provides guidance for completing the wiring phase.

---

## Completed Work

### ✅ Phase 1: Screen Creation (100% Complete)

Created 4 screens with comprehensive tests:

1. **StrategySelectScreen** (Commit: 5f66a9d)
   - 26 tests (100% passing)
   - Generic BaseSelectScreen[CaptureStrategy]
   - Strategy selection: Quick vs Manual

2. **EventFormScreen** (Commit: 429b819)
   - 21 tests (100% passing)
   - Wraps existing CaptureForm model
   - Supports Quick/Manual strategies

3. **EventReviewScreen** (Commit: 465b591)
   - 13 tests (100% passing)
   - Displays event + bursts + facts
   - Navigation actions for editing

4. **EventSubmitScreen** (Commit: c145212)
   - 20 tests (100% passing)
   - Async submission with progress
   - Error handling

**Total**: 80 tests, 1,662 lines of code

### ✅ Phase 2: ScreenResultHandler Implementation (100% Complete)

Added handler interface to CaptureEventIntent (Commit: 3234c06):

- ✅ ScreenResultHandler interface compliance
- ✅ 4 handler methods implemented
- ✅ 4 transition helpers (stubs)
- ✅ All 22 existing tests passing

---

## Pending Work

### 🚧 Phase 3: Screen Wiring (0% Complete)

**Blocker**: Import cycle between `intents` and `screens/capture`
- `screens/capture` imports `intents` (for CaptureStrategy type)
- `intents` would import `screens/capture` (for screen constructors)
- This creates a circular dependency

**Solutions** (choose one):

#### Option 1: Move CaptureStrategy to Shared Package (Recommended)
```
internal/cli/types/
  └── capture_types.go  // CaptureStrategy, constants
```
- Both packages import `types`
- No cycle
- Clean separation of concerns

#### Option 2: Use String Strategy
- Change screens to accept `string` instead of `CaptureStrategy`
- Convert at intent boundary
- Less type-safe but simpler

#### Option 3: Factory Pattern
- Create factory in `screens/capture`
- Intent calls factory methods
- Screens package owns creation logic

---

## Integration Checklist

### Step 1: Resolve Import Cycle
- [ ] Choose solution (Option 1, 2, or 3)
- [ ] Implement chosen solution
- [ ] Verify no import cycles: `go build ./internal/cli/...`
- [ ] Run tests: `ginkgo ./internal/cli/intents/`

### Step 2: Update Init() Method
- [ ] Set `useScreens = true`
- [ ] Create initial StrategySelectScreen
- [ ] Assign to `activeScreen`
- [ ] Test: Intent initializes with screen

### Step 3: Update Update() Method
- [ ] Check `useScreens` flag
- [ ] If true: Delegate to `activeScreen.Update(msg)`
- [ ] Handle result via `handleScreenResult()`
- [ ] If false: Use legacy code path
- [ ] Test: Screen delegation works

### Step 4: Update View() Method
- [ ] Check `useScreens` flag
- [ ] If true: Return `activeScreen.View()`
- [ ] If false: Use legacy view
- [ ] Test: Screen renders correctly

### Step 5: Implement Transition Helpers
- [ ] `transitionToStrategyScreen()` - Create StrategySelectScreen
- [ ] `transitionToFormScreen()` - Create EventFormScreen
- [ ] `transitionToReviewScreen()` - Create EventReviewScreen
- [ ] `transitionToSubmitScreen()` - Create EventSubmitScreen
- [ ] Pass terminal info, theme, logo to all screens
- [ ] Call `Init()` on each screen
- [ ] Test: Transitions create correct screens

### Step 6: Test Integration
- [ ] Add tests for screen creation
- [ ] Add tests for transition logic
- [ ] Add E2E test for complete workflow
- [ ] Verify legacy path still works (useScreens=false)
- [ ] Run full test suite: `ginkgo ./internal/cli/intents/`

### Step 7: Enable Screens by Default
- [ ] Set `useScreens = true` in `NewCaptureEventIntent`
- [ ] Remove legacy code path (if all tests pass)
- [ ] Update documentation
- [ ] Run full regression suite

---

## Code Locations

### Screens
```
internal/cli/screens/capture/
├── strategy_select.go
├── strategy_select_test.go
├── event_form_screen.go
├── event_form_screen_test.go
├── event_review_screen.go
├── event_review_screen_test.go
├── event_submit_screen.go
└── event_submit_screen_test.go
```

### Intent
```
internal/cli/intents/
├── capture_event.go           # Types and states
├── capture_event_intent.go    # Intent implementation (handlers added)
└── capture_event_test.go      # Tests (need screen tests)
```

---

## Testing Strategy

### Unit Tests (Screens)
- ✅ All 80 screen tests passing
- Located in `*_test.go` files

### Integration Tests (Intent + Screens)
- 🚧 Need to add tests for:
  - Screen creation in transitions
  - Result handling
  - State preservation
  - Error propagation

### E2E Tests (Complete Workflow)
- 🚧 Need to add test for:
  - Strategy selection → Form → Review → Submit
  - Cancellation at each step
  - Back navigation with context preservation

---

## Migration Strategy

### Phased Approach (Recommended)

**Phase 3A**: Wire Form, Review, Submit (No import cycle)
- These screens don't need CaptureStrategy import
- Can be wired immediately
- Strategy selection stays legacy for now

**Phase 3B**: Resolve StrategySelectScreen cycle
- Implement chosen solution (Option 1, 2, or 3)
- Wire StrategySelectScreen
- Complete migration

### Benefits
- Incremental progress
- Lower risk
- Can test partially migrated intent

---

## Performance Considerations

### Screen Creation Overhead
- Screens are lightweight (BaseScreen + data)
- Creation cost: <1ms per screen
- Acceptable for UI interactions

### Memory Usage
- One active screen at a time
- Previous screen garbage collected
- Typical memory: ~10KB per screen

### Test Performance
- 80 screen tests: 0.094 seconds
- 22 intent tests: 0.069 seconds
- Total: <0.2 seconds (excellent)

---

## Known Issues

### Import Cycle (High Priority)
- **Issue**: Cannot import `screens/capture` in `intents`
- **Impact**: Blocks screen wiring
- **Solutions**: See "Pending Work" section above
- **Workaround**: Use legacy code path (useScreens=false)

### TODO Comments
- `capture_event_intent.go:1386` - "Add burst/fact enrichment here if enabled"
  - Not blocking, feature enhancement
  - Can be implemented after wiring complete

---

## Next Session Plan

1. **Decide on import cycle solution** (15 min)
2. **Implement chosen solution** (30 min)
3. **Wire screens into intent** (60 min)
4. **Add integration tests** (30 min)
5. **Enable screens by default** (15 min)

**Total Estimated Time**: 2.5 hours

---

## References

- **Task File**: `tasks/tasks-42-tui-architecture-refactor.md`
- **Pattern Library**: `docs/development/INTENT_PATTERNS_LIBRARY.md`
- **ScreenResultHandler**: `internal/cli/intents/screen_result_behavior.go`
- **Base Screens**: `internal/cli/screens/base/*.go`

---

## Success Criteria

Integration is complete when:
- ✅ All 4 screens wired into intent
- ✅ `useScreens = true` by default
- ✅ All tests passing (screen + intent)
- ✅ E2E workflow functional
- ✅ Legacy code removed
- ✅ No import cycles
- ✅ Documentation updated

---

**Status**: Ready for Phase 3 (Screen Wiring)  
**Blocker**: Import cycle resolution (choose solution and implement)  
**Next Steps**: See "Next Session Plan" above
