# Component Extraction Plan

**Last Updated**: 2026-01-03
**Current State**: 48 model files in single `internal/cli/models/` directory
**Target State**: 25+ files organized into logical packages
**Timeline**: 3-4 sprints (6-8 weeks)

---

## Executive Summary

This plan outlines a phased approach to extract and reorganize 48 supporting components into logical packages. The goal is to:

1. **Reduce file count** from 48 → 25 (organized by purpose)
2. **Improve discoverability** through package organization
3. **Clarify dependencies** with explicit imports
4. **Enable component reuse** across projects
5. **Simplify testing** with package-level test suites

---

## Current State Analysis

### File Organization (Before)

```
internal/cli/models/  (48 files)
├── Core Infrastructure (5 files)
├── Dialog/Modal (3 files)
├── List/Table (8 files)
├── Data Entry (6 files)
├── Shortcut/Navigation (4 files)
├── CV-Specific (8 files)
├── Burst/Fact (8 files)
└── Other (6 files)
```

### Problems with Current Organization

1. **Flat structure**: 48 files in single directory = hard to navigate
2. **Mixed concerns**: UI, infrastructure, utilities all mixed
3. **Unclear dependencies**: No package boundaries
4. **Testing overhead**: Single test file per component
5. **Reusability**: Hard to extract components for other projects

---

## Target State

### Proposed Package Structure

```
internal/cli/
├── models/                      (Backward compatibility layer)
│   ├── base.go                 (Re-exports BaseStandardModel)
│   ├── messages.go             (Re-exports message types)
│   └── ...
│
├── components/                  (NEW: Extracted components)
│   ├── base/
│   │   ├── standard_model.go
│   │   ├── messages.go
│   │   └── errors.go
│   │
│   ├── lists/
│   │   ├── list.go
│   │   ├── patterns.go
│   │   ├── filter.go
│   │   ├── sort.go
│   │   └── search.go
│   │
│   ├── forms/
│   │   ├── form.go
│   │   ├── validators.go
│   │   ├── fields.go
│   │   └── ...editors
│   │
│   ├── dialogs/
│   │   ├── confirmation.go
│   │   ├── success.go
│   │   └── error.go
│   │
│   ├── navigation/
│   │   ├── shortcuts.go
│   │   ├── mapper.go
│   │   ├── customizer.go
│   │   └── help.go
│   │
│   ├── cv/
│   │   ├── generator.go
│   │   ├── config.go
│   │   ├── preview.go
│   │   ├── export.go
│   │   └── quality.go
│   │
│   ├── domain/
│   │   ├── bursts.go
│   │   ├── facts.go
│   │   ├── events.go
│   │   └── cards.go
│   │
│   └── utils/
│       ├── menu.go
│       ├── help.go
│       ├── bulk.go
│       └── import.go
│
└── intents/                     (Existing intent implementations)
    ├── capture_event/
    ├── browse_timeline/
    ├── generate_cv/
    ├── export_artifact/
    └── configure_system/
```

---

## Phase-by-Phase Extraction Plan

### Phase 1: Foundation & Infrastructure (Week 1)

**Goal**: Extract core infrastructure components

**Files to Extract**:
- `standard_model.go` → `components/base/standard_model.go`
- `messages.go` → `components/base/messages.go`
- `errors.go` → `components/base/errors.go`

**Steps**:
1. Create `internal/cli/components/base/` package
2. Move 3 files with minimal changes
3. Update imports in 25+ dependent files
4. Create re-export in `internal/cli/models/` for backward compatibility
5. Run full test suite
6. Commit: "refactor(components): extract core infrastructure"

**Effort**: 3 hours
**Risk**: Medium (25+ imports to update)

---

### Phase 2: List & Table Components (Week 1)

**Goal**: Extract list rendering infrastructure

**Files to Extract**:
- `list.go` → `components/lists/list.go`
- `list_patterns.go` → `components/lists/patterns.go`
- `filter.go` → `components/lists/filter.go`
- `sort.go` → `components/lists/sort.go`
- `search.go` → `components/lists/search.go`

**Steps**:
1. Create `internal/cli/components/lists/` package
2. Move 5 files (minimal changes)
3. Update imports in 12 dependent files
4. Create re-export in `internal/cli/models/`
5. Run full test suite
6. Commit: "refactor(components): extract list infrastructure"

**Effort**: 3 hours
**Risk**: Medium (12 imports to update)

---

### Phase 3: Form & Data Entry Components (Week 2)

**Goal**: Extract form infrastructure and editors

**Files to Extract**:
- `form.go` → `components/forms/form.go`
- `fact_editor.go` → `components/forms/fact_editor.go`
- `burst_editor.go` → `components/forms/burst_editor.go`
- `metadata_editor.go` → `components/forms/metadata_editor.go`
- `role_selector.go` → `components/forms/role_selector.go`
- `audience_configurator.go` → `components/forms/audience_configurator.go`

**Steps**:
1. Create `internal/cli/components/forms/` package
2. Move 6 files (minimal changes)
3. Update imports in 8 dependent files
4. Create re-export in `internal/cli/models/`
5. Run full test suite
6. Commit: "refactor(components): extract form infrastructure"

**Effort**: 4 hours
**Risk**: Medium (8 imports to update)

---

### Phase 4: Dialog & Feedback Components (Week 2)

**Goal**: Extract dialog and feedback components

**Files to Extract**:
- `confirmation_dialog.go` → `components/dialogs/confirmation.go`
- `success.go` → `components/dialogs/success.go`
- `error_handler.go` → `components/dialogs/error.go`

**Steps**:
1. Create `internal/cli/components/dialogs/` package
2. Move 3 files (minimal changes)
3. Update imports in 8 dependent files
4. Create re-export in `internal/cli/models/`
5. Run full test suite
6. Commit: "refactor(components): extract dialog infrastructure"

**Effort**: 2 hours
**Risk**: Low (8 imports to update)

---

### Phase 5: Navigation & Shortcuts (Week 3)

**Goal**: Extract and consolidate shortcut system

**Files to Extract**:
- `shortcut_handler.go` → `components/navigation/shortcuts.go`
- `context_shortcut_handler.go` → Consolidate into shortcuts.go
- `shortcut_mapper.go` → `components/navigation/mapper.go`
- `shortcut_customizer.go` → `components/navigation/customizer.go`
- `shortcut_help_system.go` → `components/navigation/help.go`

**Steps**:
1. Create `internal/cli/components/navigation/` package
2. Consolidate 5 files into 4 (merge context handler)
3. Update imports in 25+ dependent files
4. Create unified shortcut API
5. Create re-export in `internal/cli/models/`
6. Run full test suite
7. Commit: "refactor(components): extract & consolidate navigation"

**Effort**: 5 hours
**Risk**: High (25+ imports, consolidation complexity)

---

### Phase 6: CV Components (Week 3-4)

**Goal**: Extract CV-specific components

**Files to Extract**:
- `cv_generator.go` → `components/cv/generator.go`
- `cv_config_manager.go` → `components/cv/config.go`
- `cv_config_editor.go` → `components/cv/editor.go`
- `cv_preview.go` → `components/cv/preview.go`
- `cv_export_dialog.go` → `components/cv/export.go`
- `cv_export_progress.go` → `components/cv/progress.go`
- `cv_export_success.go` → `components/cv/success.go`
- `quality_indicator.go` → `components/cv/quality.go`

**Steps**:
1. Create `internal/cli/components/cv/` package
2. Move 8 files (minimal changes)
3. Update imports in 3 dependent files (GenerateCV, ExportArtifact)
4. Create re-export in `internal/cli/models/`
5. Run full test suite
6. Commit: "refactor(components): extract CV infrastructure"

**Effort**: 4 hours
**Risk**: Low (3 imports to update)

---

### Phase 7: Domain Components (Week 4)

**Goal**: Extract burst and fact components

**Files to Extract**:
- `burst_suggestion.go` → `components/domain/burst_suggestion.go`
- `burst_details.go` → `components/domain/burst_details.go`
- `burst_card.go` → `components/domain/burst_card.go`
- `burst_list.go` → `components/domain/burst_list.go`
- `fact_details.go` → `components/domain/fact_details.go`
- `fact_card.go` → `components/domain/fact_card.go`
- `fact_search.go` → `components/domain/fact_search.go`
- `facts_results.go` → `components/domain/facts_results.go`

**Steps**:
1. Create `internal/cli/components/domain/` package
2. Move 8 files (minimal changes)
3. Update imports in 5 dependent files
4. Create re-export in `internal/cli/models/`
5. Run full test suite
6. Commit: "refactor(components): extract domain components"

**Effort**: 4 hours
**Risk**: Low (5 imports to update)

---

### Phase 8: Utility Components (Week 4-5)

**Goal**: Extract remaining utility components

**Files to Extract**:
- `view_event.go` → `components/utils/view_event.go`
- `details.go` → `components/utils/details.go`
- `metadata_review.go` → `components/utils/metadata_review.go`
- `bulk_operations.go` → `components/utils/bulk_operations.go`
- `import_review.go` → `components/utils/import_review.go`
- `action_menu.go` → `components/utils/action_menu.go`
- `menu.go` → `components/utils/menu.go`
- `help.go` → `components/utils/help.go`

**Steps**:
1. Create `internal/cli/components/utils/` package
2. Move 8 files (minimal changes)
3. Update imports in 8 dependent files
4. Create re-export in `internal/cli/models/`
5. Run full test suite
6. Commit: "refactor(components): extract utility components"

**Effort**: 5 hours
**Risk**: Low (8 imports to update)

---

### Phase 9: Cleanup & Verification (Week 5)

**Goal**: Clean up re-exports and verify everything works

**Steps**:
1. Update all re-exports in `internal/cli/models/`
2. Run full test suite
3. Run linting and formatting
4. Update documentation
5. Create migration guide for future development
6. Commit: "refactor(components): complete extraction & consolidation"

**Effort**: 3 hours
**Risk**: Low

---

## Consolidation Opportunities

### Shortcut System Consolidation (Phase 5)

Current: 5 files (shortcut_handler.go, context_shortcut_handler.go, shortcut_mapper.go, shortcut_customizer.go, shortcut_help_system.go)

**Proposed**: Consolidate into 4 files:
- `shortcuts.go` - Combined handler + context handler
- `mapper.go` - Shortcut mapping
- `customizer.go` - Customization UI
- `help.go` - Help system

**Benefits**:
- Reduce file count by 1
- Clearer API boundary
- Easier to understand shortcut flow
- Simpler testing

---

## Backward Compatibility

### Re-export Strategy

All extracted components will be re-exported from `internal/cli/models/` to maintain backward compatibility:

```go
// internal/cli/models/base.go
package models

// Re-export from components/base
export (
    BaseStandardModel from "github.com/baphled/kariya/internal/cli/components/base"
)

// internal/cli/models/messages.go
// ... similar re-exports
```

This allows:
- Existing code to continue working without changes
- Gradual migration to new imports
- Optional deprecation warnings

---

## Risk Mitigation

### High-Risk Operations

1. **Phase 5 (Shortcuts)**: 25+ imports to update
   - **Mitigation**: Automated import rewriting, careful testing
   - **Rollback**: Single git revert if needed

2. **Consolidation**: Merging related files
   - **Mitigation**: Small, atomic commits
   - **Rollback**: Easy to split back apart

### Testing Strategy

1. **Before each phase**: Run full test suite
2. **During extraction**: Update imports incrementally
3. **After each phase**: Re-run full test suite
4. **Final verification**: Race condition detection, coverage check

---

## Success Criteria

### Per-Phase Criteria

- ✅ All tests passing (978+ specs)
- ✅ No race conditions detected
- ✅ Code coverage maintained (≥85%)
- ✅ No linting errors
- ✅ Backward compatibility maintained
- ✅ Documentation updated

### Overall Success Criteria

- ✅ File count reduced from 48 → 25
- ✅ All components organized into logical packages
- ✅ Dependencies clearly visible
- ✅ All tests passing (100% pass rate)
- ✅ Zero race conditions
- ✅ Coverage maintained at ≥85%
- ✅ Shortcut system consolidated (5 → 4 files)

---

## Implementation Checklist

### Phase 1: Foundation
- [ ] Create `components/base/` package
- [ ] Move 3 files
- [ ] Update 25+ imports
- [ ] Create re-exports
- [ ] Test & commit

### Phase 2: Lists
- [ ] Create `components/lists/` package
- [ ] Move 5 files
- [ ] Update 12 imports
- [ ] Create re-exports
- [ ] Test & commit

### Phase 3: Forms
- [ ] Create `components/forms/` package
- [ ] Move 6 files
- [ ] Update 8 imports
- [ ] Create re-exports
- [ ] Test & commit

### Phase 4: Dialogs
- [ ] Create `components/dialogs/` package
- [ ] Move 3 files
- [ ] Update 8 imports
- [ ] Create re-exports
- [ ] Test & commit

### Phase 5: Navigation
- [ ] Create `components/navigation/` package
- [ ] Move 5 files (consolidate 1)
- [ ] Update 25+ imports
- [ ] Create unified API
- [ ] Create re-exports
- [ ] Test & commit

### Phase 6: CV
- [ ] Create `components/cv/` package
- [ ] Move 8 files
- [ ] Update 3 imports
- [ ] Create re-exports
- [ ] Test & commit

### Phase 7: Domain
- [ ] Create `components/domain/` package
- [ ] Move 8 files
- [ ] Update 5 imports
- [ ] Create re-exports
- [ ] Test & commit

### Phase 8: Utils
- [ ] Create `components/utils/` package
- [ ] Move 8 files
- [ ] Update 8 imports
- [ ] Create re-exports
- [ ] Test & commit

### Phase 9: Cleanup
- [ ] Verify all re-exports
- [ ] Final test run
- [ ] Linting & formatting
- [ ] Update documentation
- [ ] Create migration guide
- [ ] Final commit

---

## Timeline & Resources

### Estimated Timeline
- **Total Effort**: 31 hours
- **Timeline**: 4-5 weeks (6-8 hours/week)
- **Resources**: 1 senior engineer + automated tooling

### Weekly Breakdown
- **Week 1**: Phases 1-2 (6 hours)
- **Week 2**: Phases 3-4 (6 hours)
- **Week 3**: Phases 5-6 (9 hours)
- **Week 4**: Phases 7-8 (9 hours)
- **Week 5**: Phase 9 (3 hours)

---

## Post-Extraction Opportunities

Once extraction is complete:

1. **Extract to separate module**: `github.com/baphled/kariya-components`
2. **Create component library**: Reusable across projects
3. **Add composition patterns**: Reduce component duplication
4. **Implement component registry**: Dynamic component loading
5. **Create component storybook**: Component showcase & testing

---

**Status**: 📋 Ready for Implementation
**Generated**: 2026-01-03
**Next Step**: Execute Phase 1 when team is ready

