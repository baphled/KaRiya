---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Task List: Cleanup Orphaned Views & Refactor Components for Intent System

**PRD Reference**: `tasks/prd-cleanup-and-component-refactoring.md`

**Purpose**: Remove unused code, fix broken references, extract reusable components, and consolidate shortcut systems to prepare the codebase for aggressive intent-based replacement.

**Status**: 🚀 READY FOR IMPLEMENTATION

**Version**: 1.0 - Initial Task Generation (2026-01-03)

**Effort Estimate**: 10 hours

**Priority**: HIGH (Prerequisite for Aggressive Replacement)

---

## Relevant Files

### Files to be Deleted
- `internal/cli/models/tutorial.go` - Orphaned tutorial model (never used)
- `internal/cli/models/view_event_with_facts.go` - Orphaned event details view (duplicate of DetailsModel)
- `internal/cli/models/source_event_tracer.go` - Orphaned event lineage tracker (never used)

### Files to be Modified/Cleaned
- `internal/cli/models/view_event.go` - Partially orphaned (keep EditEventMsg and ConfirmationDialog, remove unused types)
- `internal/cli/app/app.go` - Fix broken ImportProgressModel reference

### Files to be Created (Documentation)
- `docs/SUPPORTING_COMPONENTS_REFERENCE.md` - Document all 23 supporting components with usage counts
- `docs/COMPONENT_EXTRACTION_PLAN.md` - Tier-based extraction plan for components
- `docs/UNIFIED_SHORTCUT_SYSTEM_DESIGN.md` - Design for consolidated shortcut system
- `docs/SCREEN_TO_INTENT_MAPPING.md` - Complete mapping of 25 screens to intents
- `docs/CLEANUP_SUMMARY.md` - Summary of cleanup operations performed

### Test Files
- Tests for modified components (if applicable)
- Verification tests for removed files (grep searches)

### Reference Files (No Changes)
- `docs/UNINTEGRATED_VIEWS_AUDIT.md` - Source audit (reference only)
- `docs/APP_GO_AGGRESSIVE_REPLACEMENT_PLAN.md` - Main plan (reference only)
- `internal/cli/models/` - All 52 model files (for reference during audit)

### Notes

- All file deletions should be verified with `grep` searches to ensure no orphaned references remain
- Component documentation should include reference counts from the audit
- Screen-to-intent mapping should verify all 25 integrated screens are accounted for
- No changes to functionality of integrated screens or supporting components
- All tests should continue to pass after cleanup

---

## Tasks

- [ ] 1.0 Remove All Orphaned Views (Phase 1)
  - [ ] 1.1 Delete TutorialModel (tutorial.go)
  - [ ] 1.2 Delete ViewEventWithFactsModel (view_event_with_facts.go)
  - [ ] 1.3 Delete SourceEventTracerModel (source_event_tracer.go)
  - [ ] 1.4 Clean up view_event.go (keep used types, remove unused)
  - [ ] 1.5 Verify no orphaned references remain
  - [ ] 1.6 Run tests and commit

- [ ] 2.0 Fix Missing ImportProgressModel (Phase 2)
  - [ ] 2.1 Check git history for import_progress.go
  - [ ] 2.2 Restore from git or create new implementation
  - [ ] 2.3 Update app.go references if needed
  - [ ] 2.4 Verify code compiles without errors
  - [ ] 2.5 Run tests and commit

- [ ] 3.0 Extract and Document Supporting Components (Phase 3)
  - [ ] 3.1 Document all 23 supporting components
  - [ ] 3.2 Create SUPPORTING_COMPONENTS_REFERENCE.md
  - [ ] 3.3 Create COMPONENT_EXTRACTION_PLAN.md with Tier 1/2/3
  - [ ] 3.4 Run tests and commit

- [ ] 4.0 Consolidate Shortcut Components (Phase 4)
  - [ ] 4.1 Audit 5 shortcut-related components
  - [ ] 4.2 Design unified shortcut system
  - [ ] 4.3 Create UNIFIED_SHORTCUT_SYSTEM_DESIGN.md
  - [ ] 4.4 Run tests and commit

- [ ] 5.0 Verification and Final Documentation (Phase 5)
  - [ ] 5.1 Create SCREEN_TO_INTENT_MAPPING.md
  - [ ] 5.2 Run full test suite and linting
  - [ ] 5.3 Create CLEANUP_SUMMARY.md
  - [ ] 5.4 Final commit and handoff

---

## High-Level Overview

### What This Task Accomplishes

This task list guides the implementation of a comprehensive cleanup and refactoring of the KaRiya TUI codebase:

**Phase 1: Remove Orphaned Views (2 hours)**
- Delete 4 unused view files
- Clean up partially orphaned types
- Verify no broken references

**Phase 2: Fix Missing ImportProgressModel (2 hours)**
- Restore or recreate missing model
- Update broken references in app.go
- Verify compilation

**Phase 3: Extract Supporting Components (3 hours)**
- Document all 23 supporting components
- Create extraction plan with tiering
- Prepare for intent system integration

**Phase 4: Consolidate Shortcut Components (2 hours)**
- Audit 5 shortcut-related components
- Design unified shortcut system
- Document design for implementation

**Phase 5: Verification & Documentation (1 hour)**
- Create screen-to-intent mapping
- Run full quality checks
- Create cleanup summary

### Success Criteria

✅ All 4 orphaned views removed
✅ Missing ImportProgressModel resolved
✅ All 23 components documented with usage counts
✅ Screen-to-intent mapping complete (all 25 screens)
✅ Component extraction plan created (Tier 1/2/3)
✅ Shortcut system design documented
✅ All tests passing
✅ No linting/formatting errors
✅ Zero broken references

---

## Detailed Task Breakdown (Ready for Sub-Tasks)

**Ready to proceed with detailed sub-tasks?** Respond with "Go" to generate the complete breakdown with:
- Detailed step-by-step instructions for each sub-task
- Specific commands to run
- Verification steps
- Commit messages
- Detailed file operations

---

**Document Version**: 1.0
**Status**: 🚀 READY - Awaiting confirmation to generate detailed sub-tasks
**Next Step**: Respond with "Go" to proceed to Phase 2 (detailed sub-task generation)
- **Process Guide**: `/docs/rules/master-task-prompt.md`
