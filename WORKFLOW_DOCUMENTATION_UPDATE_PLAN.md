# Workflow Documentation Update Plan

**Date**: 2026-01-13  
**Goal**: Update all workflow guides to follow Browse Timeline pattern  
**Reference**: `docs/workflows/BROWSE_TIMELINE_WORKFLOW.md` (795 lines)  
**Status**: Planning phase

---

## Current State Analysis

### Completed Workflows

| Workflow | Lines | Diagram | Modal Docs | bubbletea-overlay | Status |
|----------|-------|---------|------------|-------------------|--------|
| Browse Timeline | 795 | ❌ | ✅ (113 refs) | ✅ (6 refs) | ✅ **COMPLETE** (reference) |
| CV Generation | 937 | ✅ | ❌ (1 ref) | ❌ | ⚠️ Needs modal update |
| Event Capture | 852 | ✅ | ⚠️ (26 refs) | ❌ | ⚠️ Needs modal update |
| Manage Skills | 1,043 | ✅ | ❌ (1 ref) | ❌ | ⚠️ Needs modal + diagram |

### Modal Implementation Status

| Intent | Current Modals | Need Documentation Update | Need Implementation |
|--------|----------------|---------------------------|---------------------|
| **BrowseTimeline** | 5 (all bubbletea-overlay) | ✅ Done | ✅ Done |
| **CaptureEvent** | 3 (legacy overlay pattern) | ✅ Yes | ⚠️ Consider migration |
| **GenerateCV** | 0 (uses StandardView modals) | ⚠️ Yes | ⚠️ ProfileSelector modal? |
| **ManageSkills** | 0 (uses screens) | ⚠️ Yes | ❌ Need 2 modals |

---

## Browse Timeline Pattern (Reference)

### Structure (795 lines)

```
1. Overview (50 lines)
   - Purpose
   - When to Use
   - Prerequisites
   - Complexity

2. Workflow States & Navigation (200 lines)
   - State Machine Overview (narrative)
   - Main State
   - Modal Overlays (5 modals with bubbletea-overlay)
   - State Transitions (with timing)

3. Step-by-Step Guide (300 lines)
   - 6 major workflows with screenshots
   - Each step: Purpose, Actions, What Happens, Tips

4. Complete Keyboard Reference (150 lines)
   - One table per state/modal
   - Comprehensive key mappings

5. Navigation Patterns (50 lines)
   - Forward navigation
   - Back navigation
   - Error recovery

6. Common Workflows (30 lines)
   - Real-world examples with timing

7. Troubleshooting (15 lines)
   - Specific issues and solutions
```

### Key Features

1. **Modal-Centric Documentation**
   - Each modal has dedicated section
   - bubbletea-overlay integration explained
   - Keyboard shortcuts per modal
   - State transitions clearly documented

2. **User-Focused Language**
   - "What happens" explanations
   - Real-world workflow examples
   - Tips and best practices
   - Troubleshooting section

3. **Technical Accuracy**
   - References actual implementation files
   - Accurate line numbers
   - Test coverage statistics
   - State transition timing

4. **Visual Aids**
   - Narrative state machine (not just diagram)
   - Table-based keyboard reference
   - Clear state transition flow

---

## Update Plan by Workflow

### 1. CV Generation Workflow (Priority: Medium)

**Current**: 937 lines, has Mermaid diagram, minimal modal docs  
**Status**: Functional but needs modal documentation update  
**Estimate**: 2 hours

#### Analysis
- Uses StandardView modals (not bubbletea-overlay)
- 10 states, complex workflow
- Already has good structure and diagram
- Needs modal integration documentation

#### Required Updates

1. **Add Modal Section** (50 lines)
   - Document StandardView modal usage
   - Success/Error/Loading modal patterns
   - When modals appear in workflow

2. **Update Step-by-Step Guide** (100 lines)
   - Add modal screenshots/descriptions
   - Document modal interactions
   - Add timing estimates per step

3. **Expand Keyboard Reference** (50 lines)
   - Add modal-specific keys if any
   - Document StandardView modal shortcuts

4. **Add Modal Integration Notes** (30 lines)
   - Explain StandardView modal system
   - Link to STANDARDVIEW_GUIDE.md
   - Contrast with bubbletea-overlay approach

**Files to Update**:
- `docs/workflows/CV_GENERATION_WORKFLOW.md`

**New Sections**:
- "Modal System" subsection in Workflow States
- Modal interaction notes in Step-by-Step Guide

---

### 2. Event Capture Workflow (Priority: High)

**Current**: 852 lines, has Mermaid diagram, 26 modal references (legacy pattern)  
**Status**: Needs modal pattern update  
**Estimate**: 3 hours

#### Analysis
- Uses legacy modal overlay pattern (not bubbletea-overlay)
- 3 modals: Metadata, Burst, Fact editing
- Already has modal documentation but using old pattern
- Form-heavy workflow with modal sub-flows

#### Required Updates

1. **Update Modal Documentation** (150 lines)
   - Explain current legacy pattern
   - Document 3 modal types in detail
   - Add keyboard shortcuts per modal
   - Add timing estimates

2. **Add Migration Note** (50 lines)
   - Note: "Future migration to bubbletea-overlay planned"
   - Reference Task 42 Modal Migration Checklist
   - Explain benefits of future migration

3. **Expand Modal Sub-Flows Section** (100 lines)
   - Currently brief, needs expansion
   - Add step-by-step for each modal
   - Add screenshots/examples
   - Document modal state preservation

4. **Update Keyboard Reference** (50 lines)
   - Add detailed modal shortcuts
   - Document form navigation keys
   - Add modal-specific table sections

**Files to Update**:
- `docs/workflows/EVENT_CAPTURE_WORKFLOW.md`

**New Sections**:
- "Modal System Architecture" (legacy pattern explanation)
- "Future Modal Migration" note
- Expanded "Modal Sub-Flows" with step-by-step

---

### 3. Manage Skills Workflow (Priority: High)

**Current**: 1,043 lines, has Mermaid diagram, minimal modal docs  
**Status**: Needs modal implementation AND documentation  
**Estimate**: 4 hours (2 hours implementation + 2 hours docs)

#### Analysis
- Currently uses screens for filter/sort (no modals)
- 9 states, medium complexity
- Good structure but needs modal migration
- Task 42 identifies 2 modals needed: SkillFilterModal, SkillSortModal

#### Required Updates

**Phase 1: Implementation** (2 hours) - Separate task
- [ ] Create SkillFilterModal component
- [ ] Create SkillSortModal component
- [ ] Integrate bubbletea-overlay into ManageSkillsIntent
- [ ] Update tests
- See Task 42 Modal Migration Checklist

**Phase 2: Documentation** (2 hours)
1. **Add Modal Section** (100 lines)
   - Document 2 new modals (Filter, Sort)
   - bubbletea-overlay integration
   - Keyboard shortcuts per modal
   - State transitions with modals

2. **Update State Machine** (50 lines)
   - Change Filter/Sort from states to modals
   - Update diagram if needed
   - Update state count (9 → 7 states + 2 modals)

3. **Update Step-by-Step Guide** (150 lines)
   - Rewrite Filter Menu as "Filtering Skills (Modal)"
   - Rewrite Sort Menu as "Sorting Skills (Modal)"
   - Add modal interaction steps
   - Add timing estimates

4. **Update Keyboard Reference** (50 lines)
   - Remove Filter/Sort state sections
   - Add Filter/Sort modal sections
   - Update root state keys (f for filter modal, s for sort modal)

**Files to Update**:
- `internal/cli/intents/manage_skills_intent.go` (implementation)
- `internal/cli/components/skill_filter_modal.go` (new)
- `internal/cli/components/skill_sort_modal.go` (new)
- `docs/workflows/MANAGE_SKILLS_WORKFLOW.md` (documentation)

**New Sections**:
- "Modal Overlays" in Workflow States
- "Filter Modal" and "Sort Modal" in Step-by-Step
- Modal sections in Keyboard Reference

---

### 4. Browse Timeline Workflow (Reference)

**Current**: 795 lines, no diagram, complete modal docs  
**Status**: ✅ **COMPLETE** - Reference implementation  
**Estimate**: 1 hour (optional diagram)

#### Optional Enhancement

Add Mermaid state machine diagram to match other workflows:

1. **Add State Machine Diagram** (50 lines)
   - Create Mermaid diagram showing 1 state + 5 modal overlays
   - Place after "State Machine Overview" narrative
   - Show modal triggers and returns

2. **Diagram Generation Script** (30 min)
   - Add to `scripts/generate_workflow_diagrams.sh`
   - Auto-generate from intent implementation
   - Ensure accuracy with actual code

**Files to Update**:
- `docs/workflows/BROWSE_TIMELINE_WORKFLOW.md` (optional diagram)
- `scripts/generate_workflow_diagrams.sh` (add Browse Timeline)

---

## Standardization Checklist

All workflow guides should include:

### Required Sections

- [ ] **Overview**
  - Purpose (what the workflow does)
  - When to Use (user scenarios)
  - Prerequisites (required data/setup)
  - Complexity rating (⭐⭐⭐)

- [ ] **Workflow States & Navigation**
  - State Machine Diagram (Mermaid) OR Narrative
  - State/Modal summary table
  - State transitions with timing

- [ ] **Modal Documentation** (if applicable)
  - List of all modals with purpose
  - Modal overlay technology (bubbletea-overlay or StandardView or legacy)
  - Integration pattern (code examples if complex)
  - Link to modal component guide

- [ ] **Step-by-Step Guide**
  - One section per major step
  - Purpose, Actions, What Happens, Tips
  - Timing estimates per step
  - Real-world workflow examples

- [ ] **Complete Keyboard Reference**
  - One table per state/modal
  - Universal shortcuts at top
  - Context-specific shortcuts per section
  - Consistent key format (e.g., `Enter`, `Esc`, `Ctrl+C`)

- [ ] **Navigation Patterns**
  - Forward navigation explanation
  - Back navigation with Esc
  - Error recovery procedures

- [ ] **Common Workflows**
  - 3-5 real-world examples
  - Timing estimates
  - Step sequence

- [ ] **Troubleshooting**
  - 5-10 specific issues
  - Solutions and workarounds
  - Links to related docs

### Optional but Recommended

- [ ] Screenshots (if available)
- [ ] Video walkthrough link
- [ ] Performance notes (expected timing)
- [ ] Accessibility notes
- [ ] Mobile/small terminal notes

---

## Execution Order

### Phase 1: High Priority (7 hours)
1. **Event Capture** (3 hours) - Update modal documentation for existing modals
2. **Manage Skills** (4 hours) - Implement 2 modals + documentation

### Phase 2: Medium Priority (2 hours)
3. **CV Generation** (2 hours) - Add modal documentation for StandardView modals

### Phase 3: Optional (1 hour)
4. **Browse Timeline** (1 hour) - Add Mermaid diagram for consistency

**Total Estimate**: 10 hours (7 hours required, 3 hours optional)

---

## Success Criteria

### All Workflow Guides Will Have

1. ✅ Consistent structure matching Browse Timeline pattern
2. ✅ Complete modal documentation (technology, usage, keys)
3. ✅ Comprehensive keyboard reference tables
4. ✅ Real-world workflow examples with timing
5. ✅ Troubleshooting section
6. ✅ State machine visualization (diagram or narrative)
7. ✅ 800-1,000 lines comprehensive documentation

### Quality Checklist

- [ ] User can complete workflow using only the guide
- [ ] All keyboard shortcuts documented and tested
- [ ] Modal interactions clearly explained
- [ ] Timing estimates are realistic
- [ ] Troubleshooting covers common issues
- [ ] Links to related documentation work
- [ ] Technical details are accurate (file paths, line numbers)
- [ ] Code examples compile and run

---

## Related Documentation

**Reference Guides**:
- `docs/workflows/BROWSE_TIMELINE_WORKFLOW.md` - Reference implementation
- `docs/BUBBLETEA_OVERLAY_GUIDE.md` - Modal overlay library guide
- `docs/MODAL_PATTERNS.md` - Modal usage patterns
- `docs/STANDARDVIEW_GUIDE.md` - StandardView modal system
- `tasks/tasks-42-tui-architecture-refactor.md` - Modal Migration Checklist

**Implementation Guides**:
- `docs/TUI_DEVELOPER_GUIDE.md` - TUI development guide
- `docs/TUI_STANDARDS.md` - TUI design standards
- `docs/KEYBOARD_SHORTCUTS_GUIDE.md` - Keyboard shortcuts reference

---

## Next Steps

1. **Review this plan** with stakeholders
2. **Prioritize workflows** based on user needs
3. **Execute Phase 1** (Event Capture + Manage Skills)
4. **Execute Phase 2** (CV Generation)
5. **Execute Phase 3** (Browse Timeline diagram - optional)
6. **Update README.md** with workflow completion status

---

## Notes

- **Browse Timeline is the gold standard** - all workflows should match its quality and structure
- **Modal documentation is critical** - users need to understand modal interactions
- **Manage Skills requires implementation first** - documentation follows implementation
- **Timing estimates are conservative** - may complete faster with templates
- **Consistency is key** - users should recognize the same structure across all guides
