---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
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
| 'f' Filter shortcut not implemented | Users cannot filter events | [#46](https://github.com/baphled/KaRiya/issues/46) |
| Date range filter not implemented | Cannot filter by date | [#47](https://github.com/baphled/KaRiya/issues/47) |
| Company filter not implemented | Cannot filter by company | [#48](https://github.com/baphled/KaRiya/issues/48) |
| Categories filter not implemented | Cannot filter by category | [#49](https://github.com/baphled/KaRiya/issues/49) |
| Fact selection not connected | Cannot select facts from events | [#50](https://github.com/baphled/KaRiya/issues/50) |

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
| Review/Edit state is a stub | Cannot edit generated CV content | [#51](https://github.com/baphled/KaRiya/issues/51) |
| Hardcoded audience list | Cannot add custom audiences | N/A (part of #52) |
| Profile creation not available | Must pre-create profiles | [#52](https://github.com/baphled/KaRiya/issues/52) |

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
| PDF export defined but not implemented | Format option exists but doesn't work | [#53](https://github.com/baphled/KaRiya/issues/53) |
| Email export defined but not implemented | Destination option exists but doesn't work | [#54](https://github.com/baphled/KaRiya/issues/54) |
| Profile export returns hardcoded data | Not real profile data | [#55](https://github.com/baphled/KaRiya/issues/55) |
| CV export type not in defaults | Cannot export CVs directly | N/A (low priority) |

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
| Select type shows text input | Must type exact values, no dropdown | Medium | [#56](https://github.com/baphled/KaRiya/issues/56) |
| No validation for select options | Invalid values accepted | Medium | Part of #56 |
| Type assertion panic risk | Could crash on wrong types | **High** | [#57](https://github.com/baphled/KaRiya/issues/57) |
| Config not applied at runtime | Changes require restart | Medium | [#58](https://github.com/baphled/KaRiya/issues/58) |

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
| No file browser | Cannot select files | [#59](https://github.com/baphled/KaRiya/issues/59) |
| No CSV parsing | Cannot read CSV content | [#59](https://github.com/baphled/KaRiya/issues/59) |
| No database import | No repository calls | [#59](https://github.com/baphled/KaRiya/issues/59) |
| No field mapping | Cannot map CSV columns | [#59](https://github.com/baphled/KaRiya/issues/59) |

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
| No text inputs for editing | Cannot modify values | [#60](https://github.com/baphled/KaRiya/issues/60) |
| No entity loading | Cannot load events/facts/bursts | [#60](https://github.com/baphled/KaRiya/issues/60) |
| No persistence | Changes never saved | [#60](https://github.com/baphled/KaRiya/issues/60) |

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
| No item selection UI | Cannot select items to operate on | [#61](https://github.com/baphled/KaRiya/issues/61) |
| Delete operation not implemented | Doesn't delete anything | [#61](https://github.com/baphled/KaRiya/issues/61) |
| Tag operation not implemented | Doesn't add tags | [#61](https://github.com/baphled/KaRiya/issues/61) |
| Archive operation not implemented | Doesn't archive | [#61](https://github.com/baphled/KaRiya/issues/61) |
| Export operation not implemented | Doesn't export | [#61](https://github.com/baphled/KaRiya/issues/61) |

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

## GitHub Issues Created

### BrowseTimeline (5 issues)
1. [#46](https://github.com/baphled/KaRiya/issues/46) - feat(browse): implement 'f' filter shortcut
2. [#47](https://github.com/baphled/KaRiya/issues/47) - feat(browse): implement date range filtering
3. [#48](https://github.com/baphled/KaRiya/issues/48) - feat(browse): implement company filter
4. [#49](https://github.com/baphled/KaRiya/issues/49) - feat(browse): implement categories filter
5. [#50](https://github.com/baphled/KaRiya/issues/50) - feat(browse): connect fact selection UI

### GenerateCV (2 issues)
6. [#51](https://github.com/baphled/KaRiya/issues/51) - feat(cv): implement CV editing in review state
7. [#52](https://github.com/baphled/KaRiya/issues/52) - feat(cv): add profile creation from intent

### ExportArtifact (3 issues)
8. [#53](https://github.com/baphled/KaRiya/issues/53) - fix(export): remove or implement PDF export option
9. [#54](https://github.com/baphled/KaRiya/issues/54) - fix(export): remove or implement Email export option
10. [#55](https://github.com/baphled/KaRiya/issues/55) - fix(export): implement actual profile export

### ConfigureSystem (3 issues)
11. [#56](https://github.com/baphled/KaRiya/issues/56) - feat(config): implement select dropdown UI
12. [#57](https://github.com/baphled/KaRiya/issues/57) - fix(config): add type assertion safety checks
13. [#58](https://github.com/baphled/KaRiya/issues/58) - feat(config): apply config changes at runtime

### Secondary Intents (3 umbrella issues)
14. [#59](https://github.com/baphled/KaRiya/issues/59) - feat(import): implement ImportWizard functionality
15. [#60](https://github.com/baphled/KaRiya/issues/60) - feat(metadata): implement MetadataEditor functionality
16. [#61](https://github.com/baphled/KaRiya/issues/61) - feat(bulk): implement BulkOperations functionality

---

## Change History

| Date | Author | Description |
|------|--------|-------------|
| 2026-01-09 | AI + Human Review | Initial audit during Task 37 planning |
| 2026-01-09 | AI + Human Review | Created 16 GitHub issues (#46-#61) |
