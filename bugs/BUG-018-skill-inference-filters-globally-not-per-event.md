# BUG-018: Skill inference incorrectly filters globally instead of per-event

**Status**: Solution Identified - Refactoring Required  
**Severity**: 🔴 Critical  
**Created**: 2026-02-10  
**Updated**: 2026-02-10  

---

## Bug Summary

When pressing 'i' to infer skills from an event/burst with no linked skills, the system incorrectly reports "All detected skills already in profile" if those skills exist anywhere in the global skill repository, even though they're not linked to the current event/burst.

**Impact**: Users cannot reuse existing skills across multiple events/bursts, breaking a core workflow pattern.

---

## Affected Components

- [x] Browse Timeline Intent (`internal/cli/intents/browsetimeline/helpers.go:433-479`)
- [x] Burst Management Intent (`internal/cli/intents/burst_management/handlers.go:799-839`, `helpers.go:1086-1105`)
- [x] Skill Inference Service (`internal/service/career/skillinference/detector.go:123-145`)

**NOT Affected**:
- ✅ Review Enrichment (Capture Event) - Already implements correct solution

**Related Intents/Workflows**:
- Browse Timeline → Select Event → View Skills → Press 'i'
- Burst Management → Select Burst → View Details → Press 'i'

---

## Reproduction Steps

### Browse Timeline Path
1. Launch KaRiya and create/ensure a skill exists globally (e.g., "Go")
2. Navigate to Browse Timeline
3. Select an event that mentions "Go" in its text but has NO skills linked to it
4. Press 's' to view skills for the event
5. Confirm the skills list is empty (no linked skills)
6. Press 'i' to infer skills
7. **BUG**: See error modal "All Skills Tracked - All detected skills already in profile."
8. **Expected**: Should see skill suggestion modal with "Go" available to link

### Burst Management Path
1. Ensure skills exist globally (e.g., "Go", "PostgreSQL")
2. Navigate to Burst Management
3. Select a burst with events mentioning those skills
4. Confirm no skills are linked to the burst
5. Press 'i' to infer skills from burst details
6. **BUG**: See modal "All Skills Already Tracked"
7. **Expected**: Should see skill suggestion modal

**Consistency**: Always reproducible

**Environment**:
- Worktree: `fix-skill-inference`
- Branch: `feature/fix-skill-inference-wiring`
- Go Version: 1.24+

---

## Expected Behavior

When pressing 'i', the system should:

1. Detect skills mentioned in the event/burst text (e.g., "Go", "PostgreSQL")
2. Check which skills are **already linked to THIS SPECIFIC event/burst**
3. Filter out ONLY the skills that are already linked to this entity
4. Show Skill Suggestion Modal with all detectable skills that CAN be linked
5. Allow user to accept suggestions to link them to the current entity

**Key Semantic Distinction**:
- ✅ "Existing in profile" = **"already linked to THIS event/burst"**
- ❌ NOT: **"exists anywhere in the global skill repository"**

**Use Case**: User should be able to link the same skill (e.g., "Go") to multiple different events/bursts.

---

## Actual Behavior

When pressing 'i', the system:

1. ✅ Detects skills mentioned in event/burst text correctly
2. ❌ Checks which skills exist GLOBALLY in the repository (wrong scope)
3. ❌ Filters out ALL globally-existing skills, regardless of linkage
4. ❌ Shows error: "All detected skills already in profile"
5. ❌ Prevents user from linking existing skills to new events

**Evidence**:
- Error message displayed: "All detected skills already in profile"
- Event/burst has zero linked skills (verified via 's' key view)
- Skills like "Go" exist in repository but aren't linked to current context
- User cannot reuse existing skills across multiple events

---

## Investigation Log

### [2026-02-10 19:00] - Initial Report
- User reported: pressing 'i' in browse timeline shows "All detected skills already in profile"
- User confirmed: NO skills are currently linked to the event
- User confirmed: Same issue exists in burst management
- User confirmed: Review Enrichment (capture event) does NOT have this issue

### [2026-02-10 19:15] - Intent Layer Analysis
- **Browse Timeline** (`helpers.go:440`): Error originates from `handleSkillSuggestionsLoaded()`
- **Burst Management** (`handlers.go:826`): Same pattern in `handleSkillSuggestionsLoaded()`
- Both use filtering function: `filterNewSkillSuggestions()` / `filterNewSuggestions()`
- Filtering logic itself is CORRECT - it correctly filters based on `ExistingSkillNames`
- **Problem is UPSTREAM**: The `ExistingSkillNames` list contains the wrong skills

### [2026-02-10 19:30] - Service Layer Analysis
- Reviewed `InferSkillsFromEvents()` in `detector.go:59-88`
- Line 74 calls: `existingNames, err := s.findExistingSkillNames(ctx, suggestionMap)`
- Found `findExistingSkillNames()` at lines 123-145
- **Issue**: Line 134 calls `s.skillRepo.GetByName(ctx, name)` which checks GLOBAL existence

### [2026-02-10 19:45] - Working Solution Found (Review Enrichment)
- Analyzed Capture Event intent at `/internal/cli/intents/captureevent/handlers.go:413-425`
- **KEY FINDING**: Review Enrichment OVERRIDES EventIDs before calling service:
  ```go
  skillsWithEventID[idx].EventIDs = []string{i.reviewState.Event.ID}
  ```
- This ensures skills are linked to the CURRENT event only, regardless of global existence
- Does NOT use global filtering - shows ALL detected skills
- Service layer handles deduplication and linking via EventIDs

### [2026-02-10 20:00] - Root Cause Confirmed

**Pattern Comparison**:

| Aspect | Review Enrichment ✅ | Browse Timeline ❌ | Burst Management ❌ |
|--------|---------------------|-------------------|---------------------|
| **Filtering** | None (shows all) | Global (`filterNewSkillSuggestions`) | Global (`filterNewSuggestions`) |
| **EventID Handling** | **Overrides with current event** | Uses inference EventIDs | Uses inference EventIDs |
| **Linking** | Via `CreateSkillsFromSuggestions` | Manual + service | Via service |
| **Result** | ✅ Works correctly | ❌ Filters globally | ❌ Filters globally |

---

## Root Cause

**Status**: ✅ Identified

**Cause**: Browse Timeline and Burst Management do not follow the Review Enrichment pattern of overriding EventIDs before calling `CreateSkillsFromSuggestions()`.

**Technical Details**:

### Working Solution: Review Enrichment
**File**: `internal/cli/intents/captureevent/handlers.go`  
**Lines**: 413-425

```go
// ✅ CORRECT IMPLEMENTATION
ctx := context.Background()
skillsWithEventID := make([]skillinference.SkillSuggestion, len(accepted))
for idx, skill := range accepted {
    skillsWithEventID[idx] = skill
    skillsWithEventID[idx].EventIDs = []string{i.reviewState.Event.ID}  // ✅ Override
}

_, err := i.context.SkillInferenceService.CreateSkillsFromSuggestions(ctx, skillsWithEventID)
```

**Why This Works**:
1. Accepts ALL detected skills (no global filtering)
2. Overrides `EventIDs` to point to the current event
3. Service links to the specified event only
4. Existing skills can be reused across multiple events

---

### Broken Implementation: Browse Timeline
**File**: `internal/cli/intents/browsetimeline/helpers.go`  
**Lines**: 433-448, 464-479

```go
// ❌ BROKEN: Global filtering
func (i *Intent) handleSkillSuggestionsLoaded(msg SkillSuggestionsLoadedMsg) tea.Cmd {
    if len(msg.Suggestions) == 0 {
        i.ShowErrorModal("No Skills Detected", "No skills detected from event.")
        return nil
    }
    newSuggestions := filterNewSkillSuggestions(msg.Suggestions, msg.ExistingSkillNames)  // ❌ WRONG
    if len(newSuggestions) == 0 {
        i.ShowErrorModal("All Skills Tracked", "All detected skills already in profile.")  // ❌ MISLEADING
        return nil
    }
    // ...
}

// ❌ BROKEN: Does not override EventIDs
func (i *Intent) saveSkillFromSuggestion(suggestion skillinference.SkillSuggestion) {
    // ...
    skills, err := i.context.SkillInferenceService.CreateSkillsFromSuggestions(ctx, []skillinference.SkillSuggestion{suggestion})  // ❌ Uses inference EventIDs
    // ...
    i.context.CLIEventService.LinkSkillToEvent(ctx, i.selectedEvent.ID, skills[0].ID)  // ❌ Manual linking
}
```

**Why This Fails**:
1. Filters based on global skill existence (`ExistingSkillNames`)
2. Does NOT override `EventIDs` in suggestion
3. Relies on manual linking after service call
4. Blocks reuse of existing skills

---

### Broken Implementation: Burst Management
**File**: `internal/cli/intents/burst_management/handlers.go` & `helpers.go`  
**Lines**: 799-839, 1086-1105

```go
// ❌ BROKEN: Global filtering
func (i *Intent) handleSkillSuggestionsLoaded(msg SkillSuggestionsLoadedMsg) tea.Cmd {
    // ...
    newSuggestions := filterNewSuggestions(msg.Suggestions, msg.ExistingSkillNames)  // ❌ WRONG
    if len(newSuggestions) == 0 {
        i.ShowSuccessModal("All Skills Already Tracked", ...)  // ❌ MISLEADING
        return nil
    }
    // ...
}

// ❌ BROKEN: Does not override EventIDs
func (i *Intent) saveSkillFromSuggestion(ctx context.Context, service SkillInferenceService, suggestion skillinference.SkillSuggestion) tea.Cmd {
    return func() tea.Msg {
        skills, err := service.CreateSkillsFromSuggestions(ctx, []skillinference.SkillSuggestion{suggestion})  // ❌ Uses inference EventIDs
        // ...
    }
}
```

**Why This Fails**:
1. Filters based on global skill existence
2. Does NOT override `EventIDs` with burst event IDs
3. Relies on inference EventIDs (may be incorrect)
4. Blocks reuse of existing skills

---

## Fix Strategy

**Approach**: Apply Review Enrichment pattern to Browse Timeline and Burst Management

### Solution: Override EventIDs Before Service Call

**For Browse Timeline**:
```go
func (i *Intent) saveSkillFromSuggestion(suggestion skillinference.SkillSuggestion) {
    if i.context.SkillInferenceService == nil || i.selectedEvent == nil {
        return
    }
    ctx := i.getContext()
    
    // ✅ Override EventIDs to ensure skill is linked to THIS event only
    suggestion.EventIDs = []string{i.selectedEvent.ID}
    
    _, err := i.context.SkillInferenceService.CreateSkillsFromSuggestions(ctx, []skillinference.SkillSuggestion{suggestion})
    if err != nil {
        i.ShowErrorModal("Skill Creation Failed", err.Error())
        return
    }
}
```

**For Burst Management**:
```go
func (i *Intent) saveSkillFromSuggestion(ctx context.Context, service SkillInferenceService, suggestion skillinference.SkillSuggestion) tea.Cmd {
    return func() tea.Msg {
        if err := ctx.Err(); err != nil {
            return SkillsCreatedMsg{Error: err}
        }

        // ✅ Override EventIDs to ensure skill is linked to ALL burst events
        if i.selectedBurst != nil {
            suggestion.EventIDs = i.selectedBurst.EventIDs
        }

        skills, err := service.CreateSkillsFromSuggestions(ctx, []skillinference.SkillSuggestion{suggestion})
        if err != nil {
            return SkillsCreatedMsg{Error: err}
        }

        return SkillsCreatedMsg{Skills: skills}
    }
}
```

**Remove Global Filtering**:
Both intents should show ALL detected skills, not filter by global existence.

---

## Files to Change

### Browse Timeline
- [ ] `internal/cli/intents/browsetimeline/helpers.go:433-448`
  - Remove `filterNewSkillSuggestions()` call
  - Remove "All Skills Tracked" error modal
  - Show ALL detected skills

- [ ] `internal/cli/intents/browsetimeline/helpers.go:464-479`
  - Override `suggestion.EventIDs = []string{i.selectedEvent.ID}`
  - Remove manual `LinkSkillToEvent()` call
  - Let service handle linking

- [ ] `internal/cli/intents/browsetimeline/helpers.go:450-462`
  - (Optional) Delete unused `filterNewSkillSuggestions()` function

### Burst Management
- [ ] `internal/cli/intents/burst_management/handlers.go:799-839`
  - Remove `filterNewSuggestions()` call
  - Remove "All Skills Already Tracked" success modal
  - Show ALL detected skills

- [ ] `internal/cli/intents/burst_management/helpers.go:1086-1105`
  - Override `suggestion.EventIDs = i.selectedBurst.EventIDs`
  - Let service handle linking to all burst events

- [ ] `internal/cli/intents/burst_management/helpers.go:27-47`
  - (Optional) Delete unused `filterNewSuggestions()` function

### Tests
- [ ] `internal/cli/intents/browsetimeline/helpers_test.go`
  - Remove tests validating global filtering
  - Add test: `saveSkillFromSuggestion overrides EventIDs`

- [ ] `internal/cli/intents/burst_management/helpers_test.go`
  - Update line 746: Remove global filtering expectations
  - Add test: `saveSkillFromSuggestion overrides EventIDs with burst events`

### BDD Scenarios
- [ ] `features/browse_timeline.feature` - Add skill reuse scenarios
- [ ] `features/burst_management.feature` - Add skill reuse scenarios

---

## Testing Plan

### Phase 1: Unit Tests

**Browse Timeline**:
- [ ] Test: `saveSkillFromSuggestion` overrides EventIDs to current event
- [ ] Test: `handleSkillSuggestionsLoaded` shows all suggestions (no filtering)
- [ ] Test: Skill creation succeeds when skill exists globally but not linked

**Burst Management**:
- [ ] Test: `saveSkillFromSuggestion` overrides EventIDs to all burst events
- [ ] Test: `handleSkillSuggestionsLoaded` shows all suggestions (no filtering)
- [ ] Test: Skill is linked to ALL burst events after acceptance

---

### Phase 2: BDD Scenarios

**Browse Timeline** (`features/browse_timeline.feature`):

```gherkin
@bug-018 @critical
Scenario: Infer and link existing skill to new event
  Given I have a skill "Go" in the global repository
  And I have an event "Event A" with "Go" linked to it
  And I have an event "Event B: Built CLI tool in Go" with no skills linked
  When I navigate to browse timeline
  And I select "Event B"
  And I press "s" to view skills
  Then I should see 0 skills linked
  When I press "i" to infer skills
  Then I should see the skill suggestion modal
  And the modal should contain "Go"
  When I accept the "Go" suggestion
  Then "Go" should be linked to "Event B"
  And "Go" should still be linked to "Event A"

@bug-018 @critical
Scenario: Reuse same skill across multiple events
  Given I have a skill "Go" in the global repository
  And I have an event "Event A: Built API in Go"
  And I have an event "Event B: Created CLI tool in Go"
  And "Go" is linked to "Event A"
  And "Go" is NOT linked to "Event B"
  When I navigate to browse timeline
  And I select "Event B"
  And I press "s" to view event skills
  Then I should see 0 linked skills
  When I press "i" to infer skills
  Then I should see the skill suggestion modal
  And the modal should contain "Go" as a suggestion
  When I accept the "Go" suggestion
  Then "Go" should be linked to "Event B"
  And "Go" should still be linked to "Event A"
  And the skill "Go" should have 2 linked events in the database
```

**Burst Management** (`features/burst_management.feature`):

```gherkin
@bug-018 @critical
Scenario: Infer and link existing skill to burst events
  Given I have a skill "Go" in the global repository
  And I have a confirmed burst "Backend Development Sprint"
  And the burst contains 3 events mentioning "Go"
  And none of the burst events have skills linked
  When I navigate to burst management
  And I select the burst
  And I press "i" to infer skills
  Then I should see the skill suggestion modal
  And the modal should contain "Go"
  When I accept the "Go" suggestion
  Then all 3 events in the burst should have "Go" linked
```

---

### Phase 3: Manual Testing Checklist

**Browse Timeline**:
- [ ] Event with no skills, press 'i' → see ALL detected suggestions
- [ ] Event with 1 skill, press 'i' → see that skill again if detected (no filtering)
- [ ] Accept suggestion → skill linked to current event only
- [ ] Link skill to Event A, then Event B → both events show skill independently

**Burst Management**:
- [ ] Burst with no skills, press 'i' → see ALL detected suggestions
- [ ] Accept suggestion → skill linked to ALL burst events
- [ ] Same skill in multiple bursts → independent management

**Regression**:
- [ ] Review Enrichment still works (capture event flow)
- [ ] Manual skill creation still works
- [ ] Skill picker still works
- [ ] All existing tests pass

---

## Verification Checklist

### Code Quality
- [ ] Fix implemented following Review Enrichment pattern
- [ ] All unit tests passing (`make test`)
- [ ] All BDD scenarios passing (`make bdd`)
- [ ] No race conditions (`make test` with `-race`)
- [ ] Code coverage maintained (>80%)
- [ ] Linting passing (`make vet && make staticcheck`)
- [ ] Architecture compliance (`make check-compliance`)

### Functionality
- [ ] Bug no longer reproduces in Browse Timeline
- [ ] Bug no longer reproduces in Burst Management
- [ ] Can reuse skills across multiple events
- [ ] Can reuse skills across multiple bursts
- [ ] Review Enrichment still works correctly
- [ ] Error messages accurate and helpful

### Documentation
- [ ] Code comments explain EventID override pattern
- [ ] Bug report updated with resolution details
- [ ] Commit messages reference BUG-018

### Compliance
- [ ] Follows Go idioms and project standards
- [ ] Atomic commits with clear messages
- [ ] AI attribution via `make ai-commit`
- [ ] No breaking changes

---

## Resolution Summary

**Status**: Ready for Implementation

**Fix Description**: 
Apply the Review Enrichment pattern (from Capture Event intent) to Browse Timeline and Burst Management intents. The solution requires two changes:

1. **Remove global filtering**: Show ALL detected skills, don't filter by global existence
2. **Override EventIDs**: Before calling `CreateSkillsFromSuggestions()`, override the suggestion's EventIDs with the current event/burst events

This allows users to reuse existing skills across multiple events/bursts while maintaining correct linkage.

**Commits**: TBD

**Pull Request**: TBD

---

## Follow-Up Actions

- [ ] Implement fix following refactoring plan above
- [ ] Write regression tests FIRST (TDD)
- [ ] Update BDD scenarios
- [ ] Verify all tests pass
- [ ] Create PR targeting `next` branch
- [ ] Consider documenting the "EventID override pattern" in architecture docs

---

## Related Files

### Working Reference (Review Enrichment)
- `internal/cli/intents/captureevent/handlers.go:413-425` - **CORRECT PATTERN**

### Files to Fix (Browse Timeline)
- `internal/cli/intents/browsetimeline/helpers.go:433-448` - Remove global filtering
- `internal/cli/intents/browsetimeline/helpers.go:464-479` - Override EventIDs
- `internal/cli/intents/browsetimeline/helpers.go:450-462` - Delete unused function

### Files to Fix (Burst Management)
- `internal/cli/intents/burst_management/handlers.go:799-839` - Remove global filtering
- `internal/cli/intents/burst_management/helpers.go:1086-1105` - Override EventIDs
- `internal/cli/intents/burst_management/helpers.go:27-47` - Delete unused function

### Test Files
- `internal/cli/intents/browsetimeline/helpers_test.go` - Update tests
- `internal/cli/intents/burst_management/helpers_test.go` - Update tests
- `features/browse_timeline.feature` - Add BDD scenarios
- `features/burst_management.feature` - Add BDD scenarios

### Documentation
- `tasks/tasks-55-event-skill-management.md` - Original feature spec
- `docs/design/SKILL_INFERENCE_UI_DESIGN.md` - UI design rationale

---

## References

- **Related bugs**: None
- **Related tasks**: tasks-55 (Event Skill Management), tasks-47 (Skill Inference Service)
- **Working solution**: Review Enrichment (Capture Event) intent
- **Pattern**: EventID override before service call
- **Database**: `event_skills` junction table

---

**Last Updated**: 2026-02-10  
**Updated By**: AI Agent (Claude)  
**Status**: Ready for implementation - Solution identified, refactoring plan documented
