# Task 21: Update Documentation for Strategy System

**Created**: 2026-01-07
**Status**: In Progress
**Priority**: MEDIUM
**Estimated Time**: 1-2 hours
**Related**: Task 17 (Form Refactoring), Task 19 (Test Fixes), Task 20 (Cleanup)

---

## Overview

Update documentation to reflect the new strategy system (quick/manual) instead of the old capture mode system (Timeline Journaling, CV Backfill, Manual Entry).

The old 3-mode system was replaced with a simpler 2-strategy system in Task 17, but documentation still references the old system.

---

## Files to Update

### 1. docs/CLI_GUIDE.md
**References Found**: 9 occurrences of old mode system
**Changes Needed**:
- [ ] Remove "Timeline Journaling" mode description
- [ ] Remove "CV Backfill" mode description
- [ ] Remove "Manual Entry" mode description
- [ ] Add "Quick Capture" strategy description
- [ ] Add "Manual Capture" strategy description
- [ ] Update examples to show quick vs manual
- [ ] Update troubleshooting section

### 2. docs/TUI_STANDARDS.md
**References Found**: 3 occurrences
**Changes Needed**:
- [ ] Update footer examples to remove mode references
- [ ] Update form component examples
- [ ] Update status line examples

### 3. docs/PRD_CLI.md
**References Found**: 10+ occurrences
**Changes Needed**:
- [ ] Update requirements from 3 modes to 2 strategies
- [ ] Update feature descriptions
- [ ] Update user acceptance criteria
- [ ] Update examples

### 4. AGENTS.md
**References Found**: Recent fixes section outdated
**Changes Needed**:
- [ ] Add Tasks 19 & 20 to "Recent Fixes" section
- [ ] Update current status metrics
- [ ] Update test counts (2,078 total)

---

## Strategy System Documentation

### Old System (3 Modes)
```
1. Timeline Journaling - Recent events only (within 30 days)
2. CV Backfill - Historical events (any date)
3. Manual Entry - Flexible, any date, all fields optional
```

### New System (2 Strategies)
```
1. Quick Capture - Minimal fields (just event text + date)
   - Use for rapid event logging
   - Optional fields hidden by default
   - Toggle with 't' key to show more fields

2. Manual Capture - All fields visible
   - Use for detailed event entry
   - Company, Project, Tags, Categories visible
   - Optional fields can be toggled on/off with 't' key
```

---

## Implementation Plan

### Phase 1: Update CLI_GUIDE.md (30 min)

**Current Content to Replace**:
```markdown
#### Timeline Journaling
For logging recent career events as they happen.
- Events must be within 30 days of today
- Quick capture workflow

#### CV Backfill
For importing events from existing CVs.
- Allows older dates
- No time constraints

#### Manual Entry
For manually adding individual events.
- No date constraints
- Full control over all fields
```

**New Content**:
```markdown
### Capture Strategies

KaRiya offers two capture strategies to match your workflow:

#### Quick Capture
For rapid event logging with minimal fields:
- Enter event description and date
- Optional fields hidden by default
- Press 't' to toggle additional fields (Company, Project, Tags, Categories)
- Ideal for daily journaling

Example:
```bash
kariya capture --strategy quick
```

#### Manual Capture
For detailed event entry with all fields visible:
- All fields visible by default (Event, Date, Company, Project, Tags, Categories)
- Full control over event metadata
- Press 't' to toggle optional field visibility
- Ideal for importing historical events

Example:
```bash
kariya capture --strategy manual
```

### Field Visibility

Both strategies support field toggling:
- Press 't' to show/hide optional fields (Company, Project, Tags, Categories)
- Event text and Date are always visible
- Toggle state persists during form session
```

**Tasks**:
- [ ] Replace mode descriptions with strategy descriptions
- [ ] Update examples
- [ ] Update CLI flags/options
- [ ] Update troubleshooting section

### Phase 2: Update TUI_STANDARDS.md (15 min)

**Current Content to Replace**:
```markdown
Mode: Timeline Journaling │ Status: 3/10 events │ Ctrl+C to quit
footer.SetMode("Timeline Journaling")
- ✅ Form: Shows capture mode
```

**New Content**:
```markdown
Strategy: Quick │ Status: 3/10 events │ Ctrl+C to quit
footer.SetStrategy("Quick")
- ✅ Form: Shows capture strategy
```

**Tasks**:
- [ ] Update footer examples
- [ ] Update form component examples
- [ ] Update consistency guidelines

### Phase 3: Update PRD_CLI.md (30 min)

**Current Content to Replace**:
```markdown
1. The CLI must support three capture modes:
   - **CV Backfill Mode**: For importing events from existing CVs
   - **Timeline Journaling Mode**: For real-time event logging
   - **Manual Entry Mode**: For manually adding individual events
   - Dropdown selection for capture mode
```

**New Content**:
```markdown
1. The CLI must support two capture strategies:
   - **Quick Capture**: For rapid event logging with minimal fields
   - **Manual Capture**: For detailed event entry with all fields
   - Strategy selection via flag or interactive choice
   - Field visibility toggle ('t' key) in both strategies
```

**Tasks**:
- [ ] Update requirements
- [ ] Update acceptance criteria
- [ ] Update feature descriptions
- [ ] Update examples

### Phase 4: Update AGENTS.md (15 min)

**Add to "Recent Fixes" section**:
```markdown
### Tasks 19-20: Test Fixes and Final Cleanup (2026-01-07)

**Status**: ✅ **COMPLETE - ALL TESTS PASSING, ZERO STATICCHECK WARNINGS**

#### Task 19: Fix Remaining Test Failures
**Summary**: Fixed 3 test failures after form refactoring

**Issues Fixed**:
1. **Form Field Visibility**: Changed `showOptionalFields` default to `true` to align with manual strategy
2. **Footer Text Test**: Updated test to check for "Back" instead of outdated "Retry"

**Impact**:
- All 2,078 tests passing (100% pass rate)
- 2 files modified, 3 lines changed
- Commit: `31218ec`

#### Task 20: Final Cleanup - Staticcheck
**Summary**: Removed unused functions flagged by staticcheck

**Cleaned**:
1. `BurstManagementIntent.applyFilters()` - 35 lines
2. `BurstManagementIntent.setFailed()` - 14 lines
3. `GenerateCVIntent.getStateName()` - 27 lines

**Impact**:
- Zero staticcheck warnings (was 3)
- 80 lines of dead code removed
- All tests still passing
- Commit: `0307f4f`

#### Combined Results
- ✅ All 2,078 tests passing (100%)
- ✅ Zero staticcheck warnings
- ✅ Zero race conditions
- ✅ Build successful
- ✅ Production ready
```

**Update Current Status**:
```markdown
### 🎯 Project Status
- **Total tests**: 2,078 (100% passing)
- **Build status**: ✅ Successful
- **Staticcheck**: 0 warnings
- **Race conditions**: 0 detected
- **Code coverage**: 87%+
```

**Tasks**:
- [ ] Add Tasks 19-20 to Recent Fixes
- [ ] Update project status metrics
- [ ] Update test counts

---

## Verification

### Before Documentation Updates
```bash
# Check for old mode references
rg "Timeline Journaling|CV Backfill|Manual Entry" docs/ --type md
```

### After Documentation Updates
```bash
# Should find no old mode references (except in git history or archived docs)
rg "Timeline Journaling|CV Backfill|Manual Entry" docs/ --type md

# Should find new strategy references
rg "Quick Capture|Manual Capture|strategy" docs/ --type md
```

---

## Acceptance Criteria

### Must Have
- [ ] All 4 documentation files updated
- [ ] No references to old 3-mode system (Timeline/CV/Manual)
- [ ] Clear explanation of 2-strategy system (Quick/Manual)
- [ ] Examples updated to reflect current CLI
- [ ] AGENTS.md reflects Tasks 19-20 completion

### Should Have
- [ ] Migration guide explaining old vs new system
- [ ] Troubleshooting section updated
- [ ] User-facing benefits explained

### Nice to Have
- [ ] Screenshots showing strategy selection
- [ ] Comparison table: old vs new system

---

## Migration Guide Content

Add this section to CLI_GUIDE.md:

```markdown
## Migration from Old Capture Modes

If you're familiar with the old 3-mode system, here's how it maps to the new strategy system:

| Old Mode | New Strategy | Notes |
|----------|--------------|-------|
| Timeline Journaling | Quick Capture | Same rapid entry, no date restrictions |
| CV Backfill | Manual Capture | All fields visible, import historical events |
| Manual Entry | Manual Capture | Flexible field visibility with 't' toggle |

### Key Improvements
- **Simpler**: 2 strategies instead of 3 modes
- **Flexible**: Toggle fields on/off in both strategies
- **Consistent**: Same field behavior across strategies
- **No date restrictions**: Enter events from any date in either strategy
```

---

## Success Metrics

- ✅ All documentation accurate and up-to-date
- ✅ Zero references to deprecated mode system
- ✅ Clear user guidance for strategy selection
- ✅ Examples match current implementation
- ✅ AGENTS.md shows current project status

---

## Next Steps

After documentation updates:
- **Task 22**: Add comprehensive strategy system tests
- **Final Review**: Project ready for production

---

**Last Updated**: 2026-01-07
**Author**: AI Assistant (via OpenCode)
**Status**: Ready for implementation
