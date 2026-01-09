# E2E Implementation Gaps

**Last Updated**: 2026-01-09
**Task**: Task 37 - E2E Integration Test Enhancement

This document tracks implementation gaps discovered during E2E test development and audit.

---

## Overview

During the E2E integration test planning phase, a comprehensive audit was conducted of all 10 intents. This revealed several implementation gaps that limit full E2E testing capabilities.

### Summary by Intent Type

| Category | Fully Functional | Partial Implementation | UI Shell Only |
|----------|------------------|------------------------|---------------|
| **Core Intents** | 1 (CaptureEvent) | 4 (BrowseTimeline, GenerateCV, ExportArtifact, ConfigureSystem) | 0 |
| **Secondary Intents** | 2 (BurstManagement, FactManagement) | 0 | 3 (ImportWizard, MetadataEditor, BulkOperations) |

---

## Core Intents

### CaptureEvent

**Status**: Fully Functional

**Minor Gaps**:
| Gap | Impact | Severity |
|-----|--------|----------|
| Confirmed bursts not persisted to DB after modal | Burst acceptance doesn't save | Low |
| Enrichment errors logged but not shown to user | Silent failures | Low |

---

### BrowseTimeline

**Status**: Partial Implementation

**Critical Gaps**:
| Gap | Impact | GitHub Issue |
|-----|--------|--------------|
| 'f' Filter shortcut not implemented | Users cannot filter events | TODO |
| Date range filter not implemented | Cannot filter by date | TODO |
| Company filter not implemented | Cannot filter by company | TODO |
| Categories filter not implemented | Cannot filter by category | TODO |
| Fact selection not connected | Cannot select facts from events | TODO |

**Notes**:
- The `TimelineFilters` struct defines `DateFrom`, `DateTo`, `Companies`, and `Categories` fields
- `applyFilters()` only uses `SearchText` and `Tags`
- Footer shows "f Filter" but no handler exists for 'f' key

---

### GenerateCV

**Status**: Partial Implementation

**Gaps**:
| Gap | Impact | GitHub Issue |
|-----|--------|--------------|
| Review/Edit state is a stub | Cannot edit generated CV content | TODO |
| Hardcoded audience list | Cannot add custom audiences | TODO |
| Profile creation not available | Must pre-create profiles | TODO |

**Notes**:
- The `viewReview()` method shows placeholder text: "(Full editing interface would be implemented here)"
- Audiences are hardcoded to: `["hiring_manager", "recruiter", "peer"]`
- CV is returned in result but never persisted to database

---

### ExportArtifact

**Status**: Partial Implementation

**Gaps**:
| Gap | Impact | GitHub Issue |
|-----|--------|--------------|
| PDF export defined but not implemented | Format option exists but doesn't work | TODO |
| Email export defined but not implemented | Destination option exists but doesn't work | TODO |
| Profile export returns hardcoded data | Not real profile data | TODO |
| CV export type not in defaults | Cannot export CVs directly | TODO |

**Notes**:
- `ExportFormatPDF` and `ExportDestinationEmail` are defined in constants
- `generateProfilePreview()` returns hardcoded JSON with "John Doe"
- `startExport()` doesn't handle all artifact types

---

### ConfigureSystem

**Status**: Partial Implementation

**Gaps**:
| Gap | Impact | Severity | GitHub Issue |
|-----|--------|----------|--------------|
| Select type shows text input | Must type exact values, no dropdown | Medium | TODO |
| No validation for select options | Invalid values accepted | Medium | TODO |
| Type assertion panic risk | Could crash on wrong types | **High** | TODO |
| Config not applied at runtime | Changes require restart | Medium | TODO |

**Notes**:
- Lines 763-820 contain direct type assertions without safety checks
- Settings are saved to file but not applied to running application
- No `Reset to Default` functionality

---

## Secondary Intents

### BurstManagement

**Status**: Fully Functional

**Minor Gaps**:
| Gap | Impact | Severity |
|-----|--------|----------|
| Edit form is placeholder | Shows message instead of inputs | Low |
| Filter/sort not wired to UI | Cannot filter/sort bursts | Low |

---

### FactManagement

**Status**: Fully Functional

**Minor Gaps**:
| Gap | Impact | Severity |
|-----|--------|----------|
| Editor shows values but no inputs | Display-only in edit mode | Low |
| Quality/source filtering not implemented | Cannot filter facts | Low |

---

### ImportWizard

**Status**: UI Shell Only

**Critical Gaps** (No actual functionality):
| Gap | Impact | GitHub Issue |
|-----|--------|--------------|
| No file browser | Cannot select files | TODO |
| No CSV parsing | Cannot read CSV content | TODO |
| No database import | No repository calls | TODO |
| No field mapping | Cannot map CSV columns | TODO |

**Notes**:
- The `ImportWizardContext` has no repository or service references
- All progress tracking is manual counter increments
- `FilePath` must be set programmatically

---

### MetadataEditor

**Status**: UI Shell Only

**Critical Gaps** (No actual functionality):
| Gap | Impact | GitHub Issue |
|-----|--------|--------------|
| No text inputs for editing | Cannot modify values | TODO |
| No entity loading | Cannot load events/facts/bursts | TODO |
| No persistence | Changes never saved | TODO |

**Notes**:
- The `MetadataEditorContext` has no repository or service references
- `LoadMetadata()` just copies a map internally
- `SetFieldValue()` updates internal state but nowhere to save

---

### BulkOperations

**Status**: UI Shell Only

**Critical Gaps** (No actual functionality):
| Gap | Impact | GitHub Issue |
|-----|--------|--------------|
| No item selection UI | Cannot select items to operate on | TODO |
| Delete operation not implemented | Doesn't delete anything | TODO |
| Tag operation not implemented | Doesn't add tags | TODO |
| Archive operation not implemented | Doesn't archive | TODO |
| Export operation not implemented | Doesn't export | TODO |

**Notes**:
- The `BulkOperationsContext` has no repository or service references
- `AvailableOps` lists operations but none are implemented
- `AffectedItemCount` is set but no item selection mechanism exists

---

## E2E Test Implications

### Fully Testable
- CaptureEvent: Complete workflow testing
- BurstManagement: Full CRUD testing
- FactManagement: Full CRUD testing

### Partially Testable
- BrowseTimeline: Navigation and display only (no filters)
- GenerateCV: Generation and export work (no editing)
- ExportArtifact: Events/Facts/Bursts to JSON/CSV/YAML/TXT (no PDF/Email)
- ConfigureSystem: All domains, persistence works (careful with types)

### Navigation Only Testing
- ImportWizard: State machine and navigation only
- MetadataEditor: State machine and navigation only
- BulkOperations: State machine and navigation only

---

## Recommendations

### High Priority (Critical for User Experience)
1. **BrowseTimeline filters** - Documented feature that doesn't work
2. **ConfigureSystem type safety** - Potential crashes
3. **Remove/implement PDF/Email export** - Misleading options

### Medium Priority (Feature Completion)
4. **GenerateCV editing** - Common user need
5. **ImportWizard full implementation** - Important for onboarding
6. **BurstManagement edit form** - Complete the CRUD

### Low Priority (Polish)
7. **ConfigureSystem runtime apply** - Convenience
8. **FactManagement filtering** - Nice to have
9. **BulkOperations implementation** - Power user feature

---

## GitHub Issues to Create

### BrowseTimeline (5 issues)
1. `feat(browse): implement 'f' filter shortcut`
2. `feat(browse): implement date range filtering`
3. `feat(browse): implement company filter`
4. `feat(browse): implement categories filter`
5. `feat(browse): connect fact selection UI`

### GenerateCV (2 issues)
6. `feat(cv): implement CV editing in review state`
7. `feat(cv): add profile creation from intent`

### ExportArtifact (3 issues)
8. `fix(export): remove or implement PDF export option`
9. `fix(export): remove or implement Email export option`
10. `fix(export): implement actual profile export`

### ConfigureSystem (3 issues)
11. `feat(config): implement select dropdown UI`
12. `fix(config): add type assertion safety checks`
13. `feat(config): apply config changes at runtime`

### Secondary Intents (3 umbrella issues)
14. `feat(import): implement ImportWizard functionality`
15. `feat(metadata): implement MetadataEditor functionality`
16. `feat(bulk): implement BulkOperations functionality`

---

## Change History

| Date | Author | Description |
|------|--------|-------------|
| 2026-01-09 | AI + Human Review | Initial audit during Task 37 planning |
