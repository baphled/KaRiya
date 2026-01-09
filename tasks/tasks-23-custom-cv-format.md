# Task 23: CV Structure Selection - Standard and Narrative Templates

## Overview

**Goal**: Add CV structure selection to KaRiya's CV generation system. Users select a structure (Standard or Narrative) *before* CV generation, which determines how the CV content is organized and presented. This is separate from export format (Text/Markdown/YAML) which determines the output file type.

**Time Estimate**: ~2 weeks (5 phases)

**Prerequisites**:
- Existing CV generation system functional (✅ Complete)
- Export service supports Text/Markdown/YAML (✅ Complete)
- Understanding of `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` (✅ Complete)

**Related Docs**:
- `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` - Original proposal (needs update after implementation)
- `docs/guides/CV_GENERATION_GUIDE.md` - CV generation workflow (needs update after implementation)
- `internal/service/career/cv/export_service.go` - Current export implementation

---

## Architecture Overview

### Key Concepts

| Concept | Description | When Selected |
|---------|-------------|---------------|
| **CV Structure** | How content is organized (sections, emphasis, filtering) | Before generation |
| **Export Format** | File output format (text, markdown, yaml) | After preview |

### CV Structures

| Structure | Description | Use Case |
|-----------|-------------|----------|
| `standard` | Traditional CV with Experience, Projects, Skills, Summary sections | Most job applications |
| `narrative` | Language-agnostic professional format with Core Strengths, Technologies, What I Bring sections | Emphasizing pragmatic, cross-domain experience |

### State Flow

```
Select Profile → Select Audience → Select Structure → Generate → Preview → Export Format
                                        ↑                           ↑
                                  Standard | Narrative         Text | Markdown | YAML
```

---

## Success Criteria

- [ ] CV Structure selection appears before generation (Standard/Narrative)
- [ ] Narrative structure renders CV with correct sections in preview
- [ ] Narrative structure exports correctly to Text/Markdown formats
- [ ] YAML always uses standard structure (it's data format)
- [ ] Profile configuration available in Configure System → Profile (Phase 2 - Future)
- [ ] All tests pass with zero regressions
- [ ] Code coverage maintained >87%
- [ ] User guide complete with examples (after implementation)

---

## Implementation Phases

### Phase 1: Cleanup - Remove Incorrect Implementation

**Objective**: Remove the incorrectly implemented "Custom" export format

**Background**: Phase 1 was previously implemented incorrectly - "Custom" was added as an export format (alongside Text/Markdown/YAML) when it should be a CV structure selected before generation.

#### Files to Delete
- [x] `internal/service/career/cv/export_custom.go` - Wrong architecture
- [x] `internal/service/career/cv/export_custom_test.go` - Tests for deleted code

#### Files to Modify

##### 1.1: Remove ExportFormatCustom from export_service.go
- [x] Remove `ExportFormatCustom ExportFormat = "custom"` constant (N/A - never existed after reset)
- [x] Remove custom case from `getFileExtension()` (N/A - never existed after reset)

##### 1.2: Remove CVExportFormatCustom from generate_cv.go
- [x] Remove `CVExportFormatCustom CVExportFormat = "custom"` constant (N/A - never existed after reset)

##### 1.3: Revert generate_cv_intent.go to 3 export formats
- [x] Revert `updateExportSelectFormat()` to handle 3 formats (maxIndex = 2) (N/A - never existed after reset)
- [x] Revert `viewExportSelectFormat()` to show only Text, Markdown, YAML (N/A - never existed after reset)
- [x] Remove custom case from `exportCVAsync()` switch (N/A - never existed after reset)
- [x] Remove "Custom" from format name display in view methods (N/A - never existed after reset)

#### Verification
- [x] `go build ./...` compiles successfully
- [x] `go test ./...` passes

**Note**: Phase 1 cleanup was simplified by resetting the branch to `origin/next` which removed all incorrect commits.

---

### Phase 2: Add CV Structure Types and Selection State (TDD)

**Objective**: Add new state for CV structure selection between audience and generation

#### Step 2.1: Write Failing Tests First
**File**: `internal/cli/intents/generate_cv_structure_test.go` (new)

Tests to implement:
- [x] Transitions from select_audience to select_structure on enter
- [x] Transitions from select_structure to generating on enter
- [x] Goes back to select_audience on esc from select_structure
- [x] Cancels on q from select_structure
- [x] Shows Standard and Narrative structure options
- [x] Shows descriptions for each structure
- [x] Highlights currently selected structure
- [x] Navigates between structure options with j/k
- [x] Navigates between structure options with up/down
- [x] Stores selected structure in model
- [x] Defaults to Standard structure
- [x] Includes structure in breadcrumbs
- [x] Includes selected structure in result

#### Step 2.2: Add Types and Constants
**File**: `internal/cli/intents/generate_cv.go`

- [x] Add `CVStructure` type
- [x] Add `CVStructureStandard CVStructure = "standard"` constant
- [x] Add `CVStructureNarrative CVStructure = "narrative"` constant
- [x] Add `GenerateCVStateSelectStructure GenerateCVState = "select_structure"` constant
- [x] Add `selectedCVStructure CVStructure` to `GenerateCVModel`
- [x] Add `structureIndex int` to `GenerateCVModel`
- [x] Add `SelectedStructure CVStructure` to `GenerateCVResult`

#### Step 2.3: Implement State Handler
**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Add `case GenerateCVStateSelectStructure:` to `Update()` switch
- [x] Implement `updateSelectStructure(msg tea.Msg) tea.Cmd`
  - j/k or up/down: navigate between Standard (0) and Narrative (1)
  - enter: set structure, transition to generating, call generateCVAsync()
  - esc: go back to select_audience
  - q/ctrl+c: cancel intent
  - m: return to main menu
- [x] Implement `viewSelectStructure() string`

#### Step 2.4: Update State Transitions
- [x] Modify `updateSelectAudience()`: enter → `GenerateCVStateSelectStructure` (not generating)
- [x] Remove `generateCVAsync()` call from `updateSelectAudience()` (moved to structure selection)

#### Step 2.5: Update Helper Methods
- [x] Update `getStateContent()` - add case for select_structure
- [x] Update `getBreadcrumbs()` - add "Select Structure" crumb
- [x] Update `getContextHelp()` - add help text for structure selection

#### Step 2.6: Update Result
- [x] Modify `setCompleted()` to include `SelectedStructure` in result and metadata

---

### Phase 3: Add Structure-Aware Preview (TDD)

**Objective**: Preview displays CV using the selected structure

#### Step 3.1: Write Failing Tests First

Tests for Standard Structure Preview:
- [x] Renders CV name as header
- [x] Renders all sections in order
- [x] Renders experience with company headers and dates
- [x] Renders bullets for each group

Tests for Narrative Structure Preview:
- [x] Renders profile header with name, role, location, contact
- [x] Renders Summary section
- [x] Renders Core Strengths section
- [x] Renders Languages & Technologies section
- [x] Renders Selected Experience with confidence filtering
- [x] Renders What I Bring section
- [x] Uses default values when sections are empty
- [x] Filters bullets by confidence >= 0.75

#### Step 3.2: Create Helper File
**File**: `internal/cli/intents/generate_cv_helpers.go` (new)

- [x] Add `NarrativeProfileData` struct
- [x] Add `DefaultNarrativeProfile()` function (hardcoded for Phase 1)
- [x] Add `extractStrengthsFromSections()` helper
- [x] Add `extractTechnologiesFromSections()` helper
- [x] Add `extractValuePropositions()` helper

#### Step 3.3: Split Preview Method
**File**: `internal/cli/intents/generate_cv_intent.go`

- [x] Modify `viewPreview()` to switch on `selectedCVStructure`
- [x] Rename current preview logic to `viewPreviewStandard()`
- [x] Implement `viewPreviewNarrative()` with:
  - Profile header (using DefaultNarrativeProfile())
  - Summary section
  - Core Strengths section (using extractStrengthsFromSections or defaults)
  - Languages & Technologies section (using extractTechnologiesFromSections or defaults)
  - Selected Experience section (filtered by confidence >= 0.75)
  - What I Bring section (using extractValuePropositions or defaults)

---

### Phase 4: Add Structure-Aware Export (TDD)

**Objective**: Export methods use selected structure to render content

#### Step 4.1: Write Failing Tests First
**File**: `internal/service/career/cv/export_service_test.go`

- [ ] Exports standard structure to text format
- [ ] Exports standard structure to markdown format
- [ ] Exports narrative structure to text format
- [ ] Exports narrative structure to markdown format
- [ ] Always uses standard structure for YAML format
- [ ] Returns error for nil CV

#### Step 4.2: Create Export Helper File
**File**: `internal/service/career/cv/cv_helpers.go` (new)

- [ ] Add `CVStructure` type and constants
- [ ] Add `NarrativeProfileData` struct
- [ ] Add `DefaultNarrativeProfile()` function
- [ ] Add `extractStrengthsFromSections()` helper
- [ ] Add `extractTechnologiesFromSections()` helper
- [ ] Add `extractValuePropositions()` helper

#### Step 4.3: Add Export Method with Structure
**File**: `internal/service/career/cv/export_service.go`

- [ ] Add `Export()` method with structure parameter
- [ ] Add `renderStandard()` internal method (delegates to existing ExportToText/Markdown)
- [ ] Add `renderNarrative()` internal method (narrative structure rendering)
- [ ] Keep existing `ExportToText()`, `ExportToMarkdown()`, `ExportToYAML()` unchanged (backward compat)

#### Step 4.4: Wire Export in Intent
**File**: `internal/cli/intents/generate_cv_intent.go`

- [ ] Update `exportCVAsync()` to use new `Export()` method with structure

---

### Phase 5: Profile Configuration (Future - Not in Current Scope)

**Objective**: Allow users to configure narrative profile in Configure System

This phase is deferred. Current implementation uses hardcoded profile data.

---

## File Summary

### Files to Delete
| File | Reason |
|------|--------|
| `internal/service/career/cv/export_custom.go` | Wrong architecture |
| `internal/service/career/cv/export_custom_test.go` | Tests for deleted code |

### Files to Create
| File | Purpose |
|------|---------|
| `internal/cli/intents/generate_cv_structure_test.go` | Tests for structure selection state |
| `internal/cli/intents/generate_cv_helpers.go` | Shared helpers for narrative preview |
| `internal/service/career/cv/cv_helpers.go` | Shared helpers for narrative export |

### Files to Modify
| File | Changes |
|------|---------|
| `internal/service/career/cv/export_service.go` | Remove custom format, add Export() with structure |
| `internal/service/career/cv/export_service_test.go` | Add structure-aware export tests |
| `internal/cli/intents/generate_cv.go` | Add CVStructure type, state, model fields |
| `internal/cli/intents/generate_cv_intent.go` | Add state handler, structure-aware preview, wire export |

---

## Test Count Estimate

| Category | Tests |
|----------|-------|
| Structure selection state | ~12 |
| Preview tests | ~10 |
| Export tests | ~6 |
| **Total new tests** | ~28 |

---

## Narrative CV Output Format

The narrative structure produces content with:

```markdown
# Yomi Colledge

**Senior Software Engineer / Technical Consultant**  
Remote (UK)  
Email: [yomi@boodah.net](mailto:yomi@boodah.net)  
GitHub: https://github.com/baphled  
Portfolio: http://boodah.net

---

## Summary
[Professional summary from CV sections or default]

---

## Core Strengths
- Language-agnostic backend and systems engineering  
- System design and architectural ownership  
[...]

---

## Languages & Technologies
**Languages:** Ruby, Go, PHP, C/C++, JavaScript, Shell  
**Frontend:** Vue.js  
**Systems:** Linux, SQL, APIs, CI/CD, automation

---

## Selected Experience

### Company Name
*Jan 2020 - Present*

- High confidence bullet (>= 0.75)  
- Another achievement  

---

## What I Bring
- Languages as tools, not identity  
- Calm handling of complexity  
[...]

---

**References available on request.**
```

---

## Rollback Plan

If issues arise:

1. Revert commits atomically
2. Structure selection is additive - no existing functionality broken
3. Export formats remain unchanged (Text/Markdown/YAML)

**Impact**: Zero - feature is additive

---

## Documentation Updates (After Implementation)

### Files to Update
| File | Changes |
|------|---------|
| `docs/CUSTOM_CV_FORMAT_PROPOSAL.md` | Rewrite to reflect correct architecture |
| `docs/guides/CV_GENERATION_GUIDE.md` | Add "CV Structure Selection" section |
| `docs/guides/CV_TROUBLESHOOTING.md` | Add narrative structure troubleshooting |
| `AGENTS.md` | Update to mention CV structures |

### Files to Create
| File | Purpose |
|------|---------|
| `docs/guides/NARRATIVE_CV_GUIDE.md` | Explain narrative structure, when to use it |

---

## Resolved Decisions

| Question | Decision | Rationale |
|----------|----------|-----------|
| Terminology | Standard/Narrative (not Default/Custom) | Clearer meaning |
| Selection point | Before generation | Structure affects presentation |
| Helper duplication | Keep separate (intents vs service) | Simpler for now |
| Export backward compat | Keep existing methods unchanged | No breaking changes |
| YAML + Narrative | YAML always standard | YAML is data, not presentation |
| Documentation timing | After implementation | Avoid updating twice |

---

## References

- [CV Generation Service](../internal/service/career/cv/cv_generation_service.go)
- [Export Service](../internal/service/career/cv/export_service.go)
- [GenerateCV Intent](../internal/cli/intents/generate_cv_intent.go)

---

**Last Updated**: 2026-01-09  
**Status**: Phase 3 COMPLETE - Structure-aware preview implemented  
**Next Step**: Phase 4 - Add structure-aware export (TDD)
