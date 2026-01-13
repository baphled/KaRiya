# Workflow Documentation Status Report

**Date**: 2026-01-13  
**Assessment**: Current state of all workflow documentation  
**Goal**: Ensure all workflows match Browse Timeline quality

---

## Quick Status

| Workflow | Status | Modal Docs | Priority | Estimate |
|----------|--------|------------|----------|----------|
| **Browse Timeline** | ✅ COMPLETE | ✅ Excellent | Reference | 0h (done) |
| **Event Capture** | ⚠️ UPDATE NEEDED | ⚠️ Legacy pattern | HIGH | 3h |
| **Manage Skills** | ❌ NEEDS WORK | ❌ No modals | HIGH | 4h (2h impl + 2h docs) |
| **CV Generation** | ⚠️ UPDATE NEEDED | ❌ Minimal | MEDIUM | 2h |

**Total Work Needed**: 9 hours (7 hours high priority)

---

## Detailed Analysis

### ✅ Browse Timeline - REFERENCE IMPLEMENTATION

**File**: `docs/workflows/BROWSE_TIMELINE_WORKFLOW.md`  
**Lines**: 795  
**Quality**: ⭐⭐⭐⭐⭐ Excellent

**Strengths**:
- Complete modal documentation (5 modals)
- bubbletea-overlay integration explained
- Comprehensive keyboard reference (6 tables)
- Real-world workflow examples with timing
- Troubleshooting section
- Clear navigation patterns

**Structure**:
```
1. Overview (50 lines) - Purpose, when to use, prerequisites
2. Workflow States (200 lines) - State machine, modals, transitions
3. Step-by-Step Guide (300 lines) - 6 major workflows
4. Keyboard Reference (150 lines) - 6 comprehensive tables
5. Navigation Patterns (50 lines) - Forward, back, error recovery
6. Common Workflows (30 lines) - Real examples with timing
7. Troubleshooting (15 lines) - Specific issues
```

**Missing**: Mermaid diagram (optional, all other workflows have one)

---

### ⚠️ Event Capture - UPDATE NEEDED

**File**: `docs/workflows/EVENT_CAPTURE_WORKFLOW.md`  
**Lines**: 852  
**Quality**: ⭐⭐⭐ Good (needs modal update)

**Current State**:
- Has Mermaid diagram ✅
- 26 modal references (legacy pattern) ⚠️
- Good structure ✅
- Needs modal documentation expansion ❌

**Issues**:
1. Uses **legacy modal overlay pattern** (not bubbletea-overlay)
2. Modal sub-flows section is brief (needs expansion)
3. Doesn't explain modal architecture
4. Missing detailed keyboard reference for modals

**Required Updates**:
1. Expand modal documentation (150 lines)
   - Explain legacy pattern currently in use
   - Document 3 modals: Metadata, Burst, Fact
   - Add keyboard shortcuts per modal
2. Add migration note (50 lines)
   - Note future bubbletea-overlay migration
   - Reference Task 42 checklist
3. Expand Modal Sub-Flows (100 lines)
   - Step-by-step for each modal
   - State preservation details
4. Update Keyboard Reference (50 lines)
   - Add modal-specific tables

**Estimate**: 3 hours  
**Priority**: HIGH (users actively use this workflow)

---

### ❌ Manage Skills - NEEDS IMPLEMENTATION + DOCS

**File**: `docs/workflows/MANAGE_SKILLS_WORKFLOW.md`  
**Lines**: 1,043  
**Quality**: ⭐⭐⭐ Good (but using screens not modals)

**Current State**:
- Has Mermaid diagram ✅
- 9 states (2 should be modals) ⚠️
- Comprehensive structure ✅
- No modal documentation ❌

**Issues**:
1. Filter and Sort use **full screens** instead of modals
2. No bubbletea-overlay integration
3. Inconsistent with Browse Timeline pattern
4. More complex than necessary

**Task 42 Identified**: 2 modals needed
- SkillFilterModal (filter by category/tag)
- SkillSortModal (sort by name/usage)

**Required Work**:

**Phase 1: Implementation** (2 hours)
- [ ] Create `internal/cli/components/skill_filter_modal.go`
- [ ] Create `internal/cli/components/skill_sort_modal.go`
- [ ] Integrate into `manage_skills_intent.go`
- [ ] Remove Filter/Sort screens
- [ ] Update tests

**Phase 2: Documentation** (2 hours)
- [ ] Add Modal Overlays section (100 lines)
- [ ] Update State Machine (9 → 7 states + 2 modals)
- [ ] Rewrite Filter/Sort steps as modal workflows (150 lines)
- [ ] Update Keyboard Reference with modal tables (50 lines)

**Estimate**: 4 hours total (2h implementation + 2h documentation)  
**Priority**: HIGH (UX improvement + consistency)

---

### ⚠️ CV Generation - UPDATE NEEDED

**File**: `docs/workflows/CV_GENERATION_WORKFLOW.md`  
**Lines**: 937  
**Quality**: ⭐⭐⭐⭐ Very Good (needs modal docs)

**Current State**:
- Has Mermaid diagram ✅
- 10 states, complex workflow ✅
- Good structure ✅
- Minimal modal documentation (1 reference) ❌

**Issues**:
1. Uses StandardView modals (not bubbletea-overlay) - this is fine
2. Doesn't document modal usage
3. Missing modal keyboard shortcuts
4. No explanation of StandardView modal system

**Required Updates**:
1. Add Modal Section (50 lines)
   - Document StandardView modal usage
   - Success/Error/Loading modal patterns
2. Update Step-by-Step Guide (100 lines)
   - Add modal interaction notes
   - Document when modals appear
3. Expand Keyboard Reference (50 lines)
   - Add modal-specific shortcuts if any
4. Add Modal Integration Notes (30 lines)
   - Explain StandardView approach
   - Contrast with bubbletea-overlay

**Estimate**: 2 hours  
**Priority**: MEDIUM (workflow is functional, just needs better docs)

---

## Recommended Execution Order

### Phase 1: High Priority (7 hours)

**Week 1**:
1. **Event Capture Documentation** (3 hours)
   - Update modal documentation
   - Expand keyboard reference
   - Add migration notes

2. **Manage Skills Implementation** (2 hours)
   - Create 2 modal components
   - Integrate bubbletea-overlay
   - Update tests

**Week 2**:
3. **Manage Skills Documentation** (2 hours)
   - Document new modals
   - Update state machine
   - Rewrite affected sections

### Phase 2: Medium Priority (2 hours)

**Week 3**:
4. **CV Generation Documentation** (2 hours)
   - Add modal documentation
   - Update keyboard reference
   - Add StandardView notes

### Phase 3: Optional (1 hour)

**Week 4**:
5. **Browse Timeline Diagram** (1 hour)
   - Add Mermaid diagram for consistency
   - Update diagram generation script

---

## Success Metrics

After completion, all workflows will have:

- [x] ✅ **Browse Timeline** - Reference standard (795 lines, complete)
- [ ] ✅ **Event Capture** - Updated modal docs (900+ lines expected)
- [ ] ✅ **Manage Skills** - Modal implementation + docs (1,000+ lines expected)
- [ ] ✅ **CV Generation** - Modal docs added (1,000+ lines expected)

### Quality Criteria (All Workflows)

- [ ] Consistent structure (7 sections matching Browse Timeline)
- [ ] Complete modal documentation (technology, usage, keys)
- [ ] Comprehensive keyboard reference tables
- [ ] Real-world examples with timing estimates
- [ ] Troubleshooting section
- [ ] 800-1,000 lines comprehensive coverage

---

## Dependencies

### Before Starting

- [x] Browse Timeline workflow complete (reference)
- [x] BUBBLETEA_OVERLAY_GUIDE.md complete (700+ lines)
- [x] MODAL_PATTERNS.md updated with overlay patterns
- [x] Task 42 Modal Migration Checklist created

### During Execution

- [ ] Manage Skills: Implementation before documentation
- [ ] All: Test changes manually before documenting
- [ ] All: Verify keyboard shortcuts are accurate

---

## Related Tasks

- **Task 42**: TUI Architecture Refactor (65% complete, Browse Timeline done)
- **Task 42 Issue 4**: Apply modal patterns to other intents (Browse Timeline ✅, 4 remaining)
- **Workflow Diagram Generation**: Scripts need updating for Browse Timeline

---

## Next Action Items

1. **Review plan** - Confirm priorities and estimates
2. **Start Event Capture** - Highest priority, 3 hours
3. **Implement Manage Skills modals** - Highest priority, 2 hours
4. **Document Manage Skills** - After implementation, 2 hours
5. **Update CV Generation** - Medium priority, 2 hours

**Total**: 9 hours to bring all workflows to Browse Timeline quality standard

---

## Files to Track

**Workflow Guides**:
- `docs/workflows/BROWSE_TIMELINE_WORKFLOW.md` ✅ Complete
- `docs/workflows/EVENT_CAPTURE_WORKFLOW.md` ⚠️ Needs update
- `docs/workflows/MANAGE_SKILLS_WORKFLOW.md` ❌ Needs work
- `docs/workflows/CV_GENERATION_WORKFLOW.md` ⚠️ Needs update
- `docs/workflows/README.md` - Update after completion

**Implementation Files** (Manage Skills):
- `internal/cli/intents/manage_skills_intent.go`
- `internal/cli/components/skill_filter_modal.go` (new)
- `internal/cli/components/skill_sort_modal.go` (new)

**Reference Guides**:
- `docs/BUBBLETEA_OVERLAY_GUIDE.md` ✅ Complete
- `docs/MODAL_PATTERNS.md` ✅ Complete
- `tasks/tasks-42-tui-architecture-refactor.md` ✅ Updated
